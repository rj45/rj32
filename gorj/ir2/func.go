package ir2

import (
	"go/types"

	"github.com/rj45/rj32/gorj/ir/op"
)

type ID uint16

type ValueIndex uint16
type InstrIndex uint16

type ValueType uint16

type Value struct {
	name  ID
	typ   ValueType
	instr InstrIndex
}

type Instr struct {
	op   op.Op
	defs []ValueIndex
	args []ValueIndex
}

type NewValue struct {
	Value
	insertAt ValueIndex
}

type Block struct {
	fn *Func

	phiStart InstrIndex
	start    InstrIndex

	end     InstrIndex
	control InstrIndex
}

type Func struct {
	id  ID
	pkg *Package

	values    []Value
	newValues []NewValue
}

type Package struct {
	types  []types.Type
	typeID map[types.Type]ValueType
}
