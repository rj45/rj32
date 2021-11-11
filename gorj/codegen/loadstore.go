package codegen

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (gen *gen) genLoad(instr *ir.Value) {
	if instr.NumArgs() == 1 {
		if instr.Arg(0).Op == op.Global || instr.Arg(0).Op == op.Const {
			gen.emit("load   %s, [gp, %s]", instr, instr.Arg(0))
			return
		}

		if instr.Arg(0).Op == op.Add && instr.Arg(0).Arg(0).Reg == reg.GP {
			gen.emit("load   %s, [%s]", instr, instr.Arg(0))
			return
		}
	} else if instr.NumArgs() == 2 {
		gen.emit("load   %s, [%s, %s]", instr, instr.Arg(0), instr.Arg(1))
		return
	}

	gen.emit("; %s", instr.LongString())
}

func (gen *gen) genStore(instr *ir.Value) {
	if instr.NumArgs() == 2 {
		if instr.Arg(0).Op == op.Global || instr.Arg(0).Op == op.Const {
			gen.emit("store  [gp, %s], %s", instr.Arg(0), instr.Arg(1))
			return
		}

		if instr.Arg(0).Op == op.Add && instr.Arg(0).Arg(0).Reg == reg.GP {
			gen.emit("store  [%s], %s", instr.Arg(0), instr.Arg(1))
			return
		}
	} else if instr.NumArgs() == 3 {
		gen.emit("store  [%s, %s], %s", instr.Arg(0), instr.Arg(1), instr.Arg(2))
		return
	}
	gen.emit("; %s", instr.LongString())
}
