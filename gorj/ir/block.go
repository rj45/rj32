package ir

import (
	"fmt"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Block struct {
	id ID
	Op op.BlockOp

	Comment string

	Controls []*Value

	fn *Func

	instrs []*Value

	succs []*Block
	preds []*Block

	Idom     *Block
	Dominees []*Block
}

func (blk *Block) ID() ID {
	return blk.id
}

func (blk *Block) Func() *Func {
	return blk.fn
}

func (blk *Block) NumSuccs() int {
	return len(blk.succs)
}

func (blk *Block) Succ(i int) *Block {
	return blk.succs[i]
}

func (blk *Block) NumPreds() int {
	return len(blk.preds)
}

func (blk *Block) Pred(i int) *Block {
	return blk.preds[i]
}

func (blk *Block) NumInstrs() int {
	return len(blk.instrs)
}

func (blk *Block) Instr(i int) *Value {
	return blk.instrs[i]
}

func (blk *Block) String() string {
	return fmt.Sprintf("b%d", blk.ID())
}

func (blk *Block) LongString() string {
	str := fmt.Sprintf("%s:", blk)

	if len(blk.Comment) > 0 {
		for len(str) < 9 {
			str += " "
		}
		str += fmt.Sprintf(" ; %s", blk.Comment)
	}

	if len(blk.preds) > 0 || len(blk.succs) > 0 {
		cfg := "; CFG"

		if len(blk.preds) > 0 {
			for _, pred := range blk.preds {
				cfg += fmt.Sprintf(" %s", pred.String())
			}
			cfg += " ->"
		}

		cfg += " "
		cfg += blk.String()

		if len(blk.succs) > 0 {
			cfg += " ->"
			for _, succ := range blk.succs {
				cfg += fmt.Sprintf(" %s", succ.String())
			}
		}

		max := 40
		for (len(cfg)+max) > 68 && max > 0 {
			max--
		}

		for len(str) < max {
			str += " "
		}

		str += cfg
	}

	str += "\n"

	for i := 0; i < blk.NumInstrs(); i++ {
		instr := blk.Instr(i)
		str += fmt.Sprintf("    %s\n", instr.LongString())
	}

	opstr := fmt.Sprintf("%s ", blk.Op)

	for len(opstr) < 10 {
		opstr += " "
	}
	for i, arg := range blk.Controls {
		if i != 0 {
			opstr += ", "
		}
		opstr += arg.String()
	}

	succstr := ""
	if len(blk.succs) == 1 {
		succstr = blk.Succ(0).String()
	} else if len(blk.succs) == 2 {
		succstr = fmt.Sprintf("then %s else %s", blk.Succ(0), blk.Succ(1))
	}

	if len(blk.Controls) > 0 {
		opstr += " "
	}

	str += fmt.Sprintf("          %s%s\n", opstr, succstr)

	return str
}

func (blk *Block) InsertInstr(i int, val *Value) {
	val.block = blk
	if i < 0 || i >= len(blk.instrs) {
		val.index = len(blk.instrs)
		blk.instrs = append(blk.instrs, val)
		return
	}

	val.index = i
	blk.instrs = append(blk.instrs[:i+1], blk.instrs[i:]...)
	blk.instrs[i] = val

	for j := i + 1; j < len(blk.instrs); j++ {
		blk.instrs[j].index = j
	}
}

func (blk *Block) InsertCopy(i int, val *Value, reg reg.Reg) *Value {
	opr := op.Copy
	if reg.IsStackSlot() {
		opr = op.Store
	}
	newval := blk.fn.NewValue(opr, val.Type, val)
	newval.Reg = reg
	blk.InsertInstr(i, newval)
	return newval
}

func (blk *Block) AddSucc(succ *Block) {
	blk.succs = append(blk.succs, succ)
}

func (blk *Block) AddPred(pred *Block) {
	blk.preds = append(blk.preds, pred)
}

func (blk *Block) VisitSuccessors(fn func(*Block) bool) {
	blk.visitSuccessors(fn, make(map[ID]bool))
}

func (blk *Block) visitSuccessors(fn func(*Block) bool, visited map[ID]bool) {
	visited[blk.ID()] = true
	if !fn(blk) {
		return
	}
	for _, succ := range blk.succs {
		if !visited[succ.ID()] {
			succ.visitSuccessors(fn, visited)
		}
	}
}

func (blk *Block) RemoveInstr(val *Value) bool {
	i := val.Index()
	if i < 0 {
		return false
	}

	blk.instrs = append(blk.instrs[:i], blk.instrs[i+1:]...)

	for j := i; j < len(blk.instrs); j++ {
		blk.instrs[j].index = j
	}

	return true
}
