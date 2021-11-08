package codegen

import (
	"fmt"
	"io"

	"github.com/rj45/rj32/gorj/ir"
)

func GenerateCode(mod *ir.Module, out io.Writer) error {
	gen := &gen{
		mod:            mod,
		out:            out,
		emittedGlobals: make(map[*ir.Value]bool),
	}
	for _, fn := range mod.Funcs {
		gen.genFunc(fn)
	}
	return nil
}

type gen struct {
	mod *ir.Module
	out io.Writer

	emittedGlobals map[*ir.Value]bool

	section string
	indent  string
}

func (gen *gen) emit(fmtstr string, args ...interface{}) {
	fmt.Fprintf(gen.out, gen.indent+fmtstr+"\n", args...)
}
