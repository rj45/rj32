package ir

import (
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
)

type Module struct {
	Funcs   []*Func
	Globals []*Value

	valId idAlloc
}

func (mod *Module) LongString() string {
	str := ""

	for i, fn := range mod.Funcs {
		if i != 0 {
			str += "\n"
		}
		str += fn.LongString()
	}

	return str
}

func (mod *Module) AddGlobal(name string, typ types.Type) {
	mod.Globals = append(mod.Globals, &Value{
		id:    mod.valId.next(),
		Op:    op.Global,
		Value: constant.MakeString(name),
		Type:  typ,
	})
}
