package op

//go:generate enumer -type=Op -transform title-lower -json -text

type Def struct {
	Op      Op
	Asm     string
	Sink    bool
	Compare bool
}

type Op int

func (op Op) Asm() string {
	return op.Def().Asm
}

func (op Op) IsCompare() bool {
	return op.Def().Compare
}

func (op Op) IsSink() bool {
	return op.Def().Sink
}

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
	{Op: Call, Asm: "call"},
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
	{Op: Add, Asm: "add"},
	{Op: Sub, Asm: "sub"},
	{Op: Mul},
	{Op: Div},
	{Op: Rem},
	{Op: And, Asm: "and"},
	{Op: Or, Asm: "or"},
	{Op: Xor, Asm: "xor"},
	{Op: ShiftLeft, Asm: "shl"},
	{Op: ShiftRight, Asm: "shr"},
	{Op: AndNot},
	{Op: Equal, Compare: true},
	{Op: NotEqual, Compare: true},
	{Op: Less, Compare: true},
	{Op: LessEqual, Compare: true},
	{Op: Greater, Compare: true},
	{Op: GreaterEqual, Compare: true},
	{Op: Not},
	{Op: Negate, Asm: "neg"},
	{Op: Load},
	{Op: Invert, Asm: "not"},
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