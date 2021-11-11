package ir

import (
	"fmt"
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Value struct {
	id   ID
	Reg  reg.Reg
	Op   op.Op
	Type types.Type

	Value constant.Value

	block *Block
	index int

	args []*Value
}

func (val *Value) ID() ID {
	return val.id
}

func (val *Value) Block() *Block {
	return val.block
}

func (val *Value) Func() *Func {
	return val.block.fn
}

func (val *Value) Index() int {
	if val.block.Instrs[val.index] != val {
		panic("index out of sync")
	}
	return val.index
}

func (val *Value) NumArgs() int {
	return len(val.args)
}

func (val *Value) Arg(i int) *Value {
	return val.args[i]
}

func (val *Value) ReplaceArg(i int, arg *Value) {
	val.args[i] = arg
}

func (val *Value) RemoveArg(i int) *Value {
	oldval := val.args[i]

	val.args = append(val.args[:i], val.args[i+1:]...)

	for j := i; j < len(val.args); j++ {
		val.args[j].index = j
	}

	return oldval
}

func (val *Value) InsertArg(i int, arg *Value) {
	if i < 0 || i >= len(val.args) {
		val.args = append(val.args, arg)
		return
	}

	val.args = append(val.args[:i+1], val.args[i:]...)
	val.args[i] = arg
}

func (val *Value) Remove() {
	val.block.RemoveInstr(val)
}

func (val *Value) ReplaceWith(other *Value) bool {
	changed := false
	val.block.VisitSuccessors(func(blk *Block) bool {
		for _, instr := range blk.Instrs {
			for i, arg := range instr.args {
				if arg.ID() == val.ID() {
					instr.args[i] = other
					changed = true
				}
			}
		}
		return true
	})

	val.Remove()

	return changed
}

func (val *Value) String() string {
	if val.Reg != reg.None {
		if val.Reg.IsAReg() {
			return val.Reg.String()
		}
		if val.Reg.IsStackSlot() {
			return fmt.Sprintf("[sp, %d]", val.Reg.StackSlot())
		}
	}
	switch val.Op {
	case op.Const:
		if val.Value.Kind() == constant.Bool {
			if val.Value.String() == "true" {
				return "1"
			}
			return "0"
		}
		return val.Value.String()
	case op.Parameter, op.Func, op.Global:
		return constant.StringVal(val.Value)
	}
	return fmt.Sprintf("v%d", val.ID())
}

func (val *Value) LongString() string {
	str := ""

	if val.Op.IsSink() {
		str += "      "
	} else {
		if val.Reg != reg.None {
			str += val.Reg.String()
		} else {
			str += fmt.Sprintf("v%d", val.ID())
		}
		for len(str) < 3 {
			str += " "
		}
		str += " = "
	}
	str += fmt.Sprintf("%s ", val.Op.String())
	for len(str) < 16 {
		str += " "
	}
	for i, arg := range val.args {
		if i != 0 {
			str += ", "
		}
		if val.Op == op.Phi {
			str += fmt.Sprintf("%s:", val.Block().Preds[i])
		}
		str += arg.String()
	}

	if val.Value != nil {
		if len(val.args) > 0 {
			str += ", "
		}
		str += val.Value.String()
	}

	if val.Type != nil {
		typstr := val.Type.String()

		for (len(str) + len(typstr)) < 64 {
			str += " "
		}

		str += typstr
	}

	return str
}
