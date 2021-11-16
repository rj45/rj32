// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"fmt"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (ra *RegAlloc) liveScan() {
	fmt.Println("------------")
	fmt.Println(ra.Func.LongString())
	fmt.Println("------------")

	ra.affinities = make(map[*ir.Value][]*ir.Value)
	ra.blockInfo = make([]blockInfo, ra.Func.BlockIDCount())

	entry := ra.Func.Blocks()[0]

	if entry.Idom != nil || entry.NumPreds() != 0 {
		log.Fatalln("expecting entry block to be top of the tree")
	}

	for _, blk := range ra.Func.Blocks() {
		info := &ra.blockInfo[blk.ID()]
		info.blkKills = make(map[*ir.Value]bool)
		info.kills = make(map[*ir.Value][]*ir.Value)
		info.liveIns = make(map[*ir.Value]bool)
		info.liveOuts = make(map[*ir.Value]bool)
		info.regValues = make(map[reg.Reg]*ir.Value)
	}

	entry.VisitSuccessors(ra.scanUsage)

	// ra.scanVisit(entry, make(map[ir.ID]bool))

	for _, blk := range ra.Func.Blocks() {

		fmt.Println("{{{------------")

		fmt.Println(blk.LongString())

		info := ra.blockInfo[blk.ID()]

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

		// pretty.Println(ra.blockInfo[blk.ID()])

		fmt.Println("}}}------------")
	}
}

func (ra *RegAlloc) scanUsage(blk *ir.Block) bool {
	for i := 0; i < blk.NumInstrs(); i++ {
		def := blk.Instr(i)

		if !def.NeedsReg() {
			continue
		}

		paths := def.FindUsageSuccessorPaths()

		for _, path := range paths {
			killBlk := path[len(path)-1]
			kbInfo := ra.blockInfo[killBlk.ID()]

			found := false
			for i := 0; i < def.NumBlockUses(); i++ {
				if def.BlockUse(i) == killBlk {
					found = true
					kbInfo.blkKills[def] = true
					break
				}
			}

			var kill *ir.Value
			if !found {
				last := -1
				for i := 0; i < def.NumArgUses(); i++ {
					arg := def.ArgUse(i)
					if arg.Block() == killBlk {
						if arg.Index() > last {
							last = arg.Index()
							kill = arg
						}
					}
				}

				kbInfo.kills[kill] = append(kbInfo.kills[kill], def)
			}

			for i, pblk := range path {
				pinfo := ra.blockInfo[pblk.ID()]

				if i != 0 {
					pinfo.liveIns[def] = true
				}

				if i != len(path)-1 {
					pinfo.liveOuts[def] = true
				}
			}

			// Phis are special, they act as if the value is parallel copied on the
			// edge between blocks in the CFG:
			// - The value is killed just after the block
			//   - so it ends up in the blkKills of the pred block
			//   - it is not live-in to the last block
			//   - I think it's not live out of the pred block either
			if kill != nil && kill.Op == op.Phi {
				var pred *ir.Block
				for i := 0; i < kill.NumArgs(); i++ {
					if kill.Arg(i) == def {
						pred = kill.Block().Pred(i)
					}
				}
				pinfo := ra.blockInfo[pred.ID()]
				pinfo.blkKills[def] = true
				if pinfo.phiOuts == nil {
					pinfo.phiOuts = make(map[*ir.Value]bool)
				}
				pinfo.phiOuts[def] = true
				delete(kbInfo.liveIns, def)
				delete(pinfo.liveOuts, def) // ?

				log.Println("killed phi", kill, pred, killBlk, kbInfo.liveIns, pinfo.liveOuts)
			}
		}

		ra.trackAffinities(def, blk)
	}

	return true
}

// keep track of affinities to help with copy elimination
func (ra *RegAlloc) trackAffinities(instr *ir.Value, blk *ir.Block) {
	info := ra.blockInfo[blk.ID()]
	if instr.Op.IsCopy() && instr.NumArgs() > 0 {
		for i := 0; i < instr.NumArgs(); i++ {
			arg := instr.Arg(i)
			// make sure arg doesn't escape
			if info.liveOuts[arg] || info.phiOuts[arg] || info.blkKills[arg] {
				continue
			}

			// make sure this is marked as the last use of this arg
			found := false
			for _, k := range info.kills[instr] {
				if k == arg {
					found = true
				}
			}
			if !found {
				continue
			}

			if !arg.Op.IsConst() {
				ra.affinities[instr] = append(ra.affinities[instr], arg)
				ra.affinities[arg] = append(ra.affinities[arg], instr)
			}
		}
	}
}

// func (ra *regAlloc) scanVisit(blk *ir.Block, visited map[ir.ID]bool) {
// 	// track whether it's visited
// 	visited[blk.ID()] = true

