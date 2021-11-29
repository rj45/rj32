package rj32

import (
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
)

type Fmt int

const (
	BadFmt Fmt = iota
	LoadFmt
	StoreFmt
	MoveFmt
	CompareFmt
	BinaryFmt
	UnaryFmt
	CallFmt
	NoFmt
)

var templates = []string{
	"",
	"%s, [%s, %s]",
	"[%s, %s], %s",
	"%s, %s",
	"%s, %s",
	"%s, %s",
	"%s",
	"%s",
	"",
}

func (f Fmt) Template() string {
	return templates[f]
}

func (f Fmt) Vars(val *ir.Value) []*asm.Var {
	switch f {
	case LoadFmt:
		return []*asm.Var{varFor(val), varFor(val.Arg(0)), varFor(val.Arg(1))}
	case StoreFmt:
		return []*asm.Var{varFor(val.Arg(0)), varFor(val.Arg(1)), varFor(val.Arg(2))}
	case MoveFmt:
		return []*asm.Var{varFor(val), varFor(val.Arg(0))}
	case CompareFmt:
		return []*asm.Var{varFor(val.Arg(0)), varFor(val.Arg(1))}
	case BinaryFmt:
		return []*asm.Var{varFor(val), varFor(val.Arg(1))}
	case UnaryFmt:
		return []*asm.Var{varFor(val)}
	case CallFmt:
		return []*asm.Var{varFor(val.Arg(0))}
	}
	return nil
}

func varFor(val *ir.Value) *asm.Var {
	return &asm.Var{Value: val}
}
