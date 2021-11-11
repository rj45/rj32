package logue

import (
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

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

func Prologue(usedRegs reg.Reg, fn *ir.Func) {
	saved := savedRegs(usedRegs, fn)

	// SP currently pointing at function parameters
	// in other words, the `ArgSlots` of the previous function
	framesize := int64(len(saved) + fn.SpillSlots + fn.ArgSlots)
	entry := fn.Blocks[0]
	sp := fn.FixedReg(reg.SP)
	index := 0

	entry.InsertInstr(index, fn.NewValue(ir.Value{
		Op:  op.Sub,
		Reg: reg.SP,
		Args: []*ir.Value{
			sp,
			fn.Const(types.Typ[types.Int], constant.MakeInt64(framesize)),
		},
	}))
	index++

	for i, reg := range saved {
		offset := int64(i + fn.SpillSlots + fn.ArgSlots)
		entry.InsertInstr(index, fn.NewValue(ir.Value{
			Op: op.Store,
			Args: []*ir.Value{
				sp,
				fn.Const(types.Typ[types.Int], constant.MakeInt64(offset)),
				fn.FixedReg(reg),
			},
		}))
		index++
	}
}

func Epilogue(usedRegs reg.Reg, fn *ir.Func) {
	saved := savedRegs(usedRegs, fn)

	framesize := int64(len(saved) + fn.SpillSlots + fn.ArgSlots)
	exit := fn.Blocks[len(fn.Blocks)-1]
	sp := fn.FixedReg(reg.SP)

	for i, reg := range saved {
		offset := int64(i + fn.SpillSlots + fn.ArgSlots)
		exit.InsertInstr(-1, fn.NewValue(ir.Value{
			Op:  op.Load,
			Reg: reg,
			Args: []*ir.Value{
				sp,
				fn.Const(types.Typ[types.Int], constant.MakeInt64(offset)),
			},
		}))
	}

	exit.InsertInstr(-1, fn.NewValue(ir.Value{
		Op:  op.Add,
		Reg: reg.SP,
		Args: []*ir.Value{
			sp,
			fn.Const(types.Typ[types.Int], constant.MakeInt64(framesize)),
		},
	}))
}
