package rj32

import (
	"fmt"
	"os"

	"github.com/rj45/rj32/emurj/data"
)

type CPU struct {
	Cycles uint64

	// Registers
	Reg [16]int

	// Program Counter
	PC int

	// Skip register
	Skip bool

	// Immediate register
	Imm       int
	ImmValid  bool
	ImmExpire bool

	// Pre-decoded program memory
	Prog [65536]Inst

	BusHandler data.BusHandler

	Halt, Error bool

	Trace bool
}

// Run will run up to either the next IO request or
// the number of cycles has passed
func (cpu *CPU) Run(cycles int) {
	endCycle := cpu.Cycles + uint64(cycles)
	for ; cpu.Cycles < endCycle; cpu.Cycles++ {
		ir := cpu.Prog[cpu.PC]

		if cpu.Trace {
			fmt.Fprintf(os.Stderr, "%04x: %-15s %s\n", cpu.PC, ir, ir.PreTrace(cpu))
		}

		switch ir.Op() {
		case Nop:
			// do nothing

		case Halt:
			cpu.Halt = true
			return

		case Error:
			cpu.Error = true
			return

		case Rcsr:
			cpu.PC = cpu.Reg[ir.Rd()]
			if cpu.Trace {
				fmt.Fprintf(os.Stderr, "  temp jump, PC <- %04x\n", cpu.PC+1)
			}

		case Move:
			cpu.Reg[ir.Rd()] = cpu.rsval(ir)

		case Imm, Imm2:
			cpu.Imm = cpu.rsval(ir) << 4
			cpu.ImmExpire = true
			cpu.ImmValid = true

		case Call:
			cpu.Reg[0] = cpu.PC
			if ir.Fmt() == FmtRR {
				cpu.PC = cpu.Reg[ir.Rs()]
			} else {
				cpu.PC += cpu.imm(ir.Imm())
			}

		case Jump:
			if ir.Fmt() == FmtRR {
				cpu.PC = cpu.Reg[ir.Rs()]
			} else {
				cpu.PC += cpu.imm(ir.Imm())
			}

		case Load:
			address := (cpu.Reg[ir.Rs()] + cpu.off(ir.Imm())) & 0xffff
			bus := data.Bus(0).
				SetWE(false).
				SetAddress(address)
			bus = cpu.BusHandler.HandleBus(bus)
			if !bus.Ack() {
				panic(fmt.Sprintf("hung waiting for bus read at address %04x", address))
			}
			cpu.Reg[ir.Rd()] = bus.Data()

			// todo: handle multiple cycle memory access
			cpu.Cycles++

		case Store:
			address := (cpu.Reg[ir.Rs()] + cpu.off(ir.Imm())) & 0xffff
			bus := data.Bus(0).
				SetWE(true).
				SetAddress(address).
				SetData(cpu.Reg[ir.Rd()])
			bus = cpu.BusHandler.HandleBus(bus)
			if !bus.Ack() {
				panic(fmt.Sprintf("hung waiting for bus write at address %04x", address))
			}

			// todo: handle multiple cycle memory access
			cpu.Cycles++

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
			cpu.Reg[ir.Rd()] = signExtend((cpu.Reg[ir.Rd()]<<(cpu.rsval(ir)&0xf))&0xffff, 16)

		case Shr:
			cpu.Reg[ir.Rd()] = (cpu.Reg[ir.Rd()] & 0xffff) >> (cpu.rsval(ir) & 0xf)

		case Asr:
			cpu.Reg[ir.Rd()] = signExtend(cpu.Reg[ir.Rd()]&0xffff, 16) >> (cpu.rsval(ir) & 0xf)

		case IfEq:
			if (cpu.Reg[ir.Rd()] & 0xffff) != (cpu.rsval(ir) & 0xffff) {
				cpu.Skip = true
			}

		case IfNe:
			if (cpu.Reg[ir.Rd()] & 0xffff) == (cpu.rsval(ir) & 0xffff) {
				cpu.Skip = true
			}

		case IfLt:
			l := signExtend(cpu.Reg[ir.Rd()]&0xffff, 16)
			r := signExtend(cpu.rsval(ir)&0xffff, 16)
			if l >= r {
				cpu.Skip = true
			}

		case IfGe:
			l := signExtend(cpu.Reg[ir.Rd()]&0xffff, 16)
			r := signExtend(cpu.rsval(ir)&0xffff, 16)
			if l < r {
				cpu.Skip = true
			}

		case IfUge:
			var a uint16 = uint16(cpu.Reg[ir.Rd()] & 0xffff)
			var b uint16 = uint16(cpu.rsval(ir) & 0xffff)

			if a < b {
				cpu.Skip = true
			}

		case IfUlt:
			var a uint16 = uint16(cpu.Reg[ir.Rd()] & 0xffff)
			var b uint16 = uint16(cpu.rsval(ir) & 0xffff)

			if a >= b {
				cpu.Skip = true
			}

		default:
			panic("Op not yet implemented: " + ir.Op().String())
		}

		if !cpu.ImmExpire {
			cpu.Imm = 0
			cpu.ImmValid = false
		}
		cpu.ImmExpire = false
		cpu.PC++

		if cpu.Trace {
			fmt.Fprintln(os.Stderr, ir.PostTrace(cpu))
		}

		if cpu.Skip {
			if cpu.Prog[cpu.PC].Op() == Imm {
				cpu.PC++
			}

			cpu.PC++
			cpu.Skip = false
		}
	}
}

func (cpu *CPU) off(imm int) int {
	return (imm & 0b1111) | cpu.Imm
}

func signExtend(val, bits int) int {
	m := 1 << (bits - 1)
	return (val ^ m) - m
}

func (cpu *CPU) imm(imm int) int {
	if cpu.ImmValid {
		return cpu.Imm | (imm & 0b1111)
	}
	return signExtend(imm, 13)
}

func (cpu *CPU) rsval(ir Inst) int {
	if ir.Fmt() == FmtRR {
		return cpu.Reg[ir.Rs()]
	}
	return cpu.imm(ir.Imm())
}
