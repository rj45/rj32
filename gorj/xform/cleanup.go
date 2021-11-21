package xform

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

// func dePhi(val *ir.Value) int {
// 	if val.Op != op.Phi {
// 		return 0
// 	}

// 	for i := 0; i < val.NumArgs(); i++ {
// 		src := val.Arg(i)

// 		if src.Op != op.PhiCopy {
// 			log.Panicf("expecting phi %s to have arg %s be a PhiCopy", val, src)
// 		}

// 		if src.Reg != val.Reg {
// 			log.Panicf("expecting all args of phi %s to be assigned the same reg but saw %s:\n%s", val, src, val.LongString())
// 		}

// 		// not sure what to do here. The PhiCopies could be rearranged in a
// 		// non-conflicting order maybe?
// 	}

// 	return 0
// }

// var _ = addToPass(CleanUp, dePhi)

func deCopy(val *ir.Value) int {
	if val.Op != op.Copy && val.Op != op.PhiCopy {
		return 0
	}

	if val.Reg == val.Arg(0).Reg {
		val.ReplaceWith(val.Arg(0))
		return 1
	}

	return 0
}

var _ = addToPass(CleanUp, deCopy)

func EliminateEmptyBlocks(fn *ir.Func) {
	blks := fn.Blocks()
retry:
	for {
		for i, blk := range blks {
			if blk.NumInstrs() == 0 && blk.Op == op.Jump && blk.NumPreds() == 1 && blk.NumSuccs() == 1 {
				fn.RemoveBlock(i)
				continue retry
			}
		}
		break
	}
}
