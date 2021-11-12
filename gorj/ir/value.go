package ir

import (
	"fmt"
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Value struct {
	id   ID
	Reg  reg.Reg
	Op   op.Op
	Type types.Type

	Value constant.Value

	block *Block
	index int

	blockUses []*Block
	argUses   []*Value

	args []*Value
}

func (val *Value) ID() ID {
	return val.id
}

func (val *Value) Block() *Block {
	return val.block
}

func (val *Value) Func() *Func {
	return val.block.fn
}

func (val *Value) Index() int {
	if val.block.instrs[val.index] != val {
		panic("index out of sync")
	}
	return val.index
}

func (val *Value) NumArgs() int {
	return len(val.args)
}

func (val *Value) Arg(i int) *Value {
	return val.args[i]
}

func (val *Value) ReplaceArg(i int, arg *Value) {
	if val == arg {
		panic("attempt to replace like with like")
	}

	val.args[i].removeUse(val)
	val.args[i] = arg
	val.args[i].addUse(val)
}

func (val *Value) RemoveArg(i int) *Value {
	oldval := val.args[i]
	oldval.removeUse(val)

	val.args = append(val.args[:i], val.args[i+1:]...)

	return oldval
}

func (val *Value) InsertArg(i int, arg *Value) {
	arg.addUse(val)

	if i < 0 || i >= len(val.args) {
		val.args = append(val.args, arg)
		return
	}

	val.args = append(val.args[:i+1], val.args[i:]...)
	val.args[i] = arg
}

func (val *Value) Remove() {
	val.block.RemoveInstr(val)
}

func (val *Value) ReplaceWith(other *Value) bool {
	changed := len(val.argUses) > 0 || len(val.blockUses) > 0
	// val.block.VisitSuccessors(func(blk *Block) bool {
	// 	for i := 0; i < blk.NumInstrs(); i++ {
	// 		instr := blk.Instr(i)
	// 		for i, arg := range instr.args {
	// 			if arg.ID() == val.ID() {
	// 				instr.ReplaceArg(i, other)
	// 				instr.args[i] = other
	// 				changed = true
	// 			}
	// 		}
	// 	}
	// 	return true
	// })

	tries := 0
	for len(val.argUses) > 0 {
		tries++
		use := val.argUses[len(val.argUses)-1]
		if tries > 1000 {
			log.Panicf("bug in arguses %v, %v, %v, %v", val, other, val.argUses, use.args)
		}
		found := false
		for i, arg := range use.args {
			if arg == val {
				use.ReplaceArg(i, other)
				found = true
				break
			}
		}
		if !found {
			panic("couldn't find use!")
		}
	}

	tries = 0
	for len(val.blockUses) > 0 {
		tries++
		if tries > 1000 {
			panic("bug in block uses")
		}

		use := val.blockUses[len(val.blockUses)-1]

		found := false
		for i, ctrl := range use.controls {
			if ctrl == val {
				use.ReplaceControl(i, other)
				found = true
				break
			}
		}
		if !found {
			panic("couldn't find use!")
		}
	}

	val.block.VisitSuccessors(func(blk *Block) bool {
		for i := 0; i < blk.NumInstrs(); i++ {
			instr := blk.Instr(i)
			for _, arg := range instr.args {
				if arg.ID() == val.ID() {
					panic("leaking uses")
				}
			}
		}
		return true
	})

	val.Remove()

	return changed
}

func (val *Value) String() string {
	if val.Reg != reg.None {
		if val.Reg.IsAReg() {
			return val.Reg.String()
		}
		if val.Reg.IsStackSlot() {
			return fmt.Sprintf("[sp, %d]", val.Reg.StackSlot())
		}
	}
	switch val.Op {
	case op.Const:
		if val.Value.Kind() == constant.Bool {
			if val.Value.String() == "true" {
				return "1"
			}
			return "0"
		}
		return val.Value.String()
	case op.Parameter, op.Func, op.Global:
		return constant.StringVal(val.Value)
	}
	return fmt.Sprintf("v%d", val.ID())
}

func (val *Value) IDString() string {
	if val.block == nil {
		return fmt.Sprintf("g%d", val.ID())
	}
	return fmt.Sprintf("v%d", val.ID())
}

func (val *Value) LongString() string {
	str := ""

	if val.Op.IsSink() {
		str += "      "
	} else {
		if val.Reg != reg.None {
			str += val.Reg.String()
		} else {
			str += fmt.Sprintf("v%d", val.ID())
		}
		for len(str) < 3 {
			str += " "
		}
		str += " = "
	}
	str += fmt.Sprintf("%s ", val.Op.String())
	for len(str) < 16 {
		str += " "
	}
	for i, arg := range val.args {
		if i != 0 {
			str += ", "
		}
		if val.Op == op.Phi {
			str += fmt.Sprintf("%s:", val.Block().Pred(i))
		}
		str += arg.String()
	}

	if val.Value != nil {
		if len(val.args) > 0 {
			str += ", "
		}
		str += val.Value.String()
	}

	if val.Type != nil {
		typstr := val.Type.String()

		for (len(str) + len(typstr)) < 64 {
			str += " "
		}

		str += typstr
	}

	return str
}

func (val *Value) addUse(other *Value) {
	if other == val {
		panic("trying to add self use")
	}
	val.argUses = append(val.argUses, other)
}

func (val *Value) removeUse(other *Value) {
	if other == val {
		panic("trying to remove self use")
	}
	index := -1
	for i, use := range val.argUses {
		if use == other {
			index = i
			break
		}
	}
	if index < 0 {
		uses := ""
		for _, use := range val.argUses {
			uses += " " + use.IDString()
		}
		log.Panicf("%s:%s does not have use %s:%s, %v", val.IDString(), val.LongString(), other.IDString(), other.LongString(), uses)
	}
	val.argUses = append(val.argUses[:index], val.argUses[index+1:]...)
}
