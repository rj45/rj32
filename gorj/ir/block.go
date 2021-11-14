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

	controls []*Value

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

func (blk *Block) NumControls() int {
	return len(blk.controls)
}

func (blk *Block) Control(i int) *Value {
	return blk.controls[i]
}

func (blk *Block) SetControls(ctrls []*Value) {
	blk.controls = ctrls

	for _, c := range ctrls {
		c.blockUses = append(c.blockUses, blk)
	}
}

func (blk *Block) ReplaceControl(i int, val *Value) {
	oldval := blk.controls[i]
	found := false
	for i, use := range oldval.blockUses {
		if use == blk {
			oldval.blockUses = append(oldval.blockUses[:i], oldval.blockUses[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		panic("replacement without use!")
	}

	val.blockUses = append(val.blockUses, blk)
	blk.controls[i] = val
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
	for i, arg := range blk.controls {
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

	if len(blk.controls) > 0 {
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

	for len(val.args) > 0 {
		val.RemoveArg(len(val.args) - 1)
	}

	return true
}

// FindPathTo searches the successor graph for a specific block and
// returns the path to that block
func (blk *Block) FindPathTo(fn func(*Block) bool) []*Block {
	path, found := blk.findPathTo(fn, nil, make(map[*Block]bool))
	if found {
		return path
	}
	return nil
}

func (blk *Block) findPathTo(fn func(*Block) bool, stack []*Block, visited map[*Block]bool) ([]*Block, bool) {
	stack = append(stack, blk)

	if fn(blk) {
		return stack, true
	}

	if visited[blk] {
		return stack, false
	}
	visited[blk] = true

	for _, succ := range blk.succs {
		var found bool
		stack, found = succ.findPathTo(fn, stack, visited)
		if found {
			return stack, found
		}
	}

	stack = stack[:len(stack)-1]
	return stack, false
}
