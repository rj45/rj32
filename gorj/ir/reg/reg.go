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

func (reg Reg) StackSlot() int {
	return int(reg >> 16)
}

func StackSlot(slot int) Reg {
	return Reg(slot << 16)
}
