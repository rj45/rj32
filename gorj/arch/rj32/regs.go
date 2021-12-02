package rj32

//go:generate enumer -type=Reg -transform lower

type Reg uint

const (
	None Reg = iota

	RA
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
