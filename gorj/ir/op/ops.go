package op

//go:generate enumer -type=Op -transform title-lower

type Def struct {
	Op      Op
	Asm     string
	Sink    bool
	Compare bool
	Const   bool
	ClobArg bool
	Copy    bool
	Commute bool
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

func (op Op) IsConst() bool {
	return op.Def().Const
}

func (op Op) IsCopy() bool {
	return op.Def().Copy
}

func (op Op) IsCommutative() bool {
	return op.Def().Commute
}

func (op Op) ClobbersArg() bool {
	return op.Def().ClobArg
}

func (op Op) Opposite() Op {
	switch op {
	case Equal:
		return NotEqual
	case NotEqual:
		return Equal
	case Less:
		return GreaterEqual
	case LessEqual:
		return Greater
	case Greater:
		return LessEqual
	case GreaterEqual:
		return Less
	}
	return op
}

const (
	Invalid Op = iota
	Builtin
	Call
	ChangeInterface
	ChangeType
	Const
	Convert
	Copy
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
	Parameter
	Phi
	PhiCopy
	Range
	Reg
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
	{Op: Builtin, Const: true},
	{Op: Call, Asm: "call"},
	{Op: ChangeInterface},
	{Op: ChangeType},
	{Op: Const, Const: true},
	{Op: Convert},
	{Op: Copy, Asm: "move", Copy: true},
	{Op: Extract},
	{Op: Field},
	{Op: FieldAddr},
	{Op: FreeVar},
	{Op: Func, Const: true},
	{Op: Global, Const: true},
	{Op: Index},
	{Op: IndexAddr},
	{Op: Local},
	{Op: Lookup},
	{Op: MakeInterface},
	{Op: MakeSlice},
	{Op: Next},
	{Op: New},
	{Op: Parameter},
	{Op: Phi, Copy: true},
	{Op: PhiCopy, Asm: "move", Copy: true},
	{Op: Range},
	{Op: Reg, Copy: true},
	{Op: Slice},
	{Op: SliceToArrayPointer},
	{Op: Store, Sink: true},
	{Op: TypeAssert},
	{Op: Add, Asm: "add", ClobArg: true, Commute: true},
	{Op: Sub, Asm: "sub", ClobArg: true},
	{Op: Mul, Commute: true},
	{Op: Div},
	{Op: Rem},
	{Op: And, Asm: "and", ClobArg: true, Commute: true},
	{Op: Or, Asm: "or", ClobArg: true, Commute: true},
	{Op: Xor, Asm: "xor", ClobArg: true, Commute: true},
	{Op: ShiftLeft, Asm: "shl", ClobArg: true},
	{Op: ShiftRight, Asm: "shr", ClobArg: true},
	{Op: AndNot},
	{Op: Equal, Compare: true, Commute: true},
	{Op: NotEqual, Compare: true, Commute: true},
	{Op: Less, Compare: true},
	{Op: LessEqual, Compare: true},
	{Op: Greater, Compare: true},
	{Op: GreaterEqual, Compare: true},
	{Op: Not},
	{Op: Negate, Asm: "neg", ClobArg: true},
	{Op: Load},
	{Op: Invert, Asm: "not", ClobArg: true},
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
