// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

var dotlive = flag.Bool("dot", false, "write .dot files for reg alloc debugging")

type RegAlloc struct {
	Func *ir.Func

	usedRegs reg.Reg

	blockInfo []blockInfo

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
		ra.writeExplodedDotFile(ra.Func.Blocks())
	}

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

func (ra *RegAlloc) writeExplodedDotFile(list []*ir.Block) {
	dot, _ := os.Create(ra.Func.Name + ".ex.dot")
	defer dot.Close()

	fmt.Fprintln(dot, "digraph G {")
	fmt.Fprintln(dot, "node [fontname=\"Noto Mono\", shape=rect];")

	for _, blk := range list {
		info := &ra.blockInfo[blk.ID()]

		fmt.Fprintf(dot, "subgraph cluster_%s {\n", blk.String())

		fmt.Fprintf(dot, "label=\"%s\";\n", blk.String())
		fmt.Fprintln(dot, "labeljust=l;")
		fmt.Fprintln(dot, "color=black;")

		fmt.Fprintln(dot, "node [shape=plaintext];")

		srcs := make(map[*ir.Value]string)
		var lastStr string

		for i := 0; i < blk.NumPreds(); i++ {
			pred := blk.Pred(i)
			inname := fmt.Sprintf("in_%s_%s", blk, pred)
			ins := maptorecordsrc(inname, info.liveIns, srcs)
			if len(info.phiIns[pred]) > 0 {
				ins += maptorecordsrc(inname, info.phiIns[pred], srcs)
			}
			fmt.Fprintf(dot, "%s [label=<<table border=\"0\" cellborder=\"1\" cellspacing=\"0\"><tr><td port=\"in\">%s in</td>%s</tr></table>>];\n", inname, blk, ins)
			fmt.Fprintf(dot, "out_%s_%s:out:s -> %s:in:n;\n", pred, blk, inname)
			lastStr = inname
		}

		if blk.NumInstrs() > 0 && lastStr != "" {
			// force the instructions to be in order
			fmt.Fprintf(dot, "%s -> %s_%s [weight=100, style=invis];\n", lastStr, blk, blk.Instr(0).IDString())
		}

		for i := 0; i < blk.NumInstrs(); i++ {
			instr := blk.Instr(i)

			edges := ""
			name := fmt.Sprintf("%s_%s", blk, instr.IDString())
			lastStr = name

			label := ""
			if instr.NeedsReg() {
				label += fmt.Sprintf("<td port=\"%s\">%s</td>", instr.IDString(), valname(instr))
			}

			label += fmt.Sprintf("<td>%s</td>", instr.Op)
			for j := 0; j < instr.NumArgs(); j++ {
				arg := instr.Arg(j)

				killed := false
				for _, kill := range info.kills[instr] {
					if kill == arg {
						killed = true
						break
					}
				}

				arglabel := valname(arg)
				if killed {
					arglabel = "[" + arglabel + "]"
				}

				label += fmt.Sprintf("<td port=\"%s\">%s</td>", arg.IDString(), arglabel)

				if arg.NeedsReg() {
					edges += fmt.Sprintf("%s:s -> %s:%s:n;\n", srcs[arg], name, arg.IDString())
				}

				// chain arrows through uses
				srcs[arg] = fmt.Sprintf("%s:%s", name, arg.IDString())
			}

			srcs[instr] = fmt.Sprintf("%s:%s", name, instr.IDString())

			fmt.Fprintf(dot, "%s [label=<<table border=\"0\" cellborder=\"1\" cellspacing=\"0\"><tr>%s</tr></table>>];\n", name, label)

			fmt.Fprint(dot, edges)

			if i < blk.NumInstrs()-1 {
				// force the instructions to be in order
				fmt.Fprintf(dot, "%s -> %s_%s [weight=100, style=invis];\n", name, blk, blk.Instr(i+1).IDString())
			}
		}

		// emit block control instruction
		{
			name := fmt.Sprintf("%s_ctrl", blk)
			label := ""

			var edges string
			label += fmt.Sprintf("<td>%s</td>", blk.Op)

			for i := 0; i < blk.NumControls(); i++ {
				arg := blk.Control(i)

				arglabel := valname(arg)
				if info.blkKills[arg] {
					arglabel = "[" + arglabel + "]"
				}

				label += fmt.Sprintf("<td port=\"%s\">%s</td>", arg.IDString(), arglabel)

				if arg.NeedsReg() {
					edges += fmt.Sprintf("%s:s -> %s:%s:n;\n", srcs[arg], name, arg.IDString())
				}

				// chain arrows through uses
				srcs[arg] = fmt.Sprintf("%s:%s", name, arg.IDString())
			}

			fmt.Fprintf(dot, "%s [label=<<table border=\"0\" cellborder=\"1\" cellspacing=\"0\"><tr>%s</tr></table>>];\n", name, label)
			fmt.Fprint(dot, edges)

			if blk.NumInstrs() > 0 {
				// force the instructions to be in order
				fmt.Fprintf(dot, "%s -> %s [weight=100, style=invis];\n", lastStr, name)
			}
			lastStr = name
		}

		for i := 0; i < blk.NumSuccs(); i++ {
			succ := blk.Succ(i)
			sinfo := &ra.blockInfo[succ.ID()]

			outs := maptorecord(sinfo.liveIns)
			if len(info.phiOuts[succ]) > 0 {
				outs += maptorecord(info.phiOuts[succ])
			}

			fmt.Fprintf(dot, "out_%s_%s [label=<<table border=\"0\" cellborder=\"1\" cellspacing=\"0\"><tr><td port=\"out\">%s out</td>%s</tr></table>>];\n", blk, succ, blk, outs)
			for v := range sinfo.liveIns {
				fmt.Fprintf(dot, "%s -> out_%s_%s:%s;\n", srcs[v], blk, succ, v.IDString())
			}
			for v := range info.phiOuts[succ] {
				fmt.Fprintf(dot, "%s -> out_%s_%s:%s;\n", srcs[v], blk, succ, v.IDString())
			}

			// force the instructions to be in order
			fmt.Fprintf(dot, "%s -> out_%s_%s [weight=100, style=invis];\n", lastStr, blk, succ)
		}

		fmt.Fprintln(dot, "}")
	}

	fmt.Fprintln(dot, "}")
}

func valname(val *ir.Value) string {
	if val.NeedsReg() {
		return fmt.Sprintf("%s:%s", val.IDString(), val.Reg)
	}
	return val.String()
}

func maptorecord(l map[*ir.Value]bool) string {
	ret := ""
	for v := range l {
		ret += fmt.Sprintf("<td port=\"%s\">%s:%s</td>", v.IDString(), v.IDString(), v.Reg.String())
	}
	return ret
}

func maptorecordsrc(prefix string, l map[*ir.Value]bool, src map[*ir.Value]string) string {
	ret := ""
	for v := range l {
		ret += fmt.Sprintf("<td port=\"%s\">%s:%s</td>", v.IDString(), v.IDString(), v.Reg.String())
		src[v] = fmt.Sprintf("%s:%s", prefix, v.IDString())
	}
	return ret
}
