package xform

import (
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
	"github.com/rj45/rj32/gorj/sizes"
)

// indexAddrs converts `IndexAddr` instructions into a `mul` and `add` instruction
// The `mul` is by a constant which can be optimized into shifts and adds by
// some other piece of code.
func indexAddrs(val *ir.Value) int {
	if val.Op != op.IndexAddr {
		return 0
	}

	elem := val.Type.(*types.Pointer).Elem()

	size := sizes.Sizeof(elem)

	sizeval := val.Block.Func.Const(types.Typ[types.Int], constant.MakeInt64(size))

	mulval := &ir.Value{
		ID:   val.Block.NextInstrID(),
		Op:   op.Mul,
		Args: []*ir.Value{val.Args[1], sizeval},
		Type: types.Typ[types.Int],
	}
	val.Block.InsertInstr(val.Index, mulval)

	val.Op = op.Add
	val.Args[1] = mulval

	return 1
}

var _ = addToPass(0, indexAddrs)

func fieldAddrs(val *ir.Value) int {
	if val.Op != op.FieldAddr {
		return 0
	}

	field, ok := constant.Int64Val(val.Value)
	if !ok {
		panic("expected int constant")
	}

	elem := val.Args[0].Type.(*types.Pointer).Elem()
	strct := elem.(*types.Struct)

	fields := sizes.Fieldsof(strct)
	offsets := sizes.Offsetsof(fields)
	offset := offsets[field]

	if offset == 0 {
		// would just be adding zero, so this instruction can just be removed
		val.Block.RemoveInstr(val)
		return 1
	}

	val.Op = op.Add
	val.Value = nil
	val.Args = append(val.Args, val.Block.Func.Const(val.Type, constant.MakeInt64(offset)))

	return 1
}

var _ = addToPass(0, fieldAddrs)

func gpAdjustLoadStores(val *ir.Value) int {
	if val.Op != op.Load && val.Op != op.Store {
		return 0
	}

	if val.Args[0].Op == op.Global || val.Args[0].Op == op.Const {
		return 0
	}

	if val.Args[0].Op == op.Add && val.Args[0].Args[0].Reg == reg.GP {
		return 0
	}

	addval := &ir.Value{
		ID:   val.Block.NextInstrID(),
		Op:   op.Add,
		Args: []*ir.Value{val.Block.Func.FixedReg(reg.GP), val.Args[0]},
		Type: types.Typ[types.Int],
	}
	val.Block.InsertInstr(val.Index, addval)

	val.Args[0] = addval

	return 1
}

var _ = addToPass(0, gpAdjustLoadStores)
