// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"flag"
	"fmt"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

var debugLiveness = flag.Bool("debugliveness", false, "emit register allocation liveness logs")

func (ra *RegAlloc) liveScan() {
	if *debugLiveness {
		fmt.Println("------------")
		fmt.Println(ra.Func.LongString())
		fmt.Println("------------")
	}

	ra.blockInfo = make([]blockInfo, ra.Func.BlockIDCount())
	ra.liveThroughCalls = make(map[*ir.Value]bool)

	entry := ra.Func.Blocks()[0]

	if entry.Idom != nil || entry.NumPreds() != 0 {
		log.Fatalln("expecting entry block to be top of the tree")
	}

	for _, blk := range ra.Func.Blocks() {
		info := &ra.blockInfo[blk.ID()]
		info.regValues = make(map[reg.Reg]*ir.Value)
	}

	// order blocks by reverse succession
	list := reverseIRSuccessorSort(ra.Func.Blocks()[0], nil, make(map[*ir.Block]bool))

	enqueued := make([]bool, ra.Func.BlockIDCount())

	for _, blk := range list {
		enqueued[blk.ID()] = true
	}

	// for each block on worklist until worklist empty
	for len(list) > 0 {
		// pop block off worklist
		blk := list[0]
		list = list[1:]
		enqueued[blk.ID()] = false

		// check if the usage has changed
		if ra.scanUsage(blk) {
			// if so, for each pred block
			for i := 0; i < blk.NumPreds(); i++ {
				pred := blk.Pred(i)

				// check if it's already on the worklist, and if not
				if !enqueued[pred.ID()] {
					// append it to the worklist
					enqueued[pred.ID()] = true
					list = append(list, pred)
				}
			}
		}
	}

	if *debugLiveness {
		ra.dumpLiveness(ra.Func.Blocks())
	}
}

func (ra *RegAlloc) dumpLiveness(list []*ir.Block) {
	for _, blk := range list {
		info := &ra.blockInfo[blk.ID()]

		fmt.Println("{{{------------")

		fmt.Println(blk.LongString())

		fmt.Print("  LiveIns: ")
		for val := range info.liveIns {
			fmt.Print(val.IDString())
			if val.Reg != reg.None {
				fmt.Printf(":%s", val.Reg)
			}
			if val.String() != val.IDString() {
				fmt.Printf(":%s", val.String())
			}
			fmt.Print(" ")
		}
		fmt.Println()

		fmt.Println("  PhiIns:")
		for i := 0; i < blk.NumPreds(); i++ {
			pred := blk.Pred(i)
			if len(info.phiIns[pred]) == 0 {
				continue
			}

			fmt.Printf("    -> %s: ", pred)
			for val := range info.phiIns[pred] {
				fmt.Print(val.IDString())
				if val.Reg != reg.None {
					fmt.Printf(":%s", val.Reg)
				}
				if val.String() != val.IDString() {
					fmt.Printf(":%s", val.String())
				}
				fmt.Print(" ")
			}
			fmt.Println()
		}

		fmt.Print("  Block kills: ")
		for val := range info.blkKills {
			fmt.Print(val.IDString())
			if val.Reg != reg.None {
				fmt.Printf(":%s", val.Reg)
			}
			if val.String() != val.IDString() {
				fmt.Printf(":%s", val.String())
			}
			fmt.Print(" ")
		}
		fmt.Println()

		fmt.Println("  Instr kills: ")
		for site, kills := range info.kills {
			fmt.Printf("    %d : %s -> ", site.Index(), site.IDString())
			for _, kill := range kills {
				fmt.Print(kill.IDString())
				if kill.Reg != reg.None {
					fmt.Printf(":%s", kill.Reg)
				}
				if kill.String() != kill.IDString() {
					fmt.Printf(":%s", kill.String())
				}
				fmt.Print(" ")
			}
			fmt.Println()
		}

		fmt.Print("  LiveOuts: ")
		for val := range info.liveOuts {
			fmt.Print(val.IDString())
			if val.Reg != reg.None {
				fmt.Printf(":%s", val.Reg)
			}
			if val.String() != val.IDString() {
				fmt.Printf(":%s", val.String())
			}
			fmt.Print(" ")
		}
		fmt.Println()

		fmt.Println("  PhiOuts:")
		for i := 0; i < blk.NumSuccs(); i++ {
			succ := blk.Succ(i)
			if len(info.phiOuts[succ]) == 0 {
				continue
			}

			fmt.Printf("    <- %s: ", succ)
			for val := range info.phiOuts[succ] {
				fmt.Print(val.IDString())
				if val.Reg != reg.None {
					fmt.Printf(":%s", val.Reg)
				}
				if val.String() != val.IDString() {
					fmt.Printf(":%s", val.String())
				}
				fmt.Print(" ")
			}
			fmt.Println()
		}

		fmt.Println("}}}------------")
	}
	fmt.Printf("Live through calls: %v\n", ra.liveThroughCalls)
}

