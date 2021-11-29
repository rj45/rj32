package codegen

import (
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

var ifcc = map[op.Op]string{
	op.Equal:        "eq",
	op.NotEqual:     "ne",
	op.Less:         "lt",
	op.LessEqual:    "le",
	op.GreaterEqual: "ge",
	op.Greater:      "gt",
}

func (gen *Generator) genBlock(blk, next *ir.Block) {
	asmBlk := &asm.Block{
		Label: "." + blk.String(), // todo: move to '.' prepending to rj32 pkg
		Block: blk,
	}
	gen.fn.Blocks = append(gen.fn.Blocks, asmBlk)

	gen.emit(".%s:", blk)
	gen.indent = "  "

	gen.source(blk.Source)

	suppressedInstrs := make(map[*ir.Value]bool)

	if blk.Op == op.If && blk.Control(0).Op.IsCompare() {
		suppressedInstrs[blk.Control(0)] = true
	}

	for i := 0; i < blk.NumInstrs(); i++ {
		instr := blk.Instr(i)

		if suppressedInstrs[instr] {
			continue
		}

		asmBlk.Instrs = gen.arch.AssembleInstr(asmBlk.Instrs, instr)

		name := instr.Op.Asm()
		if name != "" {
			for len(name) < 6 {
				name += " "
			}
		}

		ret := ""
		if !instr.Op.IsSink() {
			ret = " "
			ret += instr.String()
			if instr.NumArgs() > 0 {
				ret += ","
			}
		}

		switch instr.Op {
		case op.Load:
			gen.emit("load   %s, [%s, %s]", instr, instr.Arg(0), instr.Arg(1))
			continue

		case op.Store:
			gen.emit("store  [%s, %s], %s", instr.Arg(0), instr.Arg(1), instr.Arg(2))
			continue

		case op.Call:
			gen.emit("%s %s", name, instr.Arg(0))
			continue

		case op.SwapOut:
			// ignore, the SwapIn will produce the instructions
			continue

		case op.Extract:
			if instr.Arg(0).Op == op.Call {
				continue
			}

		case op.Phi:
			gen.emit("; %s", instr.ShortString())
			continue
		}

		if instr.Op.IsCompare() {
			sign := ""
			space := " "
			if isUnsigned(instr.Arg(0).Type) && instr.Op != op.Equal && instr.Op != op.NotEqual {
				sign = "u"
				space = ""
			}
			gen.emit("move   %s, 0", instr)
			gen.emit("if.%s%s%s %s, %s", sign, ifcc[instr.Op], space, instr.Arg(0), instr.Arg(1))
			gen.emit("  move %s, 1", instr)

			continue
		}

		if name != "" {
			for len(name) < 6 {
				name += " "
			}
			if instr.Op.ClobbersArg() {
				if instr.Reg != instr.Arg(0).Reg {
					log.Panicf("expected %s to have dest and first source the same: %s", instr.Op, instr.LongString())
				}
				switch instr.NumArgs() {
				case 1:
					ret = ret[:len(ret)-1]
					gen.emit("%s%s", name, ret)
				case 2:
					gen.emit("%s%s %s", name, ret, instr.Arg(1))
				case 3:
					gen.emit("%s%s %s, %s", name, ret, instr.Arg(1), instr.Arg(1))
				default:
					gen.emit("; %s", instr.ShortString())
				}
				continue
			}
			switch instr.NumArgs() {
			case 0:
				gen.emit("%s%s", name, ret)
			case 1:
				gen.emit("%s%s %s", name, ret, instr.Arg(0))
			case 2:
				gen.emit("%s%s %s, %s", name, ret, instr.Arg(0), instr.Arg(1))
			case 3:
				gen.emit("%s%s %s, %s, %s", name, ret, instr.Arg(0), instr.Arg(1), instr.Arg(1))
			default:
				gen.emit("; %s", instr.ShortString())
			}
		} else {
			log.Panicf("unimplemented %s", instr.ShortString())
			gen.emit("; %s", instr.ShortString())
		}
	}

	flipSuccs := blk.NumSuccs() == 2 && blk.Succ(0) == next
	asmBlk.Instrs = gen.arch.AssembleBlockOp(asmBlk.Instrs, blk, flipSuccs)

	// if the last instruction refers solely to the next block, skip it
	lastInstr := asmBlk.Instrs[len(asmBlk.Instrs)-1]
	if len(lastInstr.Args) == 1 && lastInstr.Args[0].Block == next {
		asmBlk.Instrs = asmBlk.Instrs[:len(asmBlk.Instrs)-1]
	}

	switch blk.Op {
	case op.Jump:
		if blk.Succ(0) != next {
			gen.emit("jump   .%s", blk.Succ(0))
		}

	case op.Return:
		gen.emit("return")

	case op.Panic:
		gen.emit("move   a0, %s", blk.Control(0))
		gen.emit("error")

	case op.If, op.IfEqual, op.IfNotEqual, op.IfGreater,
		op.IfGreaterEqual, op.IfLessEqual, op.IfLess:
		ctrl := blk.Control(0)
		cond := op.NotEqual
		sign := ""
		space := " "
		arg1 := ctrl
		arg2 := "0"

		if blk.Op != op.If {
			cond = blk.Op.Compare()
			if isUnsigned(ctrl.Type) && cond != op.Equal && cond != op.NotEqual {
				sign = "u"
				space = ""
			}
			arg1 = ctrl
			arg2 = blk.Control(1).String()
		} else if ctrl.Op.IsCompare() {
			cond = ctrl.Op
			if isUnsigned(ctrl.Type) && ctrl.Op != op.Equal && ctrl.Op != op.NotEqual {
				sign = "u"
				space = ""
			}
			arg1 = ctrl.Arg(0)
			arg2 = ctrl.Arg(1).String()
		}

		succ0 := blk.Succ(0)
		succ1 := blk.Succ(1)
		if succ0 == next {
			cond = cond.Opposite()
			succ0, succ1 = succ1, succ0
		}

		gen.emit("if.%s%s%s %s, %s", sign, ifcc[cond], space, arg1, arg2)
		gen.emit("  jump .%s", succ0)

		if succ1 != next {
			gen.emit("jump   .%s", succ1)
		}

	default:
		panic("unimpl")
	}

	gen.indent = ""
}

func isUnsigned(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return true
	}
	return false
}
