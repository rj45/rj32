package rj32

import "fmt"

type CPU struct {
	Cycles uint64

	// Registers
	Reg [16]int

	// Program Counter
	PC int

	// Immediate register
	Imm int

	// Pre-decoded program memory
	Prog [8192]Inst

	Halt, Error bool

	Trace bool
}

func (ir Inst) String() string {
	switch ir.Fmt() {
	case FmtRR:
		return fmt.Sprintf("%-5s r%d, r%d", ir.Op(), ir.Rd(), ir.Rs())

	case FmtI11:
		return fmt.Sprintf("%-5s %d", ir.Op(), signExtend(ir.Imm(), 12))

	case FmtRI6, FmtRI8:
		return fmt.Sprintf("%-5s r%d, %d", ir.Op(), ir.Rd(), signExtend(ir.Imm(), 12))

	case FmtLS:
		if ir.Op() == Load || ir.Op() == Loadb {
			return fmt.Sprintf("%-5s r%d, [r%d, %d]", ir.Op(), ir.Rd(), ir.Rs(), ir.Imm())
		}
		return fmt.Sprintf("%-5s [r%d, %d], r%d", ir.Op(), ir.Rs(), ir.Imm(), ir.Rd())

	default:
		panic("not impl")
	}
}

func (ir Inst) PreTrace(cpu *CPU) string {
	switch ir.Fmt() {
	case FmtRR:
		return fmt.Sprintf("r%d:%d r%d:%d", ir.Rd(), cpu.Reg[ir.Rd()], ir.Rs(), cpu.Reg[ir.Rs()])

	case FmtI11:
		return fmt.Sprintf("pc:%04x rsval:%d", cpu.PC, cpu.rsval(ir))

	case FmtRI6, FmtRI8:
		return fmt.Sprintf("r%d:%d rsval:%d", ir.Rd(), cpu.Reg[ir.Rd()], cpu.rsval(ir))

	case FmtLS:
		return fmt.Sprintf("r%d:%d off:%d", ir.Rs(), cpu.Reg[ir.Rs()], cpu.off(ir.Imm()))

	default:
		panic("not impl")
	}
}

func (ir Inst) PostTrace(cpu *CPU) string {
	switch ir.Fmt() {
	case FmtRR, FmtRI6, FmtRI8:
		return fmt.Sprintf("  r%d <- %d", ir.Rd(), cpu.Reg[ir.Rd()])

	case FmtI11:
		return fmt.Sprintf("  pc <- %04x", cpu.PC)

	case FmtLS:
		if ir.Op() == Load || ir.Op() == Loadb {
			return fmt.Sprintf("  r%d <- %d", ir.Rd(), cpu.Reg[ir.Rd()])
		}
		return fmt.Sprintf("  mem <- %d", cpu.rsval(ir))

	default:
		panic("not impl")
	}
}

