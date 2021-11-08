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

	if val.Args[1].Op != op.Const {
		log.Println("op not const!")
		return 0
	}

	amt, ok := constant.Int64Val(val.Args[1].Value)
	if !ok {
		panic("expected int64 constant")
	}

	// if amt == 1 {
	// TODO: no op, remove
	// }

	if amt == 0 {
		// TODO: is zero, replace with constant zero
		return 0
	}

	i := int64(1)
	n := int64(0)
	for i = 1; i < amt; i <<= 1 {
		n++
	}
	if i != amt {
		// TODO: can use multiple shifts and adds to calculate this
		log.Println("amt", amt, i)
		return 0
	}

	val.Op = op.ShiftLeft
	val.Args[1] = val.Block.Func.Const(val.Args[1].Type, constant.MakeInt64(n))

	return 1
}

var _ = addToPass(0, mulByConst)
