package codegen

import (
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/sizes"
)

func (gen *gen) genFunc(fn *ir.Func) {
	for _, glob := range fn.Globals {
		if gen.emittedGlobals[glob] {
			continue
		}
		gen.emittedGlobals[glob] = true

		if gen.section != "data" {
			gen.emit("\n#bank data")
			gen.section = "data"
		}

		typ := glob.Type
		if ptr, ok := typ.(*types.Pointer); ok {
			typ = ptr.Elem()
		}

		size := sizes.Sizeof(typ)

		gen.emit("%s:  ; %s", constant.StringVal(glob.Value), typ)

		gen.emit("\t#res %d", size)
	}

	if gen.section != "code" {
		gen.emit("\n#bank code")
		gen.section = "code"
	}

	gen.emit("\n; %s", fn.Type)
	gen.emit("%s:", fn.Name)

	var retblock *ir.Block

	// order blocks by reverse succession
	blockList := reverseIRSuccessorSort(fn.Blocks()[0], nil, make(map[*ir.Block]bool))

	// reverse it to get succession ordering
	for i, j := 0, len(blockList)-1; i < j; i, j = i+1, j-1 {
		blockList[i], blockList[j] = blockList[j], blockList[i]
	}

	for i, blk := range blockList {
		if blk.Op == op.Return {
			if retblock != nil {
				log.Fatalf("two return blocks! %s", fn.LongString())
			}

			retblock = blk
			continue
		}

		var next *ir.Block
		if (i + 1) < len(blockList) {
			next = blockList[i+1]
		}

		gen.genBlock(blk, next)
	}

	if retblock != nil {
		gen.genBlock(retblock, nil)
	}
}

func reverseIRSuccessorSort(block *ir.Block, list []*ir.Block, visited map[*ir.Block]bool) []*ir.Block {
	visited[block] = true

	for i := block.NumSuccs() - 1; i >= 0; i-- {
		succ := block.Succ(i)
		if !visited[succ] {
			list = reverseIRSuccessorSort(succ, list, visited)
		}
	}

	return append(list, block)
}

var ifcc = map[op.Op]string{
	op.Equal:        "eq",
	op.NotEqual:     "ne",
	op.Less:         "lt",
	op.LessEqual:    "le",
	op.GreaterEqual: "ge",
	op.Greater:      "gt",
}

func (gen *gen) genBlock(blk, next *ir.Block) {
	gen.emit(".%s:", blk)
	gen.indent = "\t"

	suppressedInstrs := make(map[*ir.Value]bool)

	if blk.Op == op.If && blk.Control(0).Op.IsCompare() {
		suppressedInstrs[blk.Control(0)] = true
	}

	for i := 0; i < blk.NumInstrs(); i++ {
		instr := blk.Instr(i)

		if suppressedInstrs[instr] {
			continue
		}

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
				case 2:
					gen.emit("%s%s %s", name, ret, instr.Arg(1))
				case 3:
					gen.emit("%s%s %s, %s", name, ret, instr.Arg(1), instr.Arg(1))
				default:
					gen.emit("; %s", instr.LongString())
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
				gen.emit("; %s", instr.LongString())
			}
		} else {
			gen.emit("; %s", instr.LongString())
		}
	}

	switch blk.Op {
	case op.Jump:
		if blk.Succ(0) != next {
			gen.emit("jump   .%s", blk.Succ(0))
		}

	case op.Return:
		gen.emit("return")

	case op.If:
		ctrl := blk.Control(0)
		cond := op.NotEqual
		sign := ""
		space := " "
		arg1 := ctrl
		arg2 := "0"

		if ctrl.Op.IsCompare() {
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
