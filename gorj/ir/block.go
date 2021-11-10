package ir

import (
	"fmt"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Block struct {
	ID ID
	Op op.BlockOp

	Comment string

	Controls []*Value

	Func *Func

	Instrs []*Value

	Succs []BlockRef
	Preds []BlockRef

	Idom     *Block
	Dominees []*Block
}

func (blk *Block) NextInstrID() ID {
	return blk.Func.NextInstrID()
}

func (blk *Block) String() string {
	return fmt.Sprintf("b%d", blk.ID)
}

func (blk *Block) LongString() string {
	str := fmt.Sprintf("%s:", blk)

	if len(blk.Comment) > 0 {
		for len(str) < 9 {
			str += " "
		}
		str += fmt.Sprintf(" ; %s", blk.Comment)
	}

	if len(blk.Preds) > 0 || len(blk.Succs) > 0 {
		cfg := "; CFG"

		if len(blk.Preds) > 0 {
			for _, pred := range blk.Preds {
				cfg += fmt.Sprintf(" %s", pred.String())
			}
			cfg += " ->"
		}

		cfg += " "
		cfg += blk.String()

		if len(blk.Succs) > 0 {
			cfg += " ->"
			for _, succ := range blk.Succs {
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

	for _, instr := range blk.Instrs {
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
	if len(blk.Succs) == 1 {
		succstr = blk.Succs[0].Block.String()
	} else if len(blk.Succs) == 2 {
		succstr = fmt.Sprintf("then %s else %s", blk.Succs[0].Block, blk.Succs[1].Block)
	}

	if len(blk.Controls) > 0 {
		opstr += " "
	}

	str += fmt.Sprintf("          %s%s\n", opstr, succstr)

	return str
}

func (blk *Block) InsertInstr(i int, val *Value) {
	val.Block = blk
	if i < 0 || i >= len(blk.Instrs) {
		val.Index = len(blk.Instrs)
		blk.Instrs = append(blk.Instrs, val)
		return
	}

	val.Index = i
	blk.Instrs = append(blk.Instrs[:i+1], blk.Instrs[i:]...)
	blk.Instrs[i] = val

	for j := i + 1; j < len(blk.Instrs); j++ {
		blk.Instrs[j].Index = j
	}
}

func (blk *Block) InsertCopy(i int, val *Value, reg reg.Reg) *Value {
	opr := op.Copy
	if reg.IsStackSlot() {
		opr = op.Store
	}
	newval := blk.Func.NewValue(Value{
		Op:   opr,
		Reg:  reg,
		Args: []*Value{val},
		Type: val.Type,
	})
	blk.InsertInstr(i, newval)
	return newval
}

func (blk *Block) IndexOf(val *Value) int {
	if blk.Instrs[val.Index] != val {
		found := false
		for i, instr := range blk.Instrs {
			if val == instr {
				val.Index = i
				return i
			}
		}
		if !found {
			return -1
		}
	}
	return val.Index
}

func (blk *Block) VisitSuccessors(fn func(*Block) bool) {
	blk.visitSuccessors(fn, make(map[ID]bool))
}

func (blk *Block) visitSuccessors(fn func(*Block) bool, visited map[ID]bool) {
	visited[blk.ID] = true
	if !fn(blk) {
		return
	}
	for _, succ := range blk.Succs {
		if !visited[succ.Block.ID] {
			succ.Block.visitSuccessors(fn, visited)
		}
	}
}

func SubstituteValue(from *Value, to *Value) {
	from.Block.VisitSuccessors(func(blk *Block) bool {
		for _, instr := range blk.Instrs {
			for i, arg := range instr.Args {
				if arg.ID == from.ID {
					instr.Args[i] = to
				}
			}
		}
		return true
	})
}

func (blk *Block) RemoveInstr(val *Value) bool {
	i := blk.IndexOf(val)
	if i < 0 {
		return false
	}

	blk.Instrs = append(blk.Instrs[:i], blk.Instrs[i+1:]...)

	for j := i; j < len(blk.Instrs); j++ {
		blk.Instrs[j].Index = j
	}

	return true
}

type BlockRef struct {
	Index int
	Block *Block
}

func (ref *BlockRef) String() string {
	return fmt.Sprintf("%d:%v", ref.Index, ref.Block)
}
