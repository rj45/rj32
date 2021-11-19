// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"fmt"
	"log"
	"os"

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
	ra.liveThroughCalls = make(map[*ir.Value]bool)

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
		info.phiIns = make(map[*ir.Block]map[*ir.Value]bool)
		info.phiOuts = make(map[*ir.Block]map[*ir.Value]bool)
		info.regValues = make(map[reg.Reg]*ir.Value)
	}

	// entry.VisitSuccessors(ra.scanUsage)

	// var list []*ir.Block

	// entry.VisitSuccessors(func(b *ir.Block) bool {
	// 	list = append(list, b)
	// 	return true
	// })

	// order blocks by reverse succession
	list := reverseIRSuccessorSort(ra.Func.Blocks()[0], nil, make(map[*ir.Block]bool))

	// reverse it to get succession ordering
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	visited := make(map[*ir.Block]bool)
	for i := len(list) - 1; i >= 0; i-- {
		ra.scanUsage2(list[i], list, visited)
	}

	// ra.scanVisit(entry, make(map[ir.ID]bool))

	dot, _ := os.Create(ra.Func.Name + ".dot")
	defer dot.Close()

	fmt.Fprintln(dot, "digraph G {")
	fmt.Fprintln(dot, "labeljust=l;")
	fmt.Fprintln(dot, "node [shape=record, fontname=\"Noto Mono\", labeljust=l];")

	for _, blk := range list {
		info := &ra.blockInfo[blk.ID()]

		for i := 0; i < blk.NumPreds(); i++ {
			pred := blk.Pred(i)
			pinfo := &ra.blockInfo[pred.ID()]
			outs := maptolist(pinfo.liveOuts) + " - " + maptolist(pinfo.phiOuts[blk])
			ins := maptolist(info.liveIns) + " - " + maptolist(info.phiIns[pred])
			fmt.Fprintf(dot, "%s -> %s [headlabel=%q, taillabel=%q];\n", pred, blk, outs, ins)
		}

		liveInKills := ""
		label := fmt.Sprintf("{<%s> %s:\\l", blk, blk)
		for i := 0; i < blk.NumInstrs(); i++ {
			instr := blk.Instr(i)
			label += " | "

			kills := ""

			for i, kill := range info.kills[instr] {
				if i != 0 {
					kills += " "
				}
				kills += kill.IDString()
				// if info.liveIns[kill] || info.phiIns[kill] {
				// 	liveInKills += fmt.Sprintf("%s:%s -> %s:%s;\n", kill.Block(), kill.IDString(), blk, instr.IDString())
				// }
			}

			label += fmt.Sprintf("<%s> %s [%s]\\l", instr.IDString(), instr.ShortString(), kills)

		}
		label += "}"

		fmt.Fprintf(dot, "%s [label=\"%s\"];\n", blk, label)
		fmt.Fprintln(dot, liveInKills)

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

	fmt.Fprintln(dot, "}")
}

func maptolist(l map[*ir.Value]bool) string {
	ret := ""
	for v := range l {
		ret += v.String()
		ret += " "
	}
	return ret
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

func (ra *RegAlloc) scanUsage2(blk *ir.Block, list []*ir.Block, visited map[*ir.Block]bool) bool {
	info := &ra.blockInfo[blk.ID()]

	// todo:
	// - make sure phi copies and phis are handled properly
	// - check to make sure the liveIns of loops look correct (not so sure)

	visited[blk] = true

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

			// delete(info.liveIns, arg)

			// not sure if this needs to be live out since it's linked
			// to a PhiCopy pegged to the phi
			// info.liveOuts[val] = true // ?
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

		if def.NeedsReg() {
			ra.trackAffinities(def, blk)
		}

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

	// find any loops beginning at this block
	// for each predecessor block
	for i := 0; i < blk.NumPreds(); i++ {
		pred := blk.Pred(i)

		// if we already visited it, this is probably a loop
		if visited[pred] {
			loop := blk.FindPathTo(func(b *ir.Block) bool { return b == pred })

			if loop[len(loop)-1] != pred {
				loop = append(loop, pred)
			}

			log.Println("Loop!", loop, info.liveIns)

			// for each block in loop
			for _, lblk := range loop {
				linfo := &ra.blockInfo[lblk.ID()]

				// for each value live at the start of the loop
				for val := range info.liveIns {
					// make value live through the block, except the last block
					linfo.liveOuts[val] = true
					linfo.liveIns[val] = true
					delete(linfo.blkKills, val)
				}

				// filter kills list to not include any blk.liveIns
				for kill, kills := range linfo.kills {
					var nkills []*ir.Value
					for _, k := range kills {
						if !info.liveIns[k] {
							nkills = append(nkills, k)
						}
					}
					linfo.kills[kill] = nkills
					if len(nkills) == 0 {
						delete(linfo.kills, kill)
					}
				}
			}
		}
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
			if info.liveOuts[arg] || info.blkKills[arg] {
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
