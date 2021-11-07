package codegen

import (
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
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

		gen.emit("%s:  ; %s", constant.StringVal(glob.Value), typ)

		size := gen.sizes.Sizeof(typ) / 2
		if size < 1 {
			size = 1
		}

		gen.emit("\t#res %d", size)
	}

	if gen.section != "code" {
		gen.emit("\n#bank code")
		gen.section = "code"
	}

	gen.emit("\n; %s", fn.Type)
	gen.emit("%s:", fn.Name)

	var retblock *ir.Block

	for _, blk := range fn.Blocks {
		if blk.Op == op.Return {
			if retblock != nil {
				log.Fatalf("two return blocks! %s", fn.LongString())
			}

			retblock = blk
			continue
		}

		gen.genBlock(blk)
	}

	if retblock != nil {
		gen.genBlock(retblock)
	}
}

var ifcc = map[op.Op]string{
	op.Equal:        "eq",
	op.NotEqual:     "ne",
	op.Less:         "lt",
	op.LessEqual:    "le",
	op.GreaterEqual: "ge",
	op.Greater:      "gt",
}

func (gen *gen) genBlock(blk *ir.Block) {
	gen.emit(".%s:", blk)
	gen.indent = "\t"

	if blk.Op == op.If && blk.Controls[0].Op.IsCompare() {
		blk.RemoveInstr(blk.Controls[0])
	}

	for _, instr := range blk.Instrs {
		switch instr.Op {
		case op.Load:
			gen.genLoad(instr)
			continue

		case op.Store:
			gen.genStore(instr)
			continue

		case op.Call:
			if len(instr.Args) != 1 {
				gen.emit("; %s", instr.LongString())
				continue
			}
		}

		name := instr.Op.Asm()
		if name != "" {
			for len(name) < 6 {
				name += " "
			}
			switch len(instr.Args) {
			case 0:
				gen.emit("%s", name)
			case 1:
				gen.emit("%s %s", name, instr.Args[0])
			case 2:
				gen.emit("%s %s, %s", name, instr.Args[0], instr.Args[1])
			case 3:
				gen.emit("%s %s, %s, %s", name, instr.Args[0], instr.Args[1], instr.Args[1])
			default:
				gen.emit("; %s", instr.LongString())
			}
		} else {
			gen.emit("; %s", instr.LongString())
		}
	}

	switch blk.Op {
	case op.Jump:
		gen.emit("jump .%s", blk.Succs[0].Block)

	case op.Return:
		gen.emit("return")

	case op.If:
		ctrl := blk.Controls[0]
		if ctrl.Op.IsCompare() {
			sign := ""
			if isUnsigned(ctrl.Type) && ctrl.Op != op.Equal && ctrl.Op != op.NotEqual {
				sign = "u"
			}
			gen.emit("if.%s%s %s, %s", sign, ifcc[ctrl.Op], ctrl.Args[0], ctrl.Args[1])
		} else {
			gen.emit("if.ne %s, 0", ctrl)
		}
		gen.emit("\tjump .%s", blk.Succs[0].Block)
		gen.emit("jump .%s", blk.Succs[1].Block)

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
