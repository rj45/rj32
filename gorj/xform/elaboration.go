package xform

import (
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
	"github.com/rj45/rj32/gorj/sizes"
)

func calls(val *ir.Value) int {
	if val.Op != op.Call {
		return 0
	}

	changes := 0

	fnType := val.Args[0].Type.(*types.Signature)

	// TODO: handle multiple return values

	if fnType.Results().Len() == 1 && val.Reg != reg.A0 {
		changes++
		val.Reg = reg.A0
	}

	if len(val.Args) > 1 {
		if val.Args[1].Reg != reg.A1 {
			changes++
			val.Args[1] = val.Block.InsertCopy(val.Index, val.Args[1], reg.A1)
		}

		if len(val.Args) > 2 && val.Args[2].Reg != reg.A2 {
			changes++
			val.Args[2] = val.Block.InsertCopy(val.Index, val.Args[2], reg.A2)
		}

		slots := len(val.Args) - 3
		if slots > 0 {
			if val.Block.Func.ArgSlots < slots {
				val.Block.Func.ArgSlots = slots
			}

			for i := 0; i < slots; i++ {
				if val.Args[i+3].Reg != reg.StackSlot(i) {
					changes++
					val.Args[i+3] = val.Block.InsertCopy(val.Index, val.Args[i+3], reg.StackSlot(i))
				}
			}
		}
	}

	return changes
}

var _ = addToPass(Elaboration, calls)

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

	mulval := val.Block.Func.NewValue(ir.Value{
		Op:   op.Mul,
		Args: []*ir.Value{val.Args[1], sizeval},
		Type: types.Typ[types.Int],
	})
	val.Block.InsertInstr(val.Index, mulval)

	val.Op = op.Add
	val.Args[1] = mulval

	return 1
}

var _ = addToPass(Elaboration, indexAddrs)

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
		ir.SubstituteValue(val, val.Args[0])
		val.Block.RemoveInstr(val)
		return 1
	}

	val.Op = op.Add
	val.Value = nil
	val.Args = append(val.Args, val.Block.Func.Const(val.Type, constant.MakeInt64(offset)))

	return 1
}

var _ = addToPass(Elaboration, fieldAddrs)

func gpAdjustLoadStores(val *ir.Value) int {
	if val.Op != op.Load && val.Op != op.Store {
		return 0
	}

	if val.Args[0].Op == op.Global || val.Args[0].Op == op.Const {
		if len(val.Args) == 1 {
			val.Args = []*ir.Value{val.Block.Func.FixedReg(reg.GP), val.Args[0]}
			return 1
		}
		return 0
	}

	if val.Args[0].Op == op.Add && val.Args[0].Args[0].Reg == reg.GP {
		return 0
	}

	addval := val.Block.Func.NewValue(ir.Value{
		Op:   op.Add,
		Args: []*ir.Value{val.Block.Func.FixedReg(reg.GP), val.Args[0]},
		Type: types.Typ[types.Int],
	})
	val.Block.InsertInstr(val.Index, addval)

	val.Args[0] = addval

	return 1
}

var _ = addToPass(Elaboration, gpAdjustLoadStores)

func fixupConverts(val *ir.Value) int {
	if val.Op != op.Convert {
		return 0
	}

	if sizes.Sizeof(val.Args[0].Type) != sizes.Sizeof(val.Type) {
		log.Fatalf("Unable to convert %#v to %#v", val.Args[0].Type, val.Type)
	}

	ir.SubstituteValue(val, val.Args[0])
	val.Block.RemoveInstr(val)

	return 1
}

var _ = addToPass(Elaboration, fixupConverts)
