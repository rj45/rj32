// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rj45/rj32/gorj/codegen"
	"github.com/rj45/rj32/gorj/html"
	"github.com/rj45/rj32/gorj/parser"
	"github.com/rj45/rj32/gorj/regalloc"
	"github.com/rj45/rj32/gorj/xform"
)

type dumper interface {
	WritePhase(string, string)
	WriteAsm(string, *bytes.Buffer)
	Close()
}

type nopDumper struct{}

func (nopDumper) WritePhase(string, string)      {}
func (nopDumper) WriteAsm(string, *bytes.Buffer) {}
func (nopDumper) Close()                         {}

type nopWriteCloser struct{ w io.Writer }

func (nopWriteCloser) Close() error {
	return nil
}

func (n nopWriteCloser) Write(p []byte) (int, error) {
	return n.w.Write(p)
}

var _ io.WriteCloser = nopWriteCloser{}

func main() {
	dump := flag.String("dump", "", "Dump a function to ssa.html")
	output := flag.String("o", "", "output file for the result")

	flag.Parse()

	command := flag.Arg(0)
	printUsage := flag.NArg() < 2

	assemble := false
	run := false
	_ = run

	switch command {
	case "b", "build":
		assemble = true
	case "r", "run":
		assemble = true
		run = true
	case "s", "asm":
	default:
		printUsage = true
	}

	if printUsage {
		fmt.Fprintln(os.Stderr, "gorj - a go compiler for rj32")
		fmt.Fprintln(os.Stderr, "https://github.com/rj45/rj32/gorj")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Usage: gorj <flags> <command> <packages...>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  build: compile and assemble with customasm")
		fmt.Fprintln(os.Stderr, "  asm: compile and write assembly to file")
		fmt.Fprintln(os.Stderr, "  run: compile, assemble and run emulator")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := flag.Arg(1)

	dir := filepath.Dir(filename)
	basename := filepath.Base(filename)
	outname := "-"
	if *output != "" {
		outname = *output
	}

	var out io.WriteCloser
	if outname == "-" {
		out = nopWriteCloser{os.Stdout}
	} else {
		f, err := os.Create(outname)
		if err != nil {
			log.Fatal(err)
		}
		out = f
	}

	// var runcmd *exec.Cmd
	// if run {
	// 	runcmd = exec.Command("emu", "-novdp", "-run", "-")

	// }

	var asmcmd *exec.Cmd
	if assemble {
		// todo: if specified, allow this to not be a temp file
		tempfile, err := os.CreateTemp("", "gorj_*.asm")
		if err != nil {
			log.Fatalln("failed to pipe stdin to customasm:", err)
		}
		defer os.Remove(tempfile.Name())

		asmcmd = exec.Command("customasm", "-q", "-p", "-f", "logisim16", "/home/rj45/rj32/programs/cpudef.asm", "/home/rj45/rj32/programs/rungo.asm", tempfile.Name())
		asmcmd.Stderr = os.Stderr

		asmcmd.Stdout = out
		out = tempfile
	}

	log.SetFlags(log.Lshortfile)

	mod := parser.ParseModule(dir, basename)

	gen := codegen.NewGenerator(mod)

	for _, fn := range mod.Funcs {
		var w dumper
		w = nopDumper{}
		if *dump != "" && strings.Contains(fn.Name, *dump) {
			w = html.NewHTMLWriter("ssa.html", fn)
		}
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
		gen.Func(fn, io.MultiWriter(out, buf))
		w.WriteAsm("asm", buf)

	}
	out.Close()

	if asmcmd != nil {
		if err := asmcmd.Run(); err != nil {
			os.Exit(1)
		}
	}

}
