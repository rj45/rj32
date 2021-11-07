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

func (gen *gen) genBlock(blk *ir.Block) {
	gen.emit(".%s:", blk)
	gen.indent = "\t"

	for _, instr := range blk.Instrs {
		gen.emit("%s", instr.LongString())
	}

	switch blk.Op {
	case op.Jump:
		gen.emit("jump .%s", blk.Succs[0].Block)

	case op.Return:
		gen.emit("return")

	case op.If:
		if len(blk.Controls[0].Args) == 2 {
			gen.emit("if.%s %s, %s", blk.Controls[0].Op, blk.Controls[0].Args[0], blk.Controls[0].Args[1])
		} else {
			gen.emit("if.ne %s, 0", blk.Controls[0])
		}
		gen.emit("\tjump .%s", blk.Succs[0].Block)
		gen.emit("jump .%s", blk.Succs[1].Block)

	default:
		panic("unimpl")
	}

	gen.indent = ""
}
