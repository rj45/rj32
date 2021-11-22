// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

var dotlive = flag.Bool("dot", false, "write .dot files for reg alloc debugging")

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

	if *dotlive {
		ra.writeDotFile(ra.Func.Blocks())
	}

	log.Println("Copies eliminated:", ra.copiesEliminated, "out of potentially:", ra.potentialCopiesEliminated)

	return ra.usedRegs
}

func (ra *RegAlloc) Verify() {
	ra.verify(false)
	ra.verify(true)
}

func (ra *RegAlloc) writeDotFile(list []*ir.Block) {
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

		first := true
		kills := ""
		for kill := range info.blkKills {
			if !first {
				kills += " "
			}
			first = false
			kills += kill.IDString()
		}

		label += fmt.Sprintf("%s [%s]\\l", blk.OpString(), kills)

		label = strings.ReplaceAll(label, "\"", "\\\"")

		fmt.Fprintf(dot, "%s [label=\"%s\"];\n", blk, label)
		fmt.Fprintln(dot, liveInKills)
	}

	fmt.Fprintln(dot, "}")
}

func maptolist(l map[*ir.Value]bool) string {
	ret := ""
	first := true
	for v := range l {
		if !first {
			ret += " "
		}
		first = false
		ret += fmt.Sprintf("%s:%s", v.IDString(), v.Reg.String())
	}
	return ret
}
