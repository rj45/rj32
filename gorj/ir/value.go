package ir

import (
	"fmt"
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
)

type Value struct {
	ID   ID
	Op   op.Op
	Type types.Type

	Value constant.Value

	Block *Block

	Args []*Value
}

func (val *Value) String() string {
	switch val.Op {
	case op.Const:
		return fmt.Sprintf("$%s", val.Value.String())
	case op.Parameter:
		return fmt.Sprintf("p%s", val.Value.String())
	case op.Func:
		return fmt.Sprintf("f%s", val.Value.String())
	case op.Global:
		return fmt.Sprintf("g%s", val.Value.String())
	}
	return fmt.Sprintf("v%d", val.ID)
}

func (val *Value) LongString() string {
	str := ""

	if val.Op.Def().Sink {
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
