package codegen

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (gen *gen) genLoad(instr *ir.Value) {
	if len(instr.Args) == 1 {
		if instr.Args[0].Op == op.Global || instr.Args[0].Op == op.Const {
			gen.emit("load   %s, [gp, %s]", instr, instr.Args[0])
			return
		}

		if instr.Args[0].Op == op.Add && instr.Args[0].Args[0].Reg == reg.GP {
			gen.emit("load   %s, [%s]", instr, instr.Args[0])
			return
		}
	}

	gen.emit("; %s", instr.LongString())
}

func (gen *gen) genStore(instr *ir.Value) {
	if len(instr.Args) == 2 {
		if instr.Args[0].Op == op.Global || instr.Args[0].Op == op.Const {
			gen.emit("store  [gp, %s], %s", instr.Args[0], instr.Args[1])
			return
		}

		if instr.Args[0].Op == op.Add && instr.Args[0].Args[0].Reg == reg.GP {
			gen.emit("store  [%s], %s", instr.Args[0], instr.Args[1])
			return
		}
	}
	gen.emit("; %s", instr.LongString())
}
