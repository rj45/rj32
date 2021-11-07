package codegen

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func (gen *gen) genLoad(instr *ir.Value) {
	if len(instr.Args) == 1 && instr.Args[0].Op == op.Global {
		gen.emit("load %s, [gp, %s]", instr, instr.Args[0])
		return
	}

	gen.emit("; %s", instr.LongString())
}

func (gen *gen) genStore(instr *ir.Value) {
	if len(instr.Args) == 2 && instr.Args[0].Op == op.Global {
		gen.emit("store [gp, %s], %s", instr.Args[0], instr.Args[1])
		return
	}
	gen.emit("; %s", instr.LongString())
}
