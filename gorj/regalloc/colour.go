// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"fmt"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (ra *regAlloc) colour() {
	ra.Func.Blocks()[0].VisitSuccessors(func(blk *ir.Block) bool {
		info := &ra.blockInfo[blk.ID()]
		var used reg.Reg

		for id := range info.liveIns {
			val := ra.Func.ValueForID(id)
			used |= val.Reg
		}

		for _, val := range blk.Instrs {
			for _, id := range info.kills[val.ID()] {
				free := ra.Func.ValueForID(id).Reg
				used &^= free
			}

			if val.Reg == reg.None {
				val.Reg = ra.chooseReg(info, val, used)
			}

			if val.Reg == reg.None {
				fmt.Println(blk.LongString())
				log.Fatal("Ran out of registers, spilling not implemented")
			}

			used |= val.Reg
			ra.usedRegs |= val.Reg
		}

		fmt.Println(blk.LongString())

		return true
	})
}

func (ra *regAlloc) chooseReg(info *blockInfo, val *ir.Value, used reg.Reg) reg.Reg {
	var chosen reg.Reg
	if len(ra.affinities[val.ID()]) > 0 {
		votes := make(map[reg.Reg]int)
		for _, v := range ra.affinities[val.ID()] {
			if v.Reg != reg.None && (used&v.Reg) == 0 {
				votes[v.Reg]++
			}
		}
		max := 0
		for reg, votes := range votes {
			if votes > max {
				max = votes
				chosen = reg
			}
		}
		if chosen != reg.None {
			return chosen
		}
	}

	sets := [][]reg.Reg{reg.TempRegs, reg.ArgRegs, reg.SavedRegs}
	if info.liveOuts[val.ID()] && ra.Func.NumCalls > 0 {
		sets = [][]reg.Reg{reg.SavedRegs, reg.TempRegs, reg.ArgRegs}
	}

	for _, set := range sets {
		for _, reg := range set {
			if (used & reg) == 0 {
				return reg
			}
		}
	}

	return reg.None
}

// func safeToUse(val *ir.Value, info *blockInfo, val *ir.Value, used reg.Reg) bool {
// 	if info.liveIns[val.ID()] &&
// }
