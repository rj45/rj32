// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package regalloc

import (
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (ra *RegAlloc) colour() {
	ra.wrongGuesses = make(map[*ir.Value]bool)

	ra.Func.Blocks()[0].VisitSuccessors(ra.allocateBlock)

	ra.guessedRegs = ra.wrongGuesses

	// reallocate blocks with guessed regs
	for len(ra.guessedRegs) > 0 {
		visited := map[*ir.Block]bool{}

		ra.wrongGuesses = make(map[*ir.Value]bool)

		for val := range ra.guessedRegs {
			blk := val.Block()
			if visited[blk] {
				continue
			}
			visited[blk] = true

			for i := 0; i < blk.NumInstrs(); i++ {
				val := blk.Instr(i)
				if !ra.guessedRegs[val] {
					val.Reg = reg.None
				}
			}

			ra.allocateBlock(blk)
		}

		ra.guessedRegs = ra.wrongGuesses
	}
}

func (ra *RegAlloc) allocateBlock(blk *ir.Block) bool {
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
		otherBlk := val.Block()
		otherInfo := &ra.blockInfo[otherBlk.ID()]

		for val := range info.liveIns {
			if val.Reg != reg.None {
				otherUsed |= val.Reg
			}
		}

		otherRegs := map[reg.Reg]*ir.Value{}
		for i := 0; i < otherBlk.NumInstrs(); i++ {
			otherVal := otherBlk.Instr(i)
			if otherVal == val {
				break
			}
			otherUsed = ra.allocateValue(otherInfo, otherVal, otherUsed, otherBlk, otherRegs)
		}

		val.Reg = ra.chooseReg(otherInfo, val, otherUsed)

		if val.Reg != reg.None {
			used |= val.Reg
			info.regValues[val.Reg] = val
		}
	}

	// first handle all the phis
	// these are parallel copies, so we do this in phases
	// phase 1: expire used vars
	i := 0
	for ; i < blk.NumInstrs(); i++ {
		phi := blk.Instr(i)
		if phi.Op != op.Phi {
			break
		}

		used = ra.killUsed(info, phi, used)
	}

	// phase 2, choose registers
	i = 0
	for ; i < blk.NumInstrs(); i++ {
		val := blk.Instr(i)
		if val.Op != op.Phi {
			break
		}

		if val.Reg == reg.None {
			val.Reg = ra.chooseReg(info, val, used)
		}

		if val.Reg == reg.None {
			log.Println(blk.LongString())
			log.Fatal("Ran out of registers, spilling not implemented")
		}
	}

	// phase 3, record register usage
	i = 0
	for ; i < blk.NumInstrs(); i++ {
		val := blk.Instr(i)
		if val.Op != op.Phi {
			break
		}

		used |= val.Reg
		ra.usedRegs |= val.Reg
		info.regValues[val.Reg] = val
	}

	for ; i < blk.NumInstrs(); i++ {
		val := blk.Instr(i)

		used = ra.allocateValue(info, val, used, blk, info.regValues)
		ra.usedRegs |= val.Reg
	}

	// reload all spills before the end of the block
	// for spilled := range info.spills {
	// 	var dontcare int
	// 	used = ra.reloadSpill(blk, -1, spilled, used, info, &dontcare)
	// }

	return true
}

func (ra *RegAlloc) allocateValue(info *blockInfo, val *ir.Value, used reg.Reg, blk *ir.Block, regValues map[reg.Reg]*ir.Value) reg.Reg {
	// used = ra.reloadSpilledArgs(val, used, info, &i)

	// if val.Op == op.Call {
	// 	used = ra.spillAllTempRegs(val, used, info, &i)
	// }
	// stores and some calls don't need a reg

	for _, kill := range info.kills[val] {
		used &^= kill.Reg
		delete(regValues, kill.Reg)
	}

	if !val.NeedsReg() {
		return used
	}

	if val.Reg == reg.None || ra.guessedRegs[val] {
		val.Reg = ra.chooseReg(info, val, used)
	}

	if val.Reg == reg.None {
		log.Println(blk.LongString())
		log.Fatal("Ran out of registers, spilling not implemented")
	}

	used |= val.Reg

	regValues[val.Reg] = val

	return used
}

func (*RegAlloc) killUsed(info *blockInfo, phi *ir.Value, used reg.Reg) reg.Reg {
	for _, kill := range info.kills[phi] {
		used &^= kill.Reg
		delete(info.regValues, kill.Reg)
	}
	return used
}

func (ra *RegAlloc) reloadSpilledArgs(val *ir.Value, used reg.Reg, info *blockInfo, index *int) reg.Reg {
	for j := 0; j < val.NumArgs(); j++ {
		arg := val.Arg(j)

		if _, spilled := info.spills[arg]; spilled {
			used = ra.reloadSpill(val.Block(), val.Index(), arg, used, info, index)
		}

		if repl, found := ra.spillReloads[arg]; found {
			val.ReplaceArg(j, repl)
			arg = repl
		}
	}
	return used
}

