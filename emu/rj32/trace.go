package rj32

import "fmt"

// String returns the disassembled instruction as a string
func (ir Inst) String() string {
	switch ir.Fmt() {
	case FmtRR:
		return fmt.Sprintf("%-5s r%d, r%d", ir.Op(), ir.Rd(), ir.Rs())

	case FmtI11:
		return fmt.Sprintf("%-5s %d", ir.Op(), signExtend(ir.Imm(), 13))

	case FmtI12:
		return fmt.Sprintf("%-5s %d", ir.Op(), signExtend(ir.Imm(), 13))

	case FmtRI6, FmtRI8:
		return fmt.Sprintf("%-5s r%d, %d", ir.Op(), ir.Rd(), signExtend(ir.Imm(), 13))

	case FmtLS:
		if ir.Op() == Load || ir.Op() == Loadb {
			return fmt.Sprintf("%-5s r%d, [r%d, %d]", ir.Op(), ir.Rd(), ir.Rs(), ir.Imm())
		}
		return fmt.Sprintf("%-5s [r%d, %d], r%d", ir.Op(), ir.Rs(), ir.Imm(), ir.Rd())

	default:
		panic("not impl")
	}
}

// PreTrace returns a pre-execution debug string
func (ir Inst) PreTrace(cpu *CPU) string {
	switch ir.Fmt() {
	case FmtRR:
		return fmt.Sprintf("r%d:%d r%d:%d", ir.Rd(), cpu.Reg[ir.Rd()], ir.Rs(), cpu.Reg[ir.Rs()])

	case FmtI11:
		return fmt.Sprintf("pc:%04x rsval:%d", cpu.PC, cpu.rsval(ir))

	case FmtI12:
		return fmt.Sprintf("rsval:%d", cpu.rsval(ir))

	case FmtRI6, FmtRI8:
		return fmt.Sprintf("r%d:%d rsval:%d", ir.Rd(), cpu.Reg[ir.Rd()], cpu.rsval(ir))

	case FmtLS:
		return fmt.Sprintf("r%d:%d off:%d", ir.Rs(), cpu.Reg[ir.Rs()], cpu.off(ir.Imm()))

	default:
		panic("not impl")
	}
}

// PostTrace returns a post execution debug string
func (ir Inst) PostTrace(cpu *CPU) string {
	switch ir.Fmt() {
	case FmtRR, FmtRI6, FmtRI8:
		if ir.Op() >= IfEq {
			return fmt.Sprintf("  skip <- %v", cpu.Skip)
		}
		return fmt.Sprintf("  r%d <- %d", ir.Rd(), cpu.Reg[ir.Rd()])

	case FmtI11:
		return fmt.Sprintf("  pc <- %04x", cpu.PC)

	case FmtI12:
		return fmt.Sprintf("  imm <- %d", cpu.Imm)

	case FmtLS:
		if ir.Op() == Load || ir.Op() == Loadb {
			return fmt.Sprintf("  r%d <- %d", ir.Rd(), cpu.Reg[ir.Rd()])
		}
		return fmt.Sprintf("  mem <- %d", cpu.rsval(ir))

	default:
		panic("not impl")
	}
}
