package reg

//go:generate enumer -type=Reg -transform lower

type Reg uint

const (
	None Reg = 0

	RA Reg = 1 << iota
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

func (reg Reg) IsStackSlot() bool {
	return reg > 0xFFFF
}

func (reg Reg) IsMask() bool {
	return reg <= 0xFFFF && !reg.IsAReg()
}

func (reg Reg) IsArgReg() bool {
	return reg == A1 || reg == A2
}

func (reg Reg) IsSpecialReg() bool {
	return reg == GP || reg == SP || reg == RA
}

func (reg Reg) IsSavedReg() bool {
	for _, saved := range SavedRegs {
		if reg == saved {
			return true
		}
	}
	return false
}

func (reg Reg) CanAffinity() bool {
	return !reg.IsSpecialReg()
}

func (reg Reg) StackSlot() int {
	return int(reg>>16) - 1
}

func StackSlot(slot int) Reg {
	return Reg((slot + 1) << 16)
}

var SavedRegs = []Reg{T5, S0, S1, S2, S3}
var RevSavedRegs = []Reg{S3, S2, S1, S0, T5}
var TempRegs = []Reg{T0, T1, T2, T3, T4}
var ArgRegs = []Reg{A0, A1, A2}
