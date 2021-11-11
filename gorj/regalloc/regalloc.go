package regalloc

import (
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

// The idea is to use "tree register allocation".
//
// This is essentially a linear scan allocator but slightly modified to
// work very well on SSA IR, and does not need the phis removed first.
//
// This works as long as the dominance tree has no "critical edges" and
// forms "chordal graphs" (which it should by default).
//
// For each block in the control flow graph from the return to the entry
//  - keep track of when values are last read, mark these as kills
//    - a value may have multiple kills if its last read is in multiple blocks
//    - a value may have no kills if it's in a loop
//	- keep a list of live-ins and live-outs of each block
//  - when a phi is seen, consider it like a copy in the the pred block
//    - it will be a live-out to the phi value?
//    - it may kill the value at the end of the previous block
//  - add phis and copies to an affinity map/graph
//    - in the next pass it will try to assign these to the same register
//
// For each block in the control flow graph from the entry to exit
//  - start with a list of available registers
//    - for each live-in of the block
//      - remove the register from the available ones
//    - for each instruction
//      - if instruction kills a value
//        - free the register for that value
//      - if an instruction defines a value
//        - allocate a register
//          - check the affinity map if any other value has a register already
//            - check which register is the most common
//              - try to use that register if it's free
//          - pick the next free register
//            - use saved registers if the value is live-out
//            - use temporary registers if value does not leave block
//              - round robin can help with architectures that stall on read after write
//          - if no free registers
//            - could do "live range splitting"
//            - could spill but then I think the process needs to restart
//              - spills are best added before register allocation
//                - can do a "max cardinality search" to find them?

func Allocate(fn *ir.Func) reg.Reg {
	ra := regAlloc{Func: fn}
	ra.liveScan()
	ra.colour()
	return ra.usedRegs
}

type regAlloc struct {
	Func *ir.Func

	usedRegs reg.Reg

	affinities map[ir.ID][]*ir.Value
	blockInfo  []blockInfo
}

type blockInfo struct {
	// map[instr][]args
	kills    map[ir.ID][]ir.ID
	blkKills map[ir.ID]bool

	phiIns  map[ir.ID]bool
	phiOuts map[ir.ID]bool

	liveIns  map[ir.ID]bool
	liveOuts map[ir.ID]bool
}
