package codegen

import (
	"log"

	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func (gen *Generator) Func(fn *ir.Func) *asm.Func {
	gen.fn = &asm.Func{
		Comment: fn.Type.String(),
		Label:   fn.Name,
	}

	for _, glob := range fn.Globals {
		if gen.emittedGlobals[glob] {
			continue
		}
		gen.emittedGlobals[glob] = true

		gen.fn.Globals = append(gen.fn.Globals, gen.arch.AssembleGlobal(glob))
	}

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

	return gen.fn
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
