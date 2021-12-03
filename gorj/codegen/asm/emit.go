package asm

import (
	"fmt"
	"io"
	"strings"
)

type Emitter struct {
	out io.Writer

	src []string

	section Section
	indent  string
}

func NewEmitter(out io.Writer) *Emitter {
	return &Emitter{
		out: out,
	}
}

func (emit *Emitter) Func(fn *Func) {
	for _, glob := range fn.Globals {
		emit.global(glob)
	}

	emit.ensureSection(Code)
	emit.line("")
	emit.line("; %s", fn.Comment)
	emit.line("%s:", fn.Label)

	for _, blk := range fn.Blocks {
		emit.block(blk)
	}

	emit.line("")
}

func (emit *Emitter) global(glob *Global) {
	emit.ensureSection(glob.Section)

	emit.line("%s:", glob.Label)
	emit.indent = "  "
	if glob.Comment != "" {
		emit.line("; %s", glob.Comment)
	}
	for _, l := range glob.Strings {
		emit.line("%s", l)
	}
	emit.indent = ""
	emit.line("")
}

func (emit *Emitter) ensureSection(section Section) {
	if emit.section != section {
		emit.line("#bank %s", section)
		emit.section = section
	}
}

func (emit *Emitter) block(blk *Block) {
	emit.line("%s:", blk.Label)
	emit.indent = "  "

	// this will emit bits of the original source to aid debugging
	emit.source(blk.Block.Source)

	for _, instr := range blk.Instrs {
		emit.instr(instr)
	}
	emit.indent = ""
}

func (emit *Emitter) arg(v *Var) string {
	if v.String != "" {
		return v.String
	}
	if v.Value != nil {
		return v.Value.String()
	}
	if v.Block != nil {
		return "." + v.Block.String()
	}
	return ""
}

func (emit *Emitter) instr(instr *Instr) {
	templ := instr.Op.Fmt().Template()

	name := instr.Op.Asm()

	if instr.Indent {
		name = "  " + name
	}

	for len(name) < 6 {
		name += " "
	}

	strs := []interface{}{name}
	for _, arg := range instr.Args {
		strs = append(strs, emit.arg(arg))
	}

	emit.line("%s "+templ, strs...)
}

func (emit *Emitter) line(fmtstr string, args ...interface{}) {
	nextline := ""
	if len(emit.src) > 0 {
		nextline, emit.src = emit.src[len(emit.src)-1], emit.src[:len(emit.src)-1]
	}
	output := fmt.Sprintf(emit.indent+fmtstr, args...)

	if nextline != "" {
		for len(output) < 30 {
			output += " "
		}
		output += "; "
		output += nextline
	}

	fmt.Fprintln(emit.out, output)
}

func (emit *Emitter) source(src string) {
	if src == "" {
		return
	}

	lines := strings.Split(src, "\n")
	var revlines []string
	for i := len(lines) - 1; i >= 0; i-- {
		revlines = append(revlines, lines[i])
	}
	revlines = append(revlines, emit.src...)
	emit.src = revlines
}