// Run will run up to either the next IO request or
// the number of cycles has passed
func (cpu *CPU) Run(bus Bus, cycles int) Bus {
	for i := 0; i < cycles; i++ {
		cpu.Cycles++
		ir := cpu.Prog[cpu.PC]

		if cpu.Trace {
			fmt.Printf("%04x: %-15s %s\n", cpu.PC, ir, ir.PreTrace(cpu))
		}

		switch ir.Op() {
		case Nop:
			// do nothing

		case Halt:
			cpu.Halt = true
			return bus

		case Error:
			cpu.Error = true
			return bus

		case Rcsr:
			cpu.PC = cpu.Reg[ir.Rd()]
			if cpu.Trace {
				fmt.Printf("  temp jump, PC <- %04x\n", cpu.PC+1)
			}

		case Move:
			cpu.Reg[ir.Rd()] = cpu.rsval(ir)

		case Imm:
			cpu.Imm = cpu.rsval(ir)
			cpu.PC++
			return bus

		case Call:
			cpu.Reg[0] = cpu.PC
			if ir.Fmt() == FmtRR {
				cpu.PC = cpu.Reg[ir.Rd()]
			} else {
				cpu.PC += cpu.imm(ir.Imm())
			}

		case Jump:
			if ir.Fmt() == FmtRR {
				cpu.PC = cpu.Reg[ir.Rd()]
			} else {
				cpu.PC += cpu.imm(ir.Imm())
			}

		case Load:
			if !bus.Ack() {
				address := (cpu.Reg[ir.Rs()] + cpu.off(ir.Imm())) & 0xffff
				return bus.
					SetReq(true).
					SetWE(false).
					SetAddress(address)
			}
			cpu.Reg[ir.Rd()] = bus.Data()
			bus = bus.SetReq(false)

		case Store:
			if !bus.Ack() {
				address := (cpu.Reg[ir.Rs()] + cpu.off(ir.Imm())) & 0xffff
				return bus.
					SetReq(true).
					SetWE(true).
					SetAddress(address).
					SetData(cpu.Reg[ir.Rd()])
			}
			bus = bus.SetReq(false).SetWE(false)

		case Add:
			cpu.Reg[ir.Rd()] = cpu.Reg[ir.Rd()] + cpu.rsval(ir)

		case Sub:
			cpu.Reg[ir.Rd()] = cpu.Reg[ir.Rd()] - cpu.rsval(ir)

		case Xor:
			cpu.Reg[ir.Rd()] = cpu.Reg[ir.Rd()] ^ cpu.rsval(ir)

		case And:
			cpu.Reg[ir.Rd()] = cpu.Reg[ir.Rd()] & cpu.rsval(ir)

		case Or:
			cpu.Reg[ir.Rd()] = cpu.Reg[ir.Rd()] | cpu.rsval(ir)

		case Shl:
			cpu.Reg[ir.Rd()] = cpu.Reg[ir.Rd()] << (cpu.rsval(ir) & 0xf)

		case Shr:
			cpu.Reg[ir.Rd()] = (cpu.Reg[ir.Rd()] & 0xffff) >> (cpu.rsval(ir) & 0xf)

		case Asr:
			cpu.Reg[ir.Rd()] = signExtend(cpu.Reg[ir.Rd()]&0xffff, 16) >> (cpu.rsval(ir) & 0xf)

		case IfEq:
			if (cpu.Reg[ir.Rd()] & 0xffff) != (cpu.rsval(ir) & 0xffff) {
				cpu.PC++
			}

		case IfNe:
			if (cpu.Reg[ir.Rd()] & 0xffff) == (cpu.rsval(ir) & 0xffff) {
				cpu.PC++
			}

		case IfLt:
			l := signExtend(cpu.Reg[ir.Rd()]&0xffff, 16)
			r := signExtend(cpu.rsval(ir)&0xffff, 16)
			if l >= r {
				cpu.PC++
			}

		case IfGe:
			l := signExtend(cpu.Reg[ir.Rd()]&0xffff, 16)
			r := signExtend(cpu.rsval(ir)&0xffff, 16)
			if l < r {
				cpu.PC++
			}

		case IfHs:
			if (cpu.Reg[ir.Rd()] & 0xffff) < (cpu.rsval(ir) & 0xffff) {
				cpu.PC++
			}

		case IfLo:
			if (cpu.Reg[ir.Rd()] & 0xffff) >= (cpu.rsval(ir) & 0xffff) {
				cpu.PC++
			}

		default:
			panic("Op not yet implemented: " + ir.Op().String())
		}

		cpu.Imm = 0
		cpu.PC++

		if cpu.Trace {
			fmt.Println(ir.PostTrace(cpu))
		}
	}
	return bus
}

func (cpu *CPU) off(imm int) int {
	return (imm & 0b1111) | (cpu.Imm>>1)&0x7fff
}

func signExtend(val, bits int) int {
	m := 1 << (bits - 1)
	return (val ^ m) - m
}

func (cpu *CPU) imm(imm int) int {
	if cpu.Imm != 0 {
		return cpu.Imm | (imm & 0b11111)
	}
	return signExtend(imm, 12)
}

func (cpu *CPU) rsval(ir Inst) int {
	if ir.Fmt() == FmtRR {
		return cpu.Reg[ir.Rs()]
	}
	return cpu.imm(ir.Imm())
}

func DecodeInst(ir uint16) Inst {
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
