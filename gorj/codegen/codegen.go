package codegen

import (
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
)

type Arch interface {
	AssembleGlobal(glob *ir.Value) *asm.Global
	AssembleInstr(list []*asm.Instr, val *ir.Value) []*asm.Instr
	AssembleBlockOp(list []*asm.Instr, blk *ir.Block, flip bool) []*asm.Instr
}

var arch Arch

func SetArch(a Arch) {
	arch = a
}

type Generator struct {
	mod *ir.Package

	emittedGlobals map[*ir.Value]bool

	fn *asm.Func
}

func NewGenerator(mod *ir.Package) *Generator {
	return &Generator{
		mod:            mod,
		emittedGlobals: make(map[*ir.Value]bool),
	}
}
