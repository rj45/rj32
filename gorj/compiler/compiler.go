// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.
package compiler

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
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

var dump = flag.String("dump", "", "Dump a function to ssa.html")
var trace = flag.Bool("trace", false, "debug program with tracing info")

func Compile(outname, dir string, patterns []string, assemble, run bool) int {
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

	var runcmd *exec.Cmd
	if run {
		args := []string{"-run", "-"}
		if *trace {
			args = append(args, "-trace")
		}
		runcmd = exec.Command("emurj", args...)
		runcmd.Stderr = os.Stderr
		runcmd.Stdout = out

		var err error
		out, err = runcmd.StdinPipe()
		if err != nil {
			log.Fatalln("failed to pipe stdin to emurj:", err)
		}
	}

	var asmcmd *exec.Cmd
	if assemble {
		// todo: if specified, allow this to not be a temp file
		tempfile, err := os.CreateTemp("", "gorj_*.asm")
		if err != nil {
			log.Fatalln("failed to create temp file for customasm:", err)
		}
		defer os.Remove(tempfile.Name())

		asmcmd = exec.Command("customasm", "-q", "-p", "-f", "logisim16", "/home/rj45/rj32/programs/cpudef.asm", "/home/rj45/rj32/programs/rungo.asm", tempfile.Name())
		asmcmd.Stderr = os.Stderr

		asmcmd.Stdout = out
		out = tempfile
	}

	log.SetFlags(log.Lshortfile)

	mod := parser.ParseProgram(dir, patterns...)

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

		ra := regalloc.NewRegAlloc(fn)

		used := ra.Allocate(fn)
		w.WritePhase("allocation", "allocation")

		ra.Verify()

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
		asmcmd.Stdout.(io.WriteCloser).Close()
	}

	if runcmd != nil {
		if err := runcmd.Run(); err != nil {
			return 1
		}
		return runcmd.ProcessState.ExitCode()
	}

	return 0
}