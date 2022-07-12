package cpu

import (
	"bytes"
	"fmt"

	"github.com/rj45/rj32/llemu/cpu/de"
	"github.com/rj45/rj32/llemu/cpu/ex"
	"github.com/rj45/rj32/llemu/cpu/fe"
	"github.com/rj45/rj32/llemu/cpu/wb"
	"github.com/rj45/rj32/llemu/mem"
)

type CPU struct {
	BusHandler mem.BusHandler

	fe fe.Stage
	de de.Stage
	ex ex.Stage
	wb wb.Stage
}

func NewCPU(bus mem.BusHandler) *CPU {
	c := &CPU{
		BusHandler: bus,
	}

	c.fe.In.Log = &bytes.Buffer{}
	c.de.In.Log = &bytes.Buffer{}
	c.ex.In.Log = &bytes.Buffer{}
	c.wb.In.Log = &bytes.Buffer{}

	c.wb.Writer = &c.de

	return c
}

func (c *CPU) Jump(addr uint32) {
	c.fe.In.Npc = addr
	c.fe.In.Pc = addr
}

func (c *CPU) Run() {
	c.fe.Run()
	c.de.Run()
	c.ex.Run()
	c.wb.Run()

	c.BusHandler.HandleBus(&c.fe.Out.Bus)
}

func (c *CPU) ClockRegisters() {
	log := c.wb.In.Log
	fmt.Println(log.String())
	log.Reset()

	// wb <- ex
	c.wb.In.Log = c.ex.In.Log
	c.wb.In.PC = c.ex.Out.PC
	c.wb.In.IR = c.ex.Out.IR
	c.wb.In.RegWen = c.ex.Out.RegWen
	c.wb.In.Rd = c.ex.Out.Rd
	c.wb.In.Result = c.ex.Out.Result

	// ex <- de
	c.ex.In.Log = c.de.In.Log
	c.ex.In.PC = c.de.Out.PC
	c.ex.In.IR = c.de.Out.IR
	c.ex.In.L = c.de.Out.L
	c.ex.In.R = c.de.Out.R
	c.ex.In.AluOp = c.de.Out.AluOp
	c.ex.In.Sub = c.de.Out.Sub
	c.ex.In.RegWen = c.de.Out.RegWen
	c.ex.In.Rd = c.de.Out.Rd

	// de <- fe
	c.de.In.Log = c.fe.In.Log
	c.de.In.PC = c.fe.Out.Pc
	c.de.In.IR = c.fe.Out.IR

	// fe <- fe
	c.fe.In.Pc = c.fe.Out.Pc
	c.fe.In.Bus = c.fe.Out.Bus
	c.fe.In.Npc = c.fe.Out.Npc

	// fe <- wb
	c.fe.In.Log = log
}
