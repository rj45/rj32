package de

import (
	"fmt"

	"github.com/rj45/rj32/llemu/rj32"
)

type In struct {
	PC     uint32
	IR     uint32
	Result uint32
}

type Out struct {
	PC uint32
}

type Stage struct {
	In
	Out
}

func (s *Stage) Run() {
	ir := rj32.Inst(s.In.IR)
	_ = ir

	fmt.Printf("%08x:  %08x  %-8s %s, %s, %s \n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Rs1(), ir.Rs2())

	s.Out.PC = s.In.PC
}
