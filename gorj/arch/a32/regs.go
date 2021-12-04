package a32

import "github.com/rj45/rj32/gorj/ir/reg"

//go:generate enumer -type=Reg -transform lower

type Reg uint

const (
	Zero Reg = iota

	RA // R1
	BP // R2
	SP // R3

	A0 // R4
	A1 // R5
	A2 // R6
	A3 // R7
	A4 // R8
	A5 // R9
	A6 // R10
	A7 // R11

	T0 // R12
	T1 // R13
	T2 // R14
	T3 // R15
	T4 // R16
	T5 // R17
	T6 // R18
	T7 // R19
	T8 // R20
	T9 // R21

	S0 // R22
	S1 // R23
	S2 // R24
	S3 // R25
	S4 // R26
	S5 // R27
	S6 // R28
	S7 // R29
	S8 // R30
	S9 // R31
)

var savedRegs = []Reg{S0, S1, S2, S3, S4, S5, S6, S7, S8, S9}
var tempRegs = []Reg{T0, T1, T2, T3, T4, T5, T6, T7, T8, T9}
var argRegs = []Reg{A0, A1, A2, A3, A4, A5, A6, A7}

func (cpuArch) RegNames() []string {
	return RegStrings()
}

func regList(regs []Reg) []reg.Reg {
	ret := make([]reg.Reg, len(regs))
	for i := range regs {
		ret[i] = reg.FromRegNum(int(regs[i]))
	}
	return ret
}

func (cpuArch) SavedRegs() []reg.Reg {
	return regList(savedRegs)
}

func (cpuArch) TempRegs() []reg.Reg {
	return regList(tempRegs)
}

func (cpuArch) ArgRegs() []reg.Reg {
	return regList(argRegs)
}

func (cpuArch) SpecialRegs() map[string]reg.Reg {
	return map[string]reg.Reg{
		"SP": reg.FromRegNum(int(SP)),
		"FP": reg.FromRegNum(int(BP)),
		"GP": reg.FromRegNum(int(Zero)),
		"RA": reg.FromRegNum(int(RA)),
	}
}
