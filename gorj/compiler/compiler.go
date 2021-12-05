// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.
package compiler

import (
	"bytes"
	"flag"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rj45/rj32/gorj/codegen"
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/goenv"
	"github.com/rj45/rj32/gorj/html"
	"github.com/rj45/rj32/gorj/parser"
	"github.com/rj45/rj32/gorj/regalloc"
	"github.com/rj45/rj32/gorj/xform"
)

type Arch interface {
	Name() string
	AssemblerFormat() string
	EmulatorCmd() string
	EmulatorArgs() []string
}

func SetArch(a Arch) {
	arch = a
}

var arch Arch

type dumper interface {
	WritePhase(string, string)
	WriteAsmBuf(string, *bytes.Buffer)
	WriteAsm(string, *asm.Func)
	WriteSources(phase string, fn string, lines []string, startline int)
	Close()
}

type nopDumper struct{}

func (nopDumper) WritePhase(string, string)                                           {}
func (nopDumper) WriteAsmBuf(string, *bytes.Buffer)                                   {}
func (nopDumper) WriteAsm(string, *asm.Func)                                          {}
func (nopDumper) WriteSources(phase string, fn string, lines []string, startline int) {}
func (nopDumper) Close()                                                              {}

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
	var finalout io.WriteCloser
	var asmout io.WriteCloser

	if outname == "-" {
		finalout = nopWriteCloser{os.Stdout}
	} else {
		f, err := os.Create(outname)
		if err != nil {
			log.Fatal(err)
		}
		finalout = f
	}

	asmout = finalout

	var binfile string
	var asmcmd *exec.Cmd
	if assemble {
		// todo: if specified, allow this to not be a temp file
		asmtemp, err := os.CreateTemp("", "gorj_*.asm")
		if err != nil {
			log.Fatalln("failed to create temp asm file for customasm:", err)
		}
		defer os.Remove(asmtemp.Name())

		bintemp, err := os.CreateTemp("", "gorj_*.bin")
		if err != nil {
			log.Fatalln("failed to create temp bin file for customasm:", err)
		}
		bintemp.Close() // customasm will write to it
		binfile = bintemp.Name()
		// defer os.Remove(bintemp.Name())

		root := goenv.Get("GORJROOT")
		path := filepath.Join(root, "arch", arch.Name(), "customasm")
		cpudef := filepath.Join(path, "cpudef.asm")
		rungo := filepath.Join(path, "rungo.asm")

		asmcmd = exec.Command("customasm", "-q",
			"-f", arch.AssemblerFormat(),
			"-o", bintemp.Name(),
			cpudef, rungo, asmtemp.Name())
		log.Println(asmcmd)
		asmcmd.Stderr = os.Stderr
		asmcmd.Stdout = os.Stdout
		asmout = asmtemp
	}

	var runcmd *exec.Cmd
	if run {
		args := arch.EmulatorArgs()
		args = append(args, binfile)
		if *trace {
			args = append(args, "-trace")
		}
		runcmd = exec.Command(arch.EmulatorCmd(), args...)
		runcmd.Stderr = os.Stderr
		runcmd.Stdout = finalout
		runcmd.Stdin = os.Stdin
	}

	log.SetFlags(log.Lshortfile)

	parser := parser.NewParser(dir, patterns...)
	parser.Scan()

	pkg := parser.Package()

	gen := codegen.NewGenerator(pkg)
	emit := asm.NewEmitter(asmout)

	for fn := parser.NextUnparsedFunc(); fn != nil; fn = parser.NextUnparsedFunc() {
		var w dumper
		w = nopDumper{}
		if *dump != "" && strings.Contains(fn.Name, *dump) {
			w = html.NewHTMLWriter("ssa.html", fn)
			filename, lines, start := parser.DumpOrignalSource(fn)
			w.WriteSources("go", filename, lines, start)
			w.WriteAsmBuf("tools/go/ssa", parser.DumpOriginalSSA(fn))
		}
		defer w.Close()

		parser.ParseFunc(fn)

		w.WritePhase("initial", "initial")

		xform.AddReturnMoves(fn)
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

		asm := gen.Func(fn)
		w.WriteAsm("asm", asm)
		emit.Func(asm)
	}

	asmout.Close()

	if asmcmd != nil {
		if err := asmcmd.Run(); err != nil {
			os.Exit(1)
		}
		if !run {
			// todo: read file and emit to finalout
			f, err := os.Open(binfile)
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(finalout, f)
			f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if runcmd != nil {
		if err := runcmd.Run(); err != nil {
			return 1
		}
		return runcmd.ProcessState.ExitCode()
	}

	return 0
}
