package cpu

import (
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

func (c *CPU) Jump(addr uint32) {
	c.fe.In.PC = addr
}

func (c *CPU) Run() {
	c.fe.Run()
	c.de.Run()
	c.ex.Run()
	c.wb.Run()

	c.BusHandler.HandleBus(&c.fe.Out.Bus)
}

func (c *CPU) ClockRegisters() {
	// wb <- me
	c.wb.In.PC = c.ex.Out.PC

	// ex <- de
	c.ex.In.PC = c.de.Out.PC

	// de <- fe
	c.de.In.PC = c.fe.Out.PC
	c.de.In.IR = c.fe.Out.IR

	// de <- wb
	c.de.In.Result = c.wb.Out.Result

	// fe <- fe
	c.fe.In.PC = c.fe.Out.PC
	c.fe.In.Bus = c.fe.Out.Bus
}
