package xform

import (
	"go/types"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
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