func (ra *RegAlloc) reloadSpill(blk *ir.Block, where int, arg *ir.Value, used reg.Reg, info *blockInfo, index *int) reg.Reg {
	fn := blk.Func()

	slot := info.spills[arg]

	// reload the spilled variable
	offset := int64(slot + fn.ArgSlots)
	load := fn.NewValue(
		op.Load, arg.Type, fn.FixedReg(reg.SP),
		fn.IntConst(offset))
	blk.InsertInstr(where, load)

	// make sure to increment past this so we don't get in a loop
	*index++

	// any future references to arg need to be replaced by the load
	ra.spillReloads[arg] = load

	delete(info.spills, arg)

	load.Reg = arg.Reg

	// used |= load.Reg
	// ra.usedRegs |= load.Reg
	// info.regValues[load.Reg] = arg

	return used
}

func (ra *RegAlloc) spillAllTempRegs(call *ir.Value, used reg.Reg, info *blockInfo, index *int) reg.Reg {
	blk := call.Block()
	fn := blk.Func()

	// spill all temp regs
	for _, tmp := range reg.TempRegs {
		if used&tmp != 0 {
			val := info.regValues[tmp]

			if _, alreadySpilled := info.spills[val]; alreadySpilled {
				continue
			}

			// find a free stack slot
			slot := -1
			if len(info.freeSlots) > 0 {
				slot = info.freeSlots[len(info.freeSlots)-1]
			} else {
				slot = ra.Func.SpillSlots
				ra.Func.SpillSlots++
			}

			// spill the value to the stack
			offset := int64(slot + fn.ArgSlots)
			blk.InsertInstr(call.Index(), fn.NewValue(
				op.Store, val.Type, fn.FixedReg(reg.SP),
				fn.IntConst(offset), val))

			// make sure to increment past this so we don't get in a loop
			*index++

			if info.spills == nil {
				info.spills = make(map[*ir.Value]int)
			}
			info.spills[val] = slot

			// used &^= val.Reg
			// delete(info.regValues, val.Reg)
		}
	}

	return used
}

func (ra *RegAlloc) chooseReg(info *blockInfo, val *ir.Value, used reg.Reg) reg.Reg {
	var chosen reg.Reg
	liveThroughCalls := ra.liveThroughCalls[val]

	// todo: not sure this is correct
	if ra.guessedRegs[val] {
		oldreg := val.Reg
		delete(ra.guessedRegs, val)
		chosen := ra.chooseReg(info, val, used)
		if chosen != oldreg {
			// handle this later
			ra.wrongGuesses[val] = true
			val.Reg = chosen
		}
	}

	// a phi must have the same register assigned to itself and all args
	if val.Op == op.Phi {
		for i := 0; i < val.NumArgs(); i++ {
			arg := val.Arg(i)
			if arg.Reg != reg.None {
				return arg.Reg
			}
			liveThroughCalls = liveThroughCalls || ra.liveThroughCalls[arg]
		}
	}

	if val.Op == op.PhiCopy {
		// should have one use, which is the phi
		phi := val.ArgUse(0)
		if phi.Op != op.Phi {
			log.Panicf("expecting %s to be a phi!", phi.String())
		}

		// if the phi already has a reg, go with that
		if phi.Reg != reg.None {
			return phi.Reg
		}

		// otherwise scan the phi's args and run with the first reg assigned
		for i := 0; i < phi.NumArgs(); i++ {
			arg := phi.Arg(i)
			if arg.Reg != reg.None {
				return arg.Reg
			}
			liveThroughCalls = liveThroughCalls || ra.liveThroughCalls[arg]
		}

		liveThroughCalls = liveThroughCalls || ra.liveThroughCalls[phi]
	}

	if len(ra.affinities[val]) > 0 {
		votes := make(map[reg.Reg]int)
		if val.Reg != reg.None && (used&val.Reg) == 0 {
			votes[val.Reg]++
		}
		for _, v := range ra.affinities[val] {
			notInUse := (used&v.Reg) == 0 || (info.regValues[v.Reg] == v && val.Op.IsCopy())
			if val.Func().NumCalls > 0 && v.Reg.IsArgReg() {
				notInUse = false
			}
			if liveThroughCalls && !v.Reg.IsSavedReg() {
				notInUse = false
			}
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
			if liveThroughCalls && !chosen.IsSavedReg() {
				log.Panicf("Choosing a non saved reg! %s", chosen)
			}

			log.Println("affinity chosen", val, chosen, ra.affinities[val], votes)

			return chosen
		} else if len(ra.affinities[val]) > 0 {
			log.Println("affinity failure:", val, ra.affinities[val], used, votes)
		}
	}

	sets := [][]reg.Reg{reg.TempRegs, reg.ArgRegs, reg.RevSavedRegs}

	if liveThroughCalls {
		sets = [][]reg.Reg{reg.SavedRegs, reg.TempRegs}
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
