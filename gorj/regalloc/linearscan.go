package regalloc

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

// prerequisites:
//  - no critical edges (done)
//  - phi nodes are decomposed into parallel copies (done)

type pos int

const (
	before pos = -1
	during pos = 0
	after  pos = 1
)

func (p pos) block() int {
	return int(p) >> 16
}

func (p pos) instr() int {
	return (int(p) & 0xffff) >> 1
}

type liveRange struct {
	start pos
	end   pos

	next *liveRange
}

type interval struct {
	val *ir.Value

	ranges *liveRange
	reg    reg.Reg
}

type linearScan struct {
	fn *ir.Func

	blocks   []*ir.Block
	blockPos []pos
	visited  []bool
	info     []blockInfo2

	intervals []interval

	active       []*interval
	futureActive []*interval
}

type blockInfo2 struct {
	liveIn  map[*ir.Value]bool
	liveOut map[*ir.Value]bool
}

func (ls *linearScan) scan() {
	ls.blockPos = make([]pos, ls.fn.BlockIDCount())
	ls.visited = make([]bool, ls.fn.BlockIDCount())
	ls.info = make([]blockInfo2, ls.fn.BlockIDCount())
	ls.intervals = make([]interval, ls.fn.InstrIDCount())

	ls.fn.Blocks()[0].VisitSuccessors(ls.registerBlocks)

	ls.propagateLiveOuts()

	for i := len(ls.blocks) - 1; i >= 0; i-- {
		blk := ls.blocks[i]

		// find first instruction before PhiCopies
		j := blk.NumInstrs() - 1
		for ; j >= 0; j-- {
			val := blk.Instr(j)
			if val.Op != op.PhiCopy {
				break
			}
		}

		// firstRegularInstr := j

		// mark intervals for PhiCopies
		j = blk.NumInstrs() - 1
		for ; j >= 0; j-- {
			val := blk.Instr(j)
			if val.Op != op.PhiCopy {
				break
			}

			for i := 0; i < val.NumArgs(); i++ {
				// arg := val.Arg(0)
				// if ls.intervals[arg.ID()].
			}
		}

		// handle regular instructions
		for ; j >= 0; j-- {
			val := blk.Instr(j)
			if val.Op == op.Phi {
				break
			}

		}

		// handle phies
		for ; j >= 0; j-- {
			val := blk.Instr(j)
			if val.Op != op.Phi {
				panic("shouldn't happen")
			}

		}
	}
}

func (ls *linearScan) backPropagateLiveIns() {

	for i := len(ls.blocks) - 1; i >= 0; i-- {
		blk := ls.blocks[i]
		info := &ls.info[blk.ID()]

		for v := range info.liveOut {
			info.liveIn[v] = true
		}

		for j := blk.NumInstrs() - 1; j >= 0; j-- {
			val := blk.Instr(j)

			for k := 0; k < val.NumArgs(); k++ {
				arg := val.Arg(k)
				if arg.NeedsReg() {
					info.liveIn[arg] = true
				}
			}

			delete(info.liveIn, val)
		}
	}
}

func (ls *linearScan) propagateLiveOuts() {
	for i := len(ls.blocks) - 1; i >= 0; i-- {
		blk := ls.blocks[i]

		for i := 0; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			// find all instr inputs that are defined
			// outside this block and mark them as live out on that block
			for j := 0; j < val.NumArgs(); j++ {
				arg := val.Arg(j)
				var pred *blockInfo2

				if !arg.NeedsReg() {
					continue
				}

				if val.Op == op.Phi {
					pred = &ls.info[blk.Pred(j).ID()]
				} else if arg.Block() != blk {
					pred = &ls.info[arg.Block().ID()]
				}
				pred.liveOut[arg] = true

				// todo: traverse pred blocks until we find the def and mark
				// the blocks between has live in/out. Or we could do this as
				// a series of passes later to propagate all the values at once
			}
		}
	}
}

func (ls *linearScan) pos(val *ir.Value, offset pos) pos {
	return ls.blockPos[val.Block().ID()] + pos(val.Index()<<1) + offset
}

func (ls *linearScan) registerBlocks(blk *ir.Block) bool {
	ls.blockPos[blk.ID()] = pos(len(ls.blocks) << 16)
	ls.blocks = append(ls.blocks, blk)
	return true
}
