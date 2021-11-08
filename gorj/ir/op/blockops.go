package op

//go:generate enumer -type=BlockOp -transform title-lower

type BlockOp int

const (
	BadBlock BlockOp = iota
	Jump
	If
	Return
)
