package xform

import (
	"go/constant"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func mulByConst(val *ir.Value) int {
	if val.Op != op.Mul {
		return 0
	}

	if val.Arg(1).Op != op.Const {
		log.Println("op not const!")
		return 0
	}

	amt, ok := constant.Int64Val(val.Arg(1).Value)
	if !ok {
		panic("expected int64 constant")
	}

	if amt == 1 {
		val.ReplaceWith(val.Arg(0))
		return 1
	}

	if amt == 0 {
		val.ReplaceWith(val.Arg(1))
		return 1
	}

	i := int64(1)
	n := int64(0)
	for i = 1; i < amt; i <<= 1 {
		n++
	}
	if i != amt {
		// TODO: can use multiple shifts and adds to calculate this
		return 0
	}

	val.Op = op.ShiftLeft
	val.ReplaceArg(1, val.Func().Const(val.Arg(1).Type, constant.MakeInt64(n)))

	return 1
}

var _ = addToPass(Simplification, mulByConst)
