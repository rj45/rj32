package de

import (
	"bytes"
	"fmt"

	"github.com/rj45/rj32/llemu/rj32"
)

type In struct {
	PC uint32
	IR uint32

	Log *bytes.Buffer
}

type Out struct {
	PC uint32
	IR rj32.Inst

	L uint32
	R uint32

	AluOp rj32.Opcode
	Sub   bool

	Rd     rj32.Reg
	RegWen bool
}

type Stage struct {
	In
	Out

	Regs [32]uint32
}

func (s *Stage) WriteReg(reg rj32.Reg, val uint32) {
	s.Regs[reg] = val
}

func (s *Stage) Run() {
	ir := rj32.Inst(s.In.IR)
	s.Out.IR = ir

	switch ir.Fmt() {
	case rj32.N:
		fmt.Fprintf(s.In.Log, "%08x:  %08x  %-8s\n", s.In.PC, s.In.IR, ir.Opcode())
	case rj32.R:
		fmt.Fprintf(s.In.Log, "%08x:  %08x  %-8s %s, %s, %s\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Rs1(), ir.Rs2())
	case rj32.I:
		fmt.Fprintf(s.In.Log, "%08x:  %08x  %-8s %s, %s, %d\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Rs1(), ir.Imm())
	case rj32.B:
		fmt.Fprintf(s.In.Log, "%08x:  %08x  %-8s %s, %s, %d\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rs1(), ir.Rs2(), ir.Imm())
	case rj32.U:
		fmt.Fprintf(s.In.Log, "%08x:  %08x  %-8s %s, %#x\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Imm()<<12)
	case rj32.J:
		fmt.Fprintf(s.In.Log, "%08x:  %08x  %-8s %s, %#x\n", s.In.PC, s.In.IR, ir.Opcode(), ir.Rd(), ir.Imm())
	}

	s.Out.PC = s.In.PC

	opcode := ir.Opcode()

	// where in the opcode map is it
	row := opcode >> 3
	col := opcode & 0b111

	s.AluOp = rj32.Add
	s.Sub = false

	if (row & 0b11) == 0b00 {
		s.AluOp = col
	}

	if (row&0b11) == 0b01 && (col == 0b110 || opcode == rj32.Sub) {
		s.AluOp = col
		s.Sub = true
	}

	if row == 0b110 {
		s.Sub = true
	}

	s.Out.RegWen = true

	if (row & 0b11) == 0b10 {
		s.Out.RegWen = false
	}

	if (opcode & 0b111100) == 0b111100 {
		s.Out.RegWen = false
	}

	s.Out.Rd = ir.Rd()

	if s.Out.Rd == 0 {
		s.Out.RegWen = false
	}

	fmt.Fprintf(s.In.Log, "  de: AluOp: %s, Sub: %v, RegWen: %v\n", s.Out.AluOp, s.Out.Sub, s.Out.RegWen)

	// read happens in the second half
	s.Out.L = s.Regs[ir.Rs1()]
	s.Out.R = s.Regs[ir.Rs2()]

	if (row & 0b100) == 0b100 {
		s.Out.R = uint32(ir.Imm())
	}

	fmt.Fprintf(s.In.Log, "  de: L: 0x%08x (%d)   R: 0x%08x (%d)\n", s.Out.L, s.Out.L, s.Out.R, s.Out.R)
}
