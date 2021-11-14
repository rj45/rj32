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

		var unresolved map[*ir.Value]bool

		for val := range info.liveIns {
			if val.Reg == reg.None {
				if unresolved == nil {
					unresolved = make(map[*ir.Value]bool)
				}
				unresolved[val] = true
				continue
			}
			used |= val.Reg
			info.regValues[val.Reg] = val
		}

		for val := range unresolved {
			// need to guess at the register that will be assigned
			if ra.guessedRegs == nil {
				ra.guessedRegs = make(map[*ir.Value]bool)
			}
			ra.guessedRegs[val] = true

			otherUsed := used
			otherInfo := &ra.blockInfo[val.Block().ID()]

			for val := range info.liveIns {
				if val.Reg != reg.None {
					otherUsed |= val.Reg
				}
			}

			val.Reg = ra.chooseReg(otherInfo, val, otherUsed)

			if val.Reg != reg.None {
				used |= val.Reg
				info.regValues[val.Reg] = val
			}
		}

		for i := 0; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			for _, kill := range info.kills[val] {
				used &^= kill.Reg
				delete(info.regValues, kill.Reg)
			}

			// stores and some calls don't need a reg
			if !val.NeedsReg() {
				continue
			}

			if ra.guessedRegs[val] {
				chosen := ra.chooseReg(info, val, used)
				if chosen != val.Reg && (val.Reg&used) != 0 {
					val.Reg = chosen
				} else {
					// lucky guess, no need to follow up
					delete(ra.guessedRegs, val)
				}
			}

			if val.Reg == reg.None {
				val.Reg = ra.chooseReg(info, val, used)
			}

			if val.Reg == reg.None {
				log.Println(blk.LongString())
				log.Fatal("Ran out of registers, spilling not implemented")
			}

			used |= val.Reg
			ra.usedRegs |= val.Reg
			info.regValues[val.Reg] = val
		}

		fmt.Println(blk.LongString())

		return true
	})
}

func (ra *regAlloc) chooseReg(info *blockInfo, val *ir.Value, used reg.Reg) reg.Reg {
	var chosen reg.Reg
	if len(ra.affinities[val]) > 0 {
		votes := make(map[reg.Reg]int)
		if val.Reg != reg.None && (used&val.Reg) == 0 {
			votes[val.Reg]++
		}
		for _, v := range ra.affinities[val] {
			notInUse := (used&v.Reg) == 0 || (info.regValues[v.Reg] == v && val.Op.IsCopy())
			if v.Reg != reg.None && notInUse && v.Reg.CanAffinity() {
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
		for reg, votes := range votes {
			if votes == max {
				if reg.IsSavedReg() {
					chosen = reg
				}
			}
		}
		if chosen != reg.None {
			log.Println("affinity chosen", val, chosen, ra.affinities[val], votes)
			return chosen
		} else if len(ra.affinities[val]) > 0 {
			log.Println("affinity failure:", val, ra.affinities[val], used, votes)
		}
	}

	// if val.Op.ClobbersArg() && (used&val.Arg(0).Reg) == 0 {
	// 	return val.Arg(0).Reg
	// }

	sets := [][]reg.Reg{reg.TempRegs, reg.ArgRegs, reg.SavedRegs}

	escapes := info.liveOuts[val] || info.phiOuts[val]
	if escapes && ra.Func.NumCalls > 0 {
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
