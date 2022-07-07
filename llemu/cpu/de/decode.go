package de

import (
	"fmt"

	"github.com/rj45/rj32/llemu/rj32"
)

type In struct {
	PC uint32
	IR uint32

	Result uint32
	Rd     rj32.Reg
	RegWen bool
}

type Out struct {
	PC uint32
	IR rj32.Inst

	L uint32
	R uint32
}

type Stage struct {
	In
	Out

	Regs [32]uint32
}

func (s *Stage) Run() {
	ir := rj32.Inst(s.In.IR)
	s.Out.IR = ir

	switch ir.Fmt() {
	case rj32.N:
		fmt.Printf("%08x:  %08x  %-8s\n", s.In.PC, s.In.IR, ir.Opcode())
	case rj32.R:
		fmt.Printf("%08x:  %08x  %-8s %s, %s, %s\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Rs1(), ir.Rs2())
	case rj32.I:
		fmt.Printf("%08x:  %08x  %-8s %s, %s, %d\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Rs1(), ir.Imm())
	case rj32.B:
		fmt.Printf("%08x:  %08x  %-8s %s, %s, %d\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rs1(), ir.Rs2(), ir.Imm())
	case rj32.U:
		fmt.Printf("%08x:  %08x  %-8s %s, %#x\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Imm()<<12)
	case rj32.J:
		fmt.Printf("%08x:  %08x  %-8s %s, %#x\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Imm())
	}

	s.Out.PC = s.In.PC

	// write happens in first half of the cycle
	if s.In.RegWen {
		s.Regs[s.In.Rd] = s.In.Result
	}

	// read happens in the second half
	s.Out.L = s.Regs[ir.Rs1()]
	s.Out.R = s.Regs[ir.Rs2()]
}
