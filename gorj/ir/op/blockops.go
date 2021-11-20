package op

//go:generate enumer -type=BlockOp -transform title-lower

type BlockOp int

const (
	BadBlock BlockOp = iota
	Jump
	If
	Return
	Panic

	IfEqual
	IfNotEqual
	IfLess
	IfLessEqual
	IfGreater
	IfGreaterEqual
)

func (op BlockOp) Compare() Op {
	switch op {
	case IfEqual:
		return Equal
	case IfNotEqual:
		return NotEqual
	case IfLess:
		return Less
	case IfLessEqual:
		return LessEqual
	case IfGreater:
		return Greater
	case IfGreaterEqual:
		return GreaterEqual
	}
	return Invalid
}
