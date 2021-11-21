// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type RegAlloc struct {
	Func *ir.Func

	usedRegs reg.Reg

	guessedRegs  map[*ir.Value]bool
	wrongGuesses map[*ir.Value]bool

	blockInfo []blockInfo

	copiesEliminated          int
	potentialCopiesEliminated int

	// translation table for values that were spilled and later reloaded
	spillReloads map[*ir.Value]*ir.Value

	liveThroughCalls map[*ir.Value]bool
}

type blockInfo struct {
	// map[instr][]args
	kills    map[*ir.Value][]*ir.Value
	blkKills map[*ir.Value]bool

	phiIns  map[*ir.Block]map[*ir.Value]bool
	phiOuts map[*ir.Block]map[*ir.Value]bool

	liveIns  map[*ir.Value]bool
	liveOuts map[*ir.Value]bool

	regValues map[reg.Reg]*ir.Value
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

	log.Println("Copies eliminated:", ra.copiesEliminated, "out of potentially:", ra.potentialCopiesEliminated)

	return ra.usedRegs
}

func (ra *RegAlloc) Verify() {
	ra.verify(false)
	ra.verify(true)
}