func reverseIRSuccessorSort(block *ir.Block, list []*ir.Block, visited map[*ir.Block]bool) []*ir.Block {
	visited[block] = true

	for i := block.NumSuccs() - 1; i >= 0; i-- {
		succ := block.Succ(i)
		if !visited[succ] {
			list = reverseIRSuccessorSort(succ, list, visited)
		}
	}

	return append(list, block)
}

func (ra *RegAlloc) scanUsage(blk *ir.Block) bool {
	info := &ra.blockInfo[blk.ID()]

	oldLiveIns := info.liveIns

	info.blkKills = make(map[*ir.Value]bool)
	info.kills = make(map[*ir.Value][]*ir.Value)
	info.liveIns = make(map[*ir.Value]bool)
	info.liveOuts = make(map[*ir.Value]bool)
	info.phiIns = make(map[*ir.Block]map[*ir.Value]bool)
	info.phiOuts = make(map[*ir.Block]map[*ir.Value]bool)

	// for each successor block
	for i := 0; i < blk.NumSuccs(); i++ {
		succ := blk.Succ(i)
		sinfo := &ra.blockInfo[succ.ID()]

		// copy the live ins of successors to the live outs of this block
		for val := range sinfo.liveIns {
			info.liveIns[val] = true
			info.liveOuts[val] = true
		}

		// for each successor phi
		for j := 0; j < succ.NumInstrs(); j++ {
			val := succ.Instr(j)
			if val.Op != op.Phi {
				break
			}

			// find the index of blk in successor's pred list
			index := -1
			for k := 0; k < succ.NumPreds(); k++ {
				if succ.Pred(k) == blk {
					index = k
					break
				}
			}

			// take that arg and mark it as live in (for now)
			arg := val.Arg(index)

			if info.phiOuts[succ] == nil {
				info.phiOuts[succ] = make(map[*ir.Value]bool)
			}
			info.phiOuts[succ][arg] = true

			info.liveIns[arg] = true
		}
	}

	// for each block control
	for i := 0; i < blk.NumControls(); i++ {
		ctrl := blk.Control(i)

		if ctrl.NeedsReg() {
			info.liveIns[ctrl] = true

			if !info.liveOuts[ctrl] {
				info.blkKills[ctrl] = true
			}
		}
	}

	// for each instruction in reverse order
	for i := blk.NumInstrs() - 1; i >= 0; i-- {
		def := blk.Instr(i)

		// mark output as being seen
		delete(info.liveIns, def)

		// if this is a call, mark all "live" values during the call
		// as being live through it
		if def.Op == op.Call {
			for val := range info.liveIns {
				ra.liveThroughCalls[val] = true
			}
		}

		// for each arg
		for j := 0; j < def.NumArgs(); j++ {
			arg := def.Arg(j)

			if !arg.NeedsReg() {
				continue
			}

			// if arg is not live out and it's the first sighting
			// then mark it as killed
			if !info.liveOuts[arg] && !info.liveIns[arg] {
				info.kills[def] = append(info.kills[def], arg)
			}

			// don't mark phis as live-in
			if def.Op == op.Phi {
				pred := blk.Pred(j)
				if info.phiIns[pred] == nil {
					info.phiIns[pred] = make(map[*ir.Value]bool)
				}
				info.phiIns[pred][arg] = true
				continue
			}

			// mark it as live in
			info.liveIns[arg] = true
		}
	}

	if oldLiveIns == nil {
		return true
	}

	if len(oldLiveIns) != len(info.liveIns) {
		return true
	}

	for v := range info.liveIns {
		if !oldLiveIns[v] {
			return true
		}
	}

	for v := range oldLiveIns {
		if !info.liveIns[v] {
			return true
		}
	}

	return false
}
