package asm

import (
	"github.com/rj45/rj32/gorj/ir"
)

type Section string

const (
	Code Section = "code"
	Data Section = "data"
	Bss  Section = "bss"
)

type Program struct {
	Pkg   *ir.Package
	Funcs []Func
}

type Func struct {
	Comment string
	Label   string

	Globals []*Global
	Blocks  []*Block
	Func    *ir.Func
}

type Global struct {
	Section Section
	Comment string
	Label   string
	Strings []string
	Value   *ir.Value
}

type Op interface {
	String() string
	Fmt() Fmt
	IsMove() bool
	IsCall() bool
}

type Fmt interface {
	Template() string
	Vars(val *ir.Value) []*Var
}

type Instr struct {
	Op     Op
	Args   []*Var
	Indent bool
}

type Block struct {
	Label  string
	Instrs []*Instr

	Block *ir.Block
}

type Var struct {
	String string
	Value  *ir.Value
	Block  *ir.Block
}

type Arch interface {
	AssembleGlobal(glob *ir.Value) *Global
	AssembleInstr(list []*Instr, val *ir.Value) []*Instr
	AssembleBlockOp(list []*Instr, blk *ir.Block, flip bool) []*Instr
}