// 	// visit all children first, else block first
// 	for i := blk.NumSuccs() - 1; i >= 0; i-- {
// 		succ := blk.Succ(i)
// 		if !visited[succ.ID()] {
// 			ra.scanVisit(succ, visited)
// 		}
// 		// TODO: else do we need to copy anything into the already visited block?
// 	}

// 	// setup the block info
// 	info := &ra.blockInfo[blk.ID()]
// 	if info.kills == nil {
// 		info.kills = make(map[*ir.Value][]*ir.Value)
// 	}
// 	if info.liveIns == nil {
// 		info.liveIns = make(map[*ir.Value]bool)
// 	}

// 	if blk.Op == op.Return {
// 		// for return blocks, the controls are live-outs
// 		for i := 0; i < blk.NumControls(); i++ {
// 			if info.liveOuts == nil {
// 				info.liveOuts = make(map[*ir.Value]bool)
// 			}
// 			info.liveOuts[blk.Control(i)] = true
// 		}
// 	} else {
// 		// make sure block controls count as killed values
// 		for i := 0; i < blk.NumControls(); i++ {
// 			if !info.liveOuts[blk.Control(i)] {
// 				if info.blkKills == nil {
// 					info.blkKills = make(map[*ir.Value]bool)
// 				}
// 				info.blkKills[blk.Control(i)] = true
// 			}
// 		}
// 	}

// 	// initially copy any live-outs to live-ins
// 	for out := range info.liveOuts {
// 		info.liveIns[out] = true
// 	}

// 	// also copy phi-outs
// 	for out := range info.phiOuts {
// 		info.liveIns[out] = true
// 	}

// 	// for each instruction in the block, from last to first
// 	for i := blk.NumInstrs() - 1; i >= 0; i-- {
// 		instr := blk.Instr(i)

// 		// keep track of affinities to help with copy elimination
// 		if instr.Op == op.Copy || instr.Op == op.Phi {
// 			if instr.Reg.CanAffinity() {
// 				ra.affinities[instr] = append(ra.affinities[instr], instr.Arg(0))
// 				for j := 0; j < instr.NumArgs(); j++ {
// 					arg := instr.Arg(j)
// 					ra.affinities[arg] = append(ra.affinities[arg], instr)
// 				}
// 			}
// 		}

// 		// try to also assign the same register to the first arg if it's clobbered
// 		if instr.Op.ClobbersArg() {
// 			ra.affinities[instr] = append(ra.affinities[instr], instr.Arg(0))
// 			ra.affinities[instr.Arg(0)] = append(ra.affinities[instr.Arg(0)], instr)
// 		}

// 		// handle the definition
// 		{
// 			if info.liveIns[instr] {
// 				// no longer a live in
// 				delete(info.liveIns, instr)
// 			}
// 		}

// 		// phi are treated specially
// 		if instr.Op == op.Phi {
// 			for i := 0; i < instr.NumArgs(); i++ {
// 				arg := instr.Arg(i)
// 				if arg.Op.IsConst() {
// 					continue
// 				}

// 				// find the pred block
// 				pred := blk.Pred(i)

// 				// mark the pred block as having the phiOut
// 				pinfo := &ra.blockInfo[pred.ID()]
// 				if pinfo.phiOuts == nil {
// 					pinfo.phiOuts = make(map[*ir.Value]bool)
// 				}
// 				pinfo.phiOuts[arg] = true

// 				// not marking the live-in because it doesn't come in
// 				// from all blocks, just some. Marking as phiIn instead
// 				if info.phiIns == nil {
// 					info.phiIns = make(map[*ir.Value]bool)
// 				}
// 				info.phiIns[arg] = true
// 			}
// 			continue
// 		}

// 		// for each value this instr reads
// 		for i := 0; i < instr.NumArgs(); i++ {
// 			arg := instr.Arg(i)
// 			if arg.Op.IsConst() {
// 				continue
// 			}

// 			// is this the first read?
// 			if !info.liveOuts[arg] && !info.phiOuts[arg] && !info.liveIns[arg] && !info.blkKills[arg] {
// 				info.kills[instr] = append(info.kills[instr], arg)
// 				info.liveIns[arg] = true
// 			}
// 		}
// 	}

// 	// copy the live-ins to the live-outs of pred blocks
// 	for i := 0; i < blk.NumPreds(); i++ {
// 		pred := blk.Pred(i)
// 		pinfo := &ra.blockInfo[pred.ID()]
// 		if pinfo.liveOuts == nil {
// 			pinfo.liveOuts = make(map[*ir.Value]bool)
// 		}
// 		for id := range info.liveIns {
// 			pinfo.liveOuts[id] = true
// 		}
// 	}
// }
