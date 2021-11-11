package regalloc

import (
	"fmt"
	"log"

	"github.com/kr/pretty"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func (ra *regAlloc) liveScan() {
	fmt.Println("------------")
	fmt.Println(ra.Func.LongString())
	fmt.Println("------------")

	ra.affinities = make(map[ir.ID][]*ir.Value)
	ra.blockInfo = make([]blockInfo, ra.Func.BlockIDCount())

	if ra.Func.Blocks[0].Idom != nil {
		log.Fatalln("expecting first block to be top of dom tree")
	}

	// for each unvisited block
	visited := make(map[ir.ID]bool)

	ra.scanVisit(ra.Func.Blocks[0], visited)

	for _, blk := range ra.Func.Blocks {

		fmt.Println("{{{------------")

		fmt.Println(blk.LongString())

		pretty.Println(ra.blockInfo[blk.ID])

		fmt.Println("}}}------------")

	}
}

func (ra *regAlloc) scanVisit(blk *ir.Block, visited map[ir.ID]bool) {
	// track whether it's visited
	visited[blk.ID] = true

	// visit all children first, else block first
	for i := len(blk.Succs) - 1; i >= 0; i-- {
		succ := blk.Succs[i].Block
		if !visited[succ.ID] {
			ra.scanVisit(succ, visited)
		}
		// TODO: else do we need to copy anything into the already visited block?
	}

	// setup the block info
	info := &ra.blockInfo[blk.ID]
	info.kills = make(map[ir.ID][]ir.ID)
	info.liveIns = make(map[ir.ID]bool)

	if blk.Op == op.Return {
		// for return blocks, the controls are live-outs
		for _, arg := range blk.Controls {
			info.liveOuts[arg.ID] = true
		}
	} else {
		// make sure block controls count as killed values
		for _, arg := range blk.Controls {
			if !info.liveOuts[arg.ID] {
				if info.blkKills == nil {
					info.blkKills = make(map[ir.ID]bool)
				}
				info.blkKills[arg.ID] = true
			}
		}
	}

	// initially copy any live-outs to live-ins
	for out := range info.liveOuts {
		info.liveIns[out] = true
	}

	// also copy phi-outs
	for out := range info.phiOuts {
		info.liveIns[out] = true
	}

	// for each instruction in the block, from last to first
	for i := len(blk.Instrs) - 1; i >= 0; i-- {
		instr := blk.Instrs[i]

		// keep track of affinities to help with copy elimination
		if instr.Op == op.Copy || instr.Op == op.Phi {
			if instr.Reg.CanAffinity() {
				ra.affinities[instr.ID] = append(ra.affinities[instr.ID], instr.Args[0])
				for _, arg := range instr.Args {
					ra.affinities[arg.ID] = append(ra.affinities[arg.ID], instr)
				}
			}
		}

		// try to also assign the same register to the first arg if it's clobbered
		if instr.Op.ClobbersArg() {
			ra.affinities[instr.ID] = append(ra.affinities[instr.ID], instr.Args[0])
			ra.affinities[instr.Args[0].ID] = append(ra.affinities[instr.Args[0].ID], instr)
		}

		// handle the definition
		{
			if info.liveIns[instr.ID] {
				// no longer a live in
				delete(info.liveIns, instr.ID)
			}
		}

		// phi are treated specially
		if instr.Op == op.Phi {
			for i, arg := range instr.Args {
				if arg.Op.IsConst() {
					continue
				}

				// find the pred block
				var pred *ir.Block
				for _, ref := range blk.Preds {
					if ref.Index == i {
						pred = ref.Block
					}
				}

				// mark the pred block as having the phiOut
				pinfo := &ra.blockInfo[pred.ID]
				if pinfo.phiOuts == nil {
					pinfo.phiOuts = make(map[ir.ID]bool)
				}
				pinfo.phiOuts[arg.ID] = true

				// not marking the live-in because it doesn't come in
				// from all blocks, just some. Marking as phiIn instead
				if info.phiIns == nil {
					info.phiIns = make(map[ir.ID]bool)
				}
				info.phiIns[arg.ID] = true
			}
			continue
		}

		// for each value this instr reads
		for _, arg := range instr.Args {
			if arg.Op.IsConst() {
				continue
			}

			// is this the first read?
			if !info.liveOuts[arg.ID] && !info.phiOuts[arg.ID] && !info.liveIns[arg.ID] && !info.blkKills[arg.ID] {
				info.kills[instr.ID] = append(info.kills[instr.ID], arg.ID)
				info.liveIns[arg.ID] = true
			}
		}
	}

	// copy the live-ins to the live-outs of pred blocks
	for _, ref := range blk.Preds {
		pred := ref.Block

		pinfo := &ra.blockInfo[pred.ID]
		if pinfo.liveOuts == nil {
			pinfo.liveOuts = make(map[ir.ID]bool)
		}
		for id := range info.liveIns {
			pinfo.liveOuts[id] = true
		}
	}
}
