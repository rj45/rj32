package rj32

import (
	"io"
	"os"

	"github.com/rj45/rj32/emu/data"
)

func (cpu *CPU) LoadProgram(filename string) error {
	var buf []byte
	var err error
	if filename == "-" {
		buf, err = io.ReadAll(os.Stdin)
	} else {
		buf, err = os.ReadFile(filename)
	}
	if err != nil {
		return err
	}
	data.Load(16, buf, func(a int, val uint64) {
		cpu.Prog[a] = decodeInst(uint16(val))
	})

	return nil
}

func decodeInst(ir uint16) Inst {
	inst := Inst(0)

	fmt := Fmt(ir & 3)
	if fmt == FmtExt {
		if ir&4 == 0 {
			fmt = FmtRI8
		} else {
			fmt = FmtI11
		}
	}

	inst = inst.SetFmt(fmt)

	switch fmt {
	case FmtRR:
		i := InstRR(ir)
		inst = inst.
			SetRd(i.Rd()).
			SetRs(i.Rs()).
			SetOp(i.Op())
	case FmtLS:
		i := InstLS(ir)
		inst = inst.
			SetRd(i.Rd()).
			SetRs(i.Rs()).
			SetImm(i.Imm()).
			SetOp(i.Op() | 0b01100)
	case FmtRI6:
		i := InstRI6(ir)
		inst = inst.
			SetRd(i.Rd()).
			SetImm(signExtend(i.Imm(), 6)).
			SetOp(i.Op() | 0b10000)
	case FmtRI8:
		i := InstRI8(ir)
		inst = inst.
			SetRd(i.Rd()).
			SetImm(signExtend(i.Imm(), 8)).
			SetOp(i.Op() | 0b00110)
	case FmtI11:
		i := InstI11(ir)
		inst = inst.
			SetImm(signExtend(i.Imm(), 11)).
			SetOp(i.Op() | 0b01000)
	default:
		panic("unknown fmt")
	}

	return inst
}
