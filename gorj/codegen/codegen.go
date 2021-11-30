package codegen

import (
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/codegen/asm/rj32"
	"github.com/rj45/rj32/gorj/ir"
)

type Generator struct {
	mod *ir.Package

	emittedGlobals map[*ir.Value]bool

	arch asm.Arch
	fn   *asm.Func
}

func NewGenerator(mod *ir.Package) *Generator {
	return &Generator{
		mod:            mod,
		emittedGlobals: make(map[*ir.Value]bool),
		arch:           rj32.Rj32{},
	}
}
