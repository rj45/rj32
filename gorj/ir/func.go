package ir

import (
	"fmt"
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Func struct {
	Name string
	Type *types.Signature

	Mod *Package

	NumCalls int

	blocks     []*Block
	blockStore []Block

	values  []Value
	Consts  []*Value
	Params  []*Value
	Globals []*Value

	blockID idAlloc
	instrID idAlloc

	SpillSlots int
	ArgSlots   int
}

func (fn *Func) NextBlockID() ID {
	return fn.blockID.next()
}

func (fn *Func) BlockIDCount() int {
	return fn.blockID.count()
}

func (fn *Func) InstrIDCount() int {
	return fn.instrID.count()
}

func (fn *Func) String() string {
	return fn.Name
}

func (fn *Func) NewValue(op op.Op, typ types.Type, args ...*Value) *Value {
	fn.values = append(fn.values, Value{})
	val := &fn.values[len(fn.values)-1]
	val.id = fn.instrID.next()
	val.Op = op
	val.Type = typ

	for _, arg := range args {
		val.InsertArg(-1, arg)
	}

	if val.ID() != ID(len(fn.values)-1) {
		log.Panicln("value leak:", val.ID(), len(fn.values))
	}

	return val
}

func (fn *Func) NewRegValue(op op.Op, typ types.Type, reg reg.Reg, args ...*Value) *Value {
	val := fn.NewValue(op, typ, args...)
	val.Reg = reg
	return val
}

func (fn *Func) Blocks() []*Block {
	return fn.blocks
}

func (fn *Func) NewBlock(blk Block) *Block {
	blk.id = fn.NextBlockID()
	blk.fn = fn
	fn.blockStore = append(fn.blockStore, blk)
	return &fn.blockStore[len(fn.blockStore)-1]
}

func (fn *Func) InsertBlock(i int, blk *Block) {
	blk.fn = fn
	if i < 0 || i >= len(fn.blocks) {
		fn.blocks = append(fn.blocks, blk)
		return
	}

	fn.blocks = append(fn.blocks[:i+1], fn.blocks[i:]...)
	fn.blocks[i] = blk
}

func (fn *Func) RemoveBlock(i int) {
	blk := fn.blocks[i]

	if len(blk.preds) == 1 && len(blk.succs) == 1 {
		replPred := blk.preds[0]
		replSucc := blk.succs[0]

		for j, pred := range replSucc.preds {
			if pred == blk {
				replSucc.preds[j] = replPred
			}
		}

		for j, succ := range replPred.succs {
			if succ == blk {
				replPred.succs[j] = replSucc
			}
		}
	} else {
		panic("can't remove block")
	}

	fn.blocks = append(fn.blocks[:i], fn.blocks[i+1:]...)

}

func (fn *Func) ValueForID(id ID) *Value {
	return &fn.values[id]
}

func (fn *Func) BlockForID(id ID) *Block {
	return &fn.blockStore[id]
}

func (fn *Func) Const(typ types.Type, val constant.Value) *Value {
	for _, c := range fn.Consts {
		if types.Identical(c.Type, typ) && c.Value != nil && c.Value.ExactString() == val.ExactString() {
			return c
		}
	}

	con := fn.NewValue(op.Const, typ)
	con.Value = val
	fn.Consts = append(fn.Consts, con)
	return con
}

func (fn *Func) IntConst(val int64) *Value {
	inttype := types.Typ[types.Int]
	for _, c := range fn.Consts {
		if types.Identical(c.Type, inttype) && c.Value != nil {
			if v, ok := constant.Int64Val(c.Value); ok && v == val {
				return c
			}
		}
	}
	return fn.Const(inttype, constant.MakeInt64(val))
}

func (fn *Func) FixedReg(reg reg.Reg) *Value {
	for _, c := range fn.Consts {
		if c.Value == nil && c.Reg == reg {
			return c
		}
	}

	con := fn.NewValue(op.Reg, types.Typ[types.Int])
	con.Reg = reg
	fn.Consts = append(fn.Consts, con)
	return con
}

func (fn *Func) LongString() string {
	str := fmt.Sprintf("%s: ", fn.Name)

	typ := fmt.Sprintf("; %v", fn.Type)

	max := 40
	for (len(str)+len(typ)+max) > 80 && max > 0 {
		max--
	}

	for len(str) < max {
		str += " "
	}

	str += typ
	str += "\n"

	for _, blk := range fn.Blocks() {
		str += blk.LongString()
	}

	return str
}
