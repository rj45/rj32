// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type RegAlloc struct {
	Func *ir.Func

	usedRegs reg.Reg

	guessedRegs  map[*ir.Value]bool
	wrongGuesses map[*ir.Value]bool

	affinities map[*ir.Value][]*ir.Value
	blockInfo  []blockInfo

	// translation table for values that were spilled and later reloaded
	spillReloads map[*ir.Value]*ir.Value

	liveThroughCalls map[*ir.Value]bool
}

func NewRegAlloc(fn *ir.Func) *RegAlloc {
	return &RegAlloc{
		Func:         fn,
		spillReloads: make(map[*ir.Value]*ir.Value),
	}
}

func (ra *RegAlloc) Allocate(fn *ir.Func) reg.Reg {
	// ra.alloc.scan()

	ra.liveScan()
	ra.colour()
	return ra.usedRegs
}

func (ra *RegAlloc) Verify() {
	ra.verify(false)
	ra.verify(true)
}

type blockInfo struct {
	// map[instr][]args
	kills    map[*ir.Value][]*ir.Value
	blkKills map[*ir.Value]bool

	phiIns  map[*ir.Block]map[*ir.Value]bool
	phiOuts map[*ir.Block]map[*ir.Value]bool

	liveIns  map[*ir.Value]bool
	liveOuts map[*ir.Value]bool

	spills    map[*ir.Value]int
	freeSlots []int

	regValues map[reg.Reg]*ir.Value
}
