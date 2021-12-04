package codegen

import (
	"go/constant"
	"strings"

	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func (gen *Generator) genBlock(blk, next *ir.Block) {
	asmBlk := &asm.Block{
		Label: "." + blk.String(), // todo: move to '.' prepending to rj32 pkg
		Block: blk,
	}
	gen.fn.Blocks = append(gen.fn.Blocks, asmBlk)

	suppressedInstrs := make(map[*ir.Value]bool)

	if blk.Op == op.If && blk.Control(0).Op.IsCompare() {
		suppressedInstrs[blk.Control(0)] = true
	}

	for i := 0; i < blk.NumInstrs(); i++ {
		instr := blk.Instr(i)

		if suppressedInstrs[instr] {
			continue
		}

		if instr.Op == op.InlineAsm {
			asmBlk.Instrs = inlineAsm(asmBlk.Instrs, instr)
			return
		}

		asmBlk.Instrs = arch.AssembleInstr(asmBlk.Instrs, instr)
	}

	flipSuccs := blk.NumSuccs() == 2 && blk.Succ(0) == next
	asmBlk.Instrs = arch.AssembleBlockOp(asmBlk.Instrs, blk, flipSuccs)

	// if the last instruction refers solely to the next block, skip it
	lastInstr := asmBlk.Instrs[len(asmBlk.Instrs)-1]
	if len(lastInstr.Args) == 1 && lastInstr.Args[0].Block == next && next != nil {
		asmBlk.Instrs = asmBlk.Instrs[:len(asmBlk.Instrs)-1]
	}
}

type inlineAsmFmt struct {
}

func (inlineAsmFmt) Template() string {
	return ""
}

func (inlineAsmFmt) Vars(val *ir.Value) []*asm.Var {
	return nil
}

type inlineAsmOp struct {
	line string
}

func (ia *inlineAsmOp) Asm() string {
	return ia.line
}

func (ia *inlineAsmOp) Fmt() asm.Fmt {
	return inlineAsmFmt{}
}

func (ia *inlineAsmOp) IsMove() bool {
	return false
}

func (ia *inlineAsmOp) IsCall() bool {
	return false
}

func inlineAsm(list []*asm.Instr, val *ir.Value) []*asm.Instr {
	text := constant.StringVal(val.Value)
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		list = append(list, &asm.Instr{
			Op: &inlineAsmOp{line: line},
		})
	}

	return list
}
