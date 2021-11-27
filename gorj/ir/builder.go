package ir

import (
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Builder struct {
	blk     *Block
	index   int
	prev    *Value
	replace bool
}

type prevVal struct{}

func BuildAt(blk *Block, index int) Builder {
	i := index
	if index < 0 {
		i = len(blk.instrs) - 1
	}
	var val *Value
	if len(blk.instrs) > 0 {
		val = blk.instrs[i]
	}
	return Builder{blk, index, val, false}
}

func BuildBefore(val *Value) Builder {
	return Builder{val.block, val.index, val, false}
}

func BuildAfter(val *Value) Builder {
	return Builder{val.block, val.index + 1, val, false}
}

func BuildReplacement(val *Value) Builder {
	return Builder{val.block, val.index, val, true}
}

func PrevBuildVal() interface{} {
	return prevVal{}
}

func (bd Builder) Op(op op.Op, typ types.Type, args ...interface{}) Builder {
	nargs := make([]*Value, len(args))

	fn := bd.blk.Func()
	for i, arg := range args {
		switch arg := arg.(type) {
		case *Value:
			nargs[i] = arg
		case prevVal:
			nargs[i] = bd.PrevVal()
		case int:
			nargs[i] = fn.IntConst(int64(arg))
		case int64:
			nargs[i] = fn.IntConst(arg)
		case reg.Reg:
			nargs[i] = fn.FixedReg(arg)
		}
	}

	if bd.replace {
		val := bd.blk.instrs[bd.index]
		val.Op = op
		val.Type = typ

		for i, arg := range nargs {
			if len(val.args) <= i {
				val.InsertArg(i, arg)
			} else if arg != val.args[i] {
				val.ReplaceArg(i, arg)
			}
		}
		for len(nargs) < len(val.args) {
			val.RemoveArg(len(val.args) - 1)
		}

		return Builder{bd.blk, bd.index, val, false}
	}

	val := fn.NewValue(op, typ, nargs...)
	bd.blk.InsertInstr(bd.index, val)

	if bd.index >= 0 {
		return Builder{bd.blk, val.index + 1, val, false}
	}
	return Builder{bd.blk, bd.index, val, false}
}

// PrevVal returns the previously inserted value
func (bd Builder) PrevVal() *Value {
	return bd.prev
}
