package rj32

import "github.com/rj45/rj32/gorj/ir/reg"

//go:generate enumer -type=Reg -transform lower

type Reg uint

const (
	RA Reg = iota
	A0
	A1
	A2

	S0
	S1
	S2
	S3

	T0
	T1
	T2
	T3

	T4
	T5
	GP
	SP
)

var SavedRegs = []Reg{S0, S1, S2, S3, T5}
var TempRegs = []Reg{T0, T1, T2, T3, T4}
var ArgRegs = []Reg{A0, A1, A2}

func (Rj32) RegNames() []string {
	return RegStrings()
}

func regList(regs []Reg) []reg.Reg {
	ret := make([]reg.Reg, len(regs))
	for i := range regs {
		ret[i] = reg.FromRegNum(int(regs[i]))
	}
	return ret
}

func (Rj32) SavedRegs() []reg.Reg {
	return regList(SavedRegs)
}

func (Rj32) TempRegs() []reg.Reg {
	return regList(TempRegs)
}

func (Rj32) ArgRegs() []reg.Reg {
	return regList(ArgRegs)
}

func (Rj32) SpecialRegs() map[string]reg.Reg {
	return map[string]reg.Reg{
		"SP": reg.FromRegNum(int(SP)),
		"GP": reg.FromRegNum(int(GP)),
		"RA": reg.FromRegNum(int(RA)),
	}
}
