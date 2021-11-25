package ir

import (
	"fmt"
	"log"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Block struct {
	id ID
	Op op.BlockOp

	Comment string
	Source  string

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

func (blk *Block) InsertControl(i int, val *Value) {
	val.blockUses = append(val.blockUses, blk)

	if i < 0 || i >= len(blk.controls) {
		blk.controls = append(blk.controls, val)
		return
	}

	blk.controls = append(blk.controls[:i+1], blk.controls[i:]...)
}

func (blk *Block) RemoveControl(i int) {
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

	blk.controls = append(blk.controls[:i], blk.controls[i+1:]...)
}

func (blk *Block) IsAfter(other *Block) bool {
	return blk.isAfter(other, make(map[*Block]bool))
}

func (blk *Block) isAfter(other *Block, visited map[*Block]bool) bool {
	visited[blk] = true
	if blk == other || visited[blk] {
		return true
	}

	for _, pred := range blk.preds {
		if pred.isAfter(other, visited) {
			return true
		}
	}

	return false
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

	str += blk.OpString()

	return str
}

func (blk *Block) OpString() string {
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

	return fmt.Sprintf("          %s%s\n", opstr, succstr)
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

func (blk *Block) SwapInstr(a *Value, b *Value) {
	i := a.Index()
	j := b.Index()

	blk.instrs[i], blk.instrs[j] = blk.instrs[j], blk.instrs[i]

	a.index = j
	b.index = i
}

func (blk *Block) InsertCopy(i int, val *Value, r reg.Reg) *Value {
	var newval *Value
	if r.IsStackSlot() {
		if !val.NeedsReg() {
			val = blk.InsertCopy(i, val, reg.None)
			i++
		}
		newval = blk.fn.NewValue(op.Store, val.Type, blk.fn.FixedReg(reg.SP), blk.fn.IntConst(int64(r.StackSlot())), val)
	} else {
		newval = blk.fn.NewValue(op.Copy, val.Type, val)
		newval.Reg = r
	}
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
	list := blk.reverseVisitSuccessors(nil, make([]bool, blk.fn.BlockIDCount()))
	for i := len(list) - 1; i >= 0; i-- {
		if !fn(list[i]) {
			return
		}
	}
}

func (blk *Block) reverseVisitSuccessors(list []*Block, visited []bool) []*Block {
	visited[blk.ID()] = true

	for i := blk.NumSuccs() - 1; i >= 0; i-- {
		succ := blk.Succ(i)
		if !visited[succ.ID()] {
			list = succ.reverseVisitSuccessors(list, visited)
		}
	}

	return append(list, blk)
}

func (blk *Block) RemoveInstr(val *Value) bool {
	i := val.Index()
	if i < 0 {
		return false
	}

	if len(val.argUses) > 0 {
		log.Panicf("attempted to remove val %s:%s still in use", val.IDString(), val.String())
	}

	if len(val.blockUses) > 0 {
		log.Panicf("attempted to remove val %s:%s still in use", val.IDString(), val.String())
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
