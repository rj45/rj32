package ir

import (
	"fmt"
	"go/types"
)

type Func struct {
	Name string
	Type *types.Signature

	Mod *Module

	Blocks []*Block

	Consts []*Value
	Params []*Value
	Calls  []*Value

	blockID idAlloc
	instrID idAlloc
}

func (fn *Func) NextBlockID() ID {
	return fn.blockID.next()
}

func (fn *Func) NextInstrID() ID {
	return fn.instrID.next()
}

func (fn *Func) String() string {
	return fn.Name
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
