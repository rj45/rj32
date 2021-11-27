package rj32

import (
	"fmt"
	"io"
	"os"

	"github.com/rj45/rj32/emurj/data"
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
		if a >= 0x10000 {
			if cpu.Trace {
				fmt.Fprintf(os.Stderr, "data %04x: %04x\n", a-0x10000+0x8000, val)
			}
			// data bank starts at at 0x10000 in the file
			// but the data is actually stored at 0x8000
			// in the data address space
			bus := data.Bus(0).
				SetWE(true).
				SetAddress(a - 0x10000 + 0x8000).
				SetData(int(val))
			cpu.BusHandler.HandleBus(bus)
			return
		}
		if cpu.Trace && val != 0 {
			fmt.Fprintf(os.Stderr, "code %04x: %04x\n", a, val)
		}
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
		} else if ir&8 == 0 {
			fmt = FmtI11
		} else {
			fmt = FmtI12
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

	case FmtI12:
		i := InstI12(ir)
		inst = inst.
			SetImm(signExtend(i.Imm(), 12)).
			SetOp(Imm)

	default:
		panic("unknown fmt")
	}

	return inst
}
