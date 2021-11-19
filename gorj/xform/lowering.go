package xform

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func AddPhiCopies(val *ir.Value) int {
	if val.Op != op.Phi {
		return 0
	}

	changes := 0

	for i := 0; i < val.NumArgs(); i++ {
		src := val.Arg(i)

		if src.Op == op.PhiCopy {
			continue
		}

		copy := val.Func().NewValue(op.PhiCopy, src.Type, src)
		pred := val.Block().Pred(i)

		pred.InsertInstr(-1, copy)
		val.ReplaceArg(i, copy)
		changes++
	}

	return changes
}

var _ = addToPass(Lowering, AddPhiCopies)
