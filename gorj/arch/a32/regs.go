package a32

import "github.com/rj45/rj32/gorj/ir/reg"

//go:generate enumer -type=Reg -transform lower

type Reg uint

const (
	Zero Reg = iota
	RA
	BP
	SP

	A0
	A1
	A2
	A3
	A4
	A5
	A6
	A7

	T0
	T1
	T2
	T3
	T4
	T5
	T6
	T7
	T8
	T9

	S0
	S1
	S2
	S3
	S4
	S5
	S6
	S7
	S8
	S9
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
		"GP": reg.FromRegNum(int(Zero)),
		"RA": reg.FromRegNum(int(RA)),
	}
}
