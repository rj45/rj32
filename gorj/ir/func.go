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

	Mod *Module

	NumCalls int

	Blocks []*Block

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

func (fn *Func) NextInstrID() ID {
	return fn.instrID.next()
}

func (fn *Func) String() string {
	return fn.Name
}

func (fn *Func) NewValue(val Value) *Value {
	val.ID = fn.NextInstrID()
	if val.ID != ID(len(fn.values)) {
		// TODO: may be prudent to make ID private to avoid this
		log.Fatalln("value leak:", val.ID, len(fn.values))
	}
	fn.values = append(fn.values, val)
	return &fn.values[len(fn.values)-1]
}

func (fn *Func) ValueForID(id ID) *Value {
	return &fn.values[id]
}

func (fn *Func) Const(typ types.Type, val constant.Value) *Value {
	for _, c := range fn.Consts {
		if types.Identical(c.Type, typ) && c.Value.ExactString() == val.ExactString() {
			return c
		}
	}

	con := fn.NewValue(Value{
		Op:    op.Const,
		Type:  typ,
		Value: val,
	})
	fn.Consts = append(fn.Consts, con)
	return con
}

func (fn *Func) FixedReg(reg reg.Reg) *Value {
	for _, c := range fn.Consts {
		if c.Value == nil && c.Reg == reg {
			return c
		}
	}

	con := fn.NewValue(Value{
		Op:   op.Reg,
		Type: types.Typ[types.Int],
		Reg:  reg,
	})
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

	for _, blk := range fn.Blocks {
		str += blk.LongString()
	}

	return str
}
