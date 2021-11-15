// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/rj45/rj32/gorj/codegen"
	"github.com/rj45/rj32/gorj/html"
	"github.com/rj45/rj32/gorj/parser"
	"github.com/rj45/rj32/gorj/regalloc"
	"github.com/rj45/rj32/gorj/xform"
)

func main() {
	log.SetFlags(log.Lshortfile)

	mod := parser.ParseModule("./testdata/seive/seive.go")
	// mod := parser.ParseModule("./testdata/fib/fib.go")

	gen := codegen.NewGenerator(mod)

	fmt.Println(mod.LongString())

	for _, fn := range mod.Funcs {
		w := html.NewHTMLWriter(fn.Name+".html", fn)
		defer w.Close()
		w.WritePhase("initial", "initial")

		xform.Transform(xform.Elaboration, fn)
		w.WritePhase("elaboration", "elaboration")

		xform.Transform(xform.Simplification, fn)
		w.WritePhase("simplification", "simplification")

		xform.Transform(xform.Lowering, fn)
		w.WritePhase("lowering", "lowering")

		used := regalloc.Allocate(fn)
		w.WritePhase("allocation", "allocation")

		xform.Transform(xform.Legalize, fn)
		w.WritePhase("legalize", "legalize")

		xform.Transform(xform.CleanUp, fn)
		w.WritePhase("cleanup", "cleanup")

		xform.ProEpiLogue(used, fn)
		xform.EliminateEmptyBlocks(fn)
		w.WritePhase("final", "final")

		buf := &bytes.Buffer{}
		gen.Func(fn, io.MultiWriter(os.Stdout, buf))
		w.WriteAsm("asm", buf)
	}

}
