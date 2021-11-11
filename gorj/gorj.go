// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package main

import (
	"fmt"
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

	mod := parser.ParseModule("./testfiles/seive/seive.go")

	fmt.Println(mod.LongString())

	for _, fn := range mod.Funcs {
		w := html.NewHTMLWriter(fn.Name+".html", fn)
		w.WritePhase("initial", "initial")

		xform.Transform(xform.Elaboration, fn)
		w.WritePhase("elaboration", "elaboration")

		xform.Transform(xform.Simplification, fn)
		w.WritePhase("simplification", "simplification")

		used := regalloc.Allocate(fn)
		w.WritePhase("allocation", "allocation")

		xform.Transform(xform.LastPass, fn)
		w.WritePhase("cleanup", "cleanup")

		xform.ProEpiLogue(used, fn)
		w.WritePhase("final", "final")

		w.Close()
	}

	fmt.Print("\n\n--------------------\n\n")

	codegen.GenerateCode(mod, os.Stdout)
}
