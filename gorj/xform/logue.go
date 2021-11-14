package xform

import (
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

/*
Here is what the call frames look like:

high mem addresses
+------------------------+
| caller saved reg 2     |  |
+------------------------+  |
| caller saved reg 1     |   > previous caller's frame
+------------------------+  |
| caller saved ra        |  |
+------------------------+  |
| stack arg 4            |  |
+------------------------+  |
| stack arg 3            |  |  <-- caller SP
+------------------------+ /
| saved reg 2            | \
+------------------------+  |
| saved reg 1            |  |
+------------------------+  |
| saved ra (if needed)   |  |
+------------------------+  |
| spill local 1          |  |
+------------------------+   > current callee's frame
| spill local 0          |  |
+------------------------+  |
| stack arg 4            |  |
+------------------------+  |
| stack arg 3            |  | <-- SP
+------------------------+ /
| saved reg 2            | \
+------------------------+  |
| saved reg 1            |   > beggings of next frame
+------------------------+  |
low mem addresses

Note that function parameters (arguments) are on the caller's frame,
in the area known as "ArgSlots".

So the order on the from the SP is:
  - ArgSlots for calls
	- SpillSlots for local variables on the stack
	- Saved registers
	- Params for the function incoming parameters

*/

func ProEpiLogue(usedRegs reg.Reg, fn *ir.Func) {
	saved := savedRegs(usedRegs, fn)
	framesize := int64(len(saved) + fn.SpillSlots + fn.ArgSlots)

	prologue(saved, framesize, fn)
	epilogue(saved, framesize, fn)
}

func prologue(saved []reg.Reg, framesize int64, fn *ir.Func) {
	// SP currently pointing at function parameters
	// in other words, the `ArgSlots` of the previous function
	entry := fn.Blocks()[0]
	sp := fn.FixedReg(reg.SP)
	index := 0

	if entry.NumPreds() > 0 {
		log.Fatalf("Entry cannot be jumped to or bad things!")
	}

	entry.InsertInstr(index, fn.NewRegValue(op.Sub, types.Typ[types.Int],
		reg.SP, sp,
		fn.IntConst(framesize)))
	index++

	for i, reg := range saved {
		offset := int64(i + fn.SpillSlots + fn.ArgSlots)
		entry.InsertInstr(index, fn.NewValue(op.Store, types.Typ[types.Int],
			sp,
			fn.IntConst(offset),
			fn.FixedReg(reg)))
		index++
	}
}

func epilogue(saved []reg.Reg, framesize int64, fn *ir.Func) {
	exit := fn.Blocks()[len(fn.Blocks())-1]
	sp := fn.FixedReg(reg.SP)

	for i, reg := range saved {
		offset := int64(i + fn.SpillSlots + fn.ArgSlots)
		exit.InsertInstr(-1, fn.NewRegValue(op.Load, types.Typ[types.Int],
			reg,
			sp,
			fn.IntConst(offset)))
	}

	exit.InsertInstr(-1, fn.NewRegValue(op.Add, types.Typ[types.Int],
		reg.SP,
		sp,
		fn.IntConst(framesize)))
}

func savedRegs(usedRegs reg.Reg, fn *ir.Func) []reg.Reg {
	var saved []reg.Reg

	if fn.NumCalls > 0 {
		saved = append(saved, reg.RA)
	}

	for _, reg := range reg.SavedRegs {
		if (usedRegs & reg) != 0 {
			saved = append(saved, reg)
		}
	}

	return saved
}