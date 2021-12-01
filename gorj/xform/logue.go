package xform

import (
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
	"github.com/rj45/rj32/gorj/sizes"
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
	countLocalSize(fn)
	framesize := int64(len(saved) + fn.SpillSlots + fn.ArgSlots + fn.LocalSlots)

	prologue(saved, framesize, fn)
	epilogue(saved, framesize, fn)
}

func countLocalSize(fn *ir.Func) {
	entry := fn.Blocks()[0]
	for i := 0; i < entry.NumInstrs(); i++ {
		val := entry.Instr(i)
		if val.Op != op.Local {
			continue
		}

		fn.LocalSlots += int(sizes.Sizeof(val.Type))
	}
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

	convertParams(entry, fn, len(saved))

	if framesize == 0 {
		return
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

func convertParams(entry *ir.Block, fn *ir.Func, saveSlots int) {
	localIndex := 0
	for i := 0; i < entry.NumInstrs(); i++ {
		val := entry.Instr(i)
		if val.Op == op.Parameter {
			index := -1
			for j, p := range fn.Params {
				if p == val {
					index = j
				}
			}
			if index < 0 {
				log.Panicln("could not find parameter", val)
			}

			switch index {
			case 0:
				val.Op = op.Copy
				val.InsertArg(-1, fn.FixedReg(reg.ArgRegs[1]))
			case 1:
				val.Op = op.Copy
				val.InsertArg(-1, fn.FixedReg(reg.ArgRegs[2]))
			default:
				val.Op = op.Load
				val.InsertArg(-1, fn.FixedReg(reg.SP))
				val.InsertArg(-1, fn.IntConst(int64(fn.ArgSlots+fn.SpillSlots+fn.LocalSlots+saveSlots+(index-2))))
			}
		}

		if val.Op == op.Local {
			offset := fn.ArgSlots + fn.SpillSlots + saveSlots + localIndex
			localIndex += int(sizes.Sizeof(val.Type))

			oval := ir.BuildBefore(val).
				Op(op.Copy, val.Type, offset).PrevVal()
			oval.Reg = val.Reg
			ir.BuildReplacement(val).Op(op.Add, val.Type, oval, reg.SP)
		}
	}
}

func epilogue(saved []reg.Reg, framesize int64, fn *ir.Func) {
	if framesize == 0 {
		return
	}

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
