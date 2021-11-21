// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

var dotlive = flag.Bool("dotlive", false, "write .dot files for liveness debugging")
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
		info.blkKills = make(map[*ir.Value]bool)
		info.kills = make(map[*ir.Value][]*ir.Value)
		info.liveIns = make(map[*ir.Value]bool)
		info.liveOuts = make(map[*ir.Value]bool)
		info.phiIns = make(map[*ir.Block]map[*ir.Value]bool)
		info.phiOuts = make(map[*ir.Block]map[*ir.Value]bool)
		info.regValues = make(map[reg.Reg]*ir.Value)
	}

	// order blocks by reverse succession
	list := reverseIRSuccessorSort(ra.Func.Blocks()[0], nil, make(map[*ir.Block]bool))

	// reverse it to get succession ordering
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}

	visited := make(map[*ir.Block]bool)
	for i := len(list) - 1; i >= 0; i-- {
		ra.scanUsage(list[i], list, visited)
	}

	if *dotlive {
		ra.writeLivenessDotFile(list)
	}
}

func (ra *RegAlloc) writeLivenessDotFile(list []*ir.Block) {
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
		label := fmt.Sprintf("%s:\\l", blk)
		for i := 0; i < blk.NumInstrs(); i++ {
			instr := blk.Instr(i)

			kills := ""

			for i, kill := range info.kills[instr] {
				if i != 0 {
					kills += " "
				}
				kills += kill.IDString()

			}

			label += fmt.Sprintf("%s [%s]\\l", instr.ShortString(), kills)
		}

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

func (ra *RegAlloc) scanUsage(blk *ir.Block, list []*ir.Block, visited map[*ir.Block]bool) bool {
	info := &ra.blockInfo[blk.ID()]

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

			if *debugLiveness {
				log.Println("Loop!", loop, info.liveIns)
			}

			// for each block in loop
			for _, lblk := range loop {
				linfo := &ra.blockInfo[lblk.ID()]

				hasCall := false
				for j := 0; j < lblk.NumInstrs(); j++ {
					if lblk.Instr(j).Op == op.Call {
						hasCall = true
					}
				}

				// for each value live at the start of the loop
				for val := range info.liveIns {
					// make value live through the block, except the last block
					linfo.liveOuts[val] = true
					linfo.liveIns[val] = true
					delete(linfo.blkKills, val)

					if hasCall {
						ra.liveThroughCalls[val] = true
					}
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
