package xform

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func dePhi(val *ir.Value) int {
	if val.Op != op.Phi {
		return 0
	}

	for i, src := range val.Args {
		if src.Op.IsConst() || val.Reg != src.Reg {
			// todo might actually need to be a swap instead
			var pred *ir.Block
			for _, ref := range val.Block.Preds {
				if ref.Index == i {
					pred = ref.Block
				}
			}
			pred.InsertInstr(-1, pred.Func.NewValue(ir.Value{
				Reg:  val.Reg,
				Op:   op.Copy,
				Args: []*ir.Value{src},
			}))
		}
	}

	return 0
}

var _ = addToPass(LastPass, dePhi)
