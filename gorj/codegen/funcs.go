package codegen

import (
	"go/constant"
	"go/types"
	"io"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/sizes"
)

func (gen *Generator) Func(fn *ir.Func, out io.Writer) {
	gen.out = out
	for _, glob := range fn.Globals {
		if gen.emittedGlobals[glob] {
			continue
		}
		gen.emittedGlobals[glob] = true

		if gen.section != "bss" {
			gen.emit("\n#bank bss")
			gen.section = "bss"
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
