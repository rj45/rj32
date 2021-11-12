// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"fmt"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (ra *regAlloc) colour() {
	ra.Func.Blocks()[0].VisitSuccessors(func(blk *ir.Block) bool {
		info := &ra.blockInfo[blk.ID()]
		var used reg.Reg

		info.regValues = make(map[reg.Reg]*ir.Value)
		for val := range info.liveIns {
			used |= val.Reg
			info.regValues[val.Reg] = val
		}

		for i := 0; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			if val.Op != op.Phi {
				for i := 0; i < val.NumArgs(); i++ {
					arg := val.Arg(i)
					if !val.Op.IsConst() && arg.Reg != reg.None && info.regValues[arg.Reg] != arg {
						log.Panicf("Attempted to read %s from reg %s, but contained %s! %s", arg.IDString(), arg.Reg, info.regValues[arg.Reg], val.LongString())
					}
				}
			}

			for _, val := range info.kills[val] {
				used &^= val.Reg
				info.regValues[val.Reg] = nil
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
		for _, v := range ra.affinities[val] {
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
	if info.liveOuts[val] && ra.Func.NumCalls > 0 {
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
