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

func (val *Value) BlockID() ID {
	return val.block.id
}

func (val *Value) Func() *Func {
	return val.block.fn
}

func (val *Value) Index() int {
	if val.block == nil {
		return -1
	}
	if val.block.instrs[val.index] != val {
		panic("index out of sync")
	}
	return val.index
}

func (val *Value) NeedsReg() bool {
	if val.Op.IsSink() {
		return false
	}
	if val.Op == op.Call {
		sig := val.Type.(*types.Signature)
		return sig.Results().Len() == 1
	}
	return !val.Op.IsConst() && val.Op != op.Reg
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

func (val *Value) ArgIndex(arg *Value) int {
	for i, a := range val.args {
		if a == arg {
			return i
		}
	}
	return -1
}

func (val *Value) IsAfter(other *Value) bool {
	return val.block.IsAfter(other.block) || val.Index() > other.Index()
}

func (val *Value) Remove() {
	val.block.RemoveInstr(val)
}

func (val *Value) ReplaceWith(other *Value) bool {
	changed := len(val.argUses) > 0 || len(val.blockUses) > 0

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

func (val *Value) ReplaceOtherUsesWith(other *Value) bool {
	changed := len(val.argUses) > 0 || len(val.blockUses) > 0

	tries := 0
	for len(val.argUses) > 1 {
		tries++
		use := val.argUses[len(val.argUses)-1]
		if use == other {
			use = val.argUses[len(val.argUses)-2]
		}
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
			if instr == other {
				continue
			}
			for _, arg := range instr.args {
				if arg == val {
					panic("leaking uses")
				}
			}
		}
		return true
	})

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
		if val.Value == nil {
			return "nil"
		}
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
	if val == nil {
		return "<nil>"
	}
	if val.block == nil {
		return fmt.Sprintf("g%d", val.ID())
	}
	return fmt.Sprintf("v%d", val.ID())
}

func (val *Value) ShortString() string {
	str := ""

	if val.Op.IsSink() {
		str += "      "
	} else {
		if val.Reg != reg.None {
			str += fmt.Sprintf("%s:%s", val.IDString(), val.Reg)
		} else {
			str += val.IDString()
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
		if arg.Reg != reg.None {
			str += fmt.Sprintf("%s:%s", arg.IDString(), arg.Reg)
		} else {
			str += arg.String()
		}
	}

	if val.Value != nil {
		if len(val.args) > 0 {
			str += ", "
		}
		str += val.Value.String()
	}

	return str
}

func (val *Value) LongString() string {
	str := val.ShortString()

	if val.Type != nil {
		typstr := val.Type.String()

		for (len(str) + len(typstr)) < 64 {
			str += " "
		}

		str += typstr
	}

	return str
}

// FindPathTo traverses the value graph until a certain value is found
// and returns the path
func (val *Value) FindPathTo(fn func(*Value) bool) []*Value {
	path, found := val.findPathTo(fn, nil, make(map[*Value]bool))
	if found {
		return path
	}
	return nil
}

func (val *Value) findPathTo(fn func(*Value) bool, stack []*Value, visited map[*Value]bool) ([]*Value, bool) {
	stack = append(stack, val)

	if fn(val) {
		return stack, true
	}

	if visited[val] {
		return stack, false
	}
	visited[val] = true

	for _, arg := range val.args {
		var found bool
		stack, found = arg.findPathTo(fn, stack, visited)
		if found {
			return stack, found
		}
	}

	stack = stack[:len(stack)-1]
	return stack, false
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

func (val *Value) NumUses() int {
	return len(val.argUses) + len(val.blockUses)
}

func (val *Value) NumArgUses() int {
	return len(val.argUses)
}

func (val *Value) ArgUse(i int) *Value {
	return val.argUses[i]
}

func (val *Value) NumBlockUses() int {
	return len(val.blockUses)
}

func (val *Value) BlockUse(i int) *Block {
	return val.blockUses[i]
}

// FindUsageSuccessorPaths will search the successor block graph for
// all unique paths to each use of this value. These will be the paths
// through the CFG where the value is live.
func (val *Value) FindUsageSuccessorPaths() [][]*Block {
	// determine all the block successor paths through which the value is live
	var paths [][]*Block

	pathCache := make(map[*Block][]*Block)

	for i := 0; i < val.NumArgUses(); i++ {
		use := val.ArgUse(i)

		path, found := pathCache[use.Block()]
		if !found {
			path = val.Block().FindPathTo(func(b *Block) bool {
				return b == use.Block()
			})
			pathCache[use.Block()] = path
		}

		paths = unifyPaths(path, paths)
	}

	for i := 0; i < val.NumBlockUses(); i++ {
		use := val.BlockUse(i)

		path, found := pathCache[use]
		if !found {
			path = val.Block().FindPathTo(func(b *Block) bool {
				return b == use
			})
			pathCache[use] = path
		}

		paths = unifyPaths(path, paths)
	}

	return paths
}

func unifyPaths(path []*Block, paths [][]*Block) [][]*Block {
	// for each path
nextpath:
	for p, epath := range paths {
		// for each node in the new path
		for k, node := range path {
			// if a node doesn't match in the existing path
			if k < len(epath) && epath[k] != node {
				// move on to the next path
				continue nextpath
			}
			// else if we can extend the existing path further do so
			if k >= len(epath) {
				paths[p] = append(paths[p], node)
			}
		}

		// if we got here, we found a matching path
		return paths
	}

	// no paths match, add this one
	paths = append(paths, path)
	return paths
}
