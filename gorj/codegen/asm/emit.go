package asm

import (
	"fmt"
	"io"
	"strings"
)

type Emitter struct {
	prog *Program
	out  io.Writer

	src []string

	section string
	indent  string
}

func NewEmitter(out io.Writer, prog *Program) *Emitter {
	return &Emitter{
		out:  out,
		prog: prog,
	}
}

func (gen *Emitter) emit(fmtstr string, args ...interface{}) {
	nextline := ""
	if len(gen.src) > 0 {
		nextline, gen.src = gen.src[len(gen.src)-1], gen.src[:len(gen.src)-1]
	}
	output := fmt.Sprintf(gen.indent+fmtstr, args...)

	if nextline != "" {
		for len(output) < 40 {
			output += " "
		}
		output += "; "
		output += nextline
	}

	fmt.Fprintln(gen.out, output)
}

func (gen *Emitter) source(src string) {
	if src == "" {
		return
	}

	lines := strings.Split(src, "\n")
	var revlines []string
	for i := len(lines) - 1; i >= 0; i-- {
		revlines = append(revlines, lines[i])
	}
	revlines = append(revlines, gen.src...)
	gen.src = revlines
}
