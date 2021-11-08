package ir

import (
	"fmt"
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Value struct {
	ID   ID
	Reg  reg.Reg
	Op   op.Op
	Type types.Type

	Value constant.Value

	Block *Block
	Index int

	Args []*Value
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
	return fmt.Sprintf("v%d", val.ID)
}

func (val *Value) LongString() string {
	str := ""

	if val.Op.IsSink() {
		str += "      "
	} else {
		str += fmt.Sprintf("v%d", val.ID)
		for len(str) < 3 {
			str += " "
		}
		str += " = "
	}
	str += fmt.Sprintf("%s ", val.Op.String())
	for len(str) < 16 {
		str += " "
	}
	for i, arg := range val.Args {
		if i != 0 {
			str += ", "
		}
		if val.Op == op.Phi {
			str += fmt.Sprintf("%s:", val.Block.Preds[i].Block)
		}
		str += arg.String()
	}

	if val.Value != nil {
		if len(val.Args) > 0 {
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
