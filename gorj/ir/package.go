package ir

import (
	"go/constant"
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
)

type Package struct {
	Funcs   []*Func
	Globals []*Value

	valId idAlloc
}

func (pkg *Package) LongString() string {
	str := ""

	for i, fn := range pkg.Funcs {
		if i != 0 {
			str += "\n"
		}
		str += fn.LongString()
	}

	return str
}

func (pkg *Package) AddGlobal(name string, typ types.Type) *Value {
	val := &Value{
		id:    pkg.valId.next(),
		Op:    op.Global,
		Value: constant.MakeString(name),
		Type:  typ,
	}
	pkg.Globals = append(pkg.Globals, val)
	return val
}
