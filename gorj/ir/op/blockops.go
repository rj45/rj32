package op

//go:generate enumer -type=BlockOp -transform title-lower -json -text

type BlockOp int

const (
	BadBlock BlockOp = iota
	Jump
	If
	Return
)
