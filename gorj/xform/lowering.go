package xform

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

// The first arg of many instructions must be assigned the same
// register as the destination. This is done by inserting copies
// at this stage, which can later be removed if the register is
// the same.
func addCopiesForArgClobbers(val *ir.Value) int {
	if !val.Op.ClobbersArg() {
		return 0
	}

	if val.NumArgs() < 1 {
		return 0
	}

	if val.Arg(0).Op == op.Copy {
		return 0
	}

	copied := val.Arg(0)
	copy := val.Func().NewValue(op.Copy, copied.Type, copied)

	val.Block().InsertInstr(val.Index(), copy)
	val.ReplaceArg(0, copy)

	return 1
}

var _ = addToPass(Lowering, addCopiesForArgClobbers)
