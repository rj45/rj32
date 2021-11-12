package xform

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func dePhi(val *ir.Value) int {
	if val.Op != op.Phi {
		return 0
	}

	for i := 0; i < val.NumArgs(); i++ {
		src := val.Arg(i)
		if src.Op.IsConst() || val.Reg != src.Reg {
			// todo might actually need to be a swap instead
			pred := val.Block().Pred(i)
			pred.InsertCopy(-1, src, val.Reg)
		}
	}

	val.Remove()

	return 1
}

var _ = addToPass(LastPass, dePhi)

func deCopy(val *ir.Value) int {
	if val.Op != op.Copy {
		return 0
	}

	if val.Reg == val.Arg(0).Reg {
		val.Remove()
		return 1
	}

	return 0
}

var _ = addToPass(LastPass, deCopy)
