package op

//go:generate enumer -type=Op -transform title-lower -json -text

type Def struct {
	Op      Op
	Sink    bool
	Compare bool
}

type Op int

const (
	Invalid Op = iota
	Builtin
	Call
	ChangeInterface
	ChangeType
	Const
	Convert
	Extract
	Field
	FieldAddr
	FreeVar
	Func
	Global
	Index
	IndexAddr
	Local
	Lookup
	MakeInterface
	MakeSlice
	Next
	New
	Panic
	Parameter
	Phi
	Range
	Slice
	SliceToArrayPointer
	Store
	TypeAssert
	Add
	Sub
	Mul
	Div
	Rem
	And
	Or
	Xor
	ShiftLeft
	ShiftRight
	AndNot
	Equal
	NotEqual
	Less
	LessEqual
	Greater
	GreaterEqual
	Not
	Negate
	Load
	Invert
)

var opDefs = []Def{
	{Op: Invalid},
	{Op: Builtin},
	{Op: Call},
	{Op: ChangeInterface},
	{Op: ChangeType},
	{Op: Const},
	{Op: Convert},
	{Op: Extract},
	{Op: Field},
	{Op: FieldAddr},
	{Op: FreeVar},
	{Op: Func},
	{Op: Global},
	{Op: Index},
	{Op: IndexAddr},
	{Op: Local},
	{Op: Lookup},
	{Op: MakeInterface},
	{Op: MakeSlice},
	{Op: Next},
	{Op: New},
	{Op: Panic},
	{Op: Parameter},
	{Op: Phi},
	{Op: Range},
	{Op: Slice},
	{Op: SliceToArrayPointer},
	{Op: Store, Sink: true},
	{Op: TypeAssert},
	{Op: Add},
	{Op: Sub},
	{Op: Mul},
	{Op: Div},
	{Op: Rem},
	{Op: And},
	{Op: Or},
	{Op: Xor},
	{Op: ShiftLeft},
	{Op: ShiftRight},
	{Op: AndNot},
	{Op: Equal, Compare: true},
	{Op: NotEqual, Compare: true},
	{Op: Less, Compare: true},
	{Op: LessEqual, Compare: true},
	{Op: Greater, Compare: true},
	{Op: GreaterEqual, Compare: true},
	{Op: Not},
	{Op: Negate},
	{Op: Load},
	{Op: Invert},
}

// sort opDefs so we don't have to worry about that
func init() {
	var newdefs []Def
	maxOp := Invalid
	for _, op := range opDefs {
		if op.Op > maxOp {
			maxOp = op.Op
		}
	}
	newdefs = make([]Def, maxOp+1)
	for _, op := range opDefs {
		newdefs[op.Op] = op
	}
	opDefs = newdefs

	for _, op := range OpValues() {
		if newdefs[op].Op != op {
			panic("missing OpDef for " + op.String())
		}
	}
}

func (op Op) Def() *Def {
	return &opDefs[op]
}
