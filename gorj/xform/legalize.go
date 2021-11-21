package xform

import (
	"log"

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

	if val.Arg(0).Reg == val.Reg {
		return 0
	}

	if val.Arg(1).Reg == val.Reg && val.Op.IsCommutative() {
		// swap
		val.InsertArg(-1, val.RemoveArg(0))
	}

	if val.Arg(0).Op == op.Copy {
		return 0
	}

	copied := val.Arg(0)
	copy := val.Func().NewValue(op.Copy, copied.Type, copied)
	copy.Reg = val.Reg

	val.Block().InsertInstr(val.Index(), copy)
	val.ReplaceArg(0, copy)

	return 1
}

var _ = addToPass(Legalize, addCopiesForArgClobbers)

// reorderPhiCopies finds any phiCopies that aren't eliminated that
// happen to clobber a register in the same PhiCopy block.
// PhiCopies are supposed to happen in parallel -- all at the same time --
// in order to avoid restricting the register allocator
func reorderPhiCopies(val *ir.Value) int {
	if val.Op != op.PhiCopy {
		return 0
	}

	// scan through previous instructions looking for a PhiCopy writing to
	// the register being read
	blk := val.Block()
	for i := val.Index() - 1; i >= 0; i-- {
		prev := blk.Instr(i)
		if prev.Op != op.PhiCopy {
			break
		}
		read := val.Arg(0)
		if prev.Reg == read.Reg && read.NeedsReg() {
			// check if this is a swap
			if val.Reg == prev.Arg(0).Reg {
				log.Panicf("need to swap these: %s and %s", val.ShortString(), prev.ShortString())
			}

			// otherwise just swap the read and the write
			blk.SwapInstr(val, prev)
		}
	}

	return 0
}

var _ = addToPass(Legalize, reorderPhiCopies)
