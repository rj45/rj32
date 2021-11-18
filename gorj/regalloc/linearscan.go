package regalloc

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

// prerequisites:
//  - no critical edges (done)
//  - phi nodes are decomposed into parallel copies

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

	intervals []interval

	active       []*interval
	futureActive []*interval
}

func (ls *linearScan) scan() {
	ls.blockPos = make([]pos, ls.fn.BlockIDCount())
	ls.visited = make([]bool, ls.fn.BlockIDCount())

	ls.fn.Blocks()[0].VisitSuccessors(ls.livenessScan)
}

func (ls *linearScan) pos(val *ir.Value, offset pos) pos {
	return ls.blockPos[val.Block().ID()] + pos(val.Index()<<1) + offset
}

func (ls *linearScan) livenessScan(blk *ir.Block) bool {

	ls.blockPos[blk.ID()] = pos(len(ls.blocks) << 16)
	ls.blocks = append(ls.blocks, blk)

	return true
}
