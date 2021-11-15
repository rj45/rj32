package codegen

import (
	"fmt"
	"io"

	"github.com/rj45/rj32/gorj/ir"
)

type Generator struct {
	mod *ir.Module
	out io.Writer

	emittedGlobals map[*ir.Value]bool

	section string
	indent  string
}

func NewGenerator(mod *ir.Module) *Generator {
	return &Generator{
		mod:            mod,
		emittedGlobals: make(map[*ir.Value]bool),
	}
}

func (gen *Generator) emit(fmtstr string, args ...interface{}) {
	fmt.Fprintf(gen.out, gen.indent+fmtstr+"\n", args...)
}
