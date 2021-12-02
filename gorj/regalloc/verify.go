package regalloc

import (
	"flag"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

var debugVerify = flag.Bool("debugverifier", false, "emit register allocation verification logs")

func (ra *RegAlloc) verify(firstPass bool) {
	entry := true

	ra.Func.Blocks()[0].VisitSuccessors(func(blk *ir.Block) bool {
		info := &ra.blockInfo[blk.ID()]
		info.regValues = make(map[reg.Reg]*ir.Value)

		info.regValues[reg.GP] = ra.Func.FixedReg(reg.GP)
		info.regValues[reg.SP] = ra.Func.FixedReg(reg.SP)

		if entry {
			entry = false

			for i := range ra.Func.Params {
				if i > len(reg.ArgRegs) {
					break
				}
				info.regValues[reg.ArgRegs[i]] = ra.Func.FixedReg(reg.ArgRegs[i])
			}
		}

		var phiIns []map[*ir.Value]bool

		for i := 0; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)
			if val.Op != op.Phi {
				break
			}

			for i := 0; i < val.NumArgs(); i++ {
				arg := val.Arg(i)

				if i >= len(phiIns) {
					phiIns = append(phiIns, make(map[*ir.Value]bool))
				}

				phiIns[i][arg] = true
			}
		}

		for i := 0; i < blk.NumPreds(); i++ {
			pred := blk.Pred(i)
			pinfo := &ra.blockInfo[pred.ID()]

			for r, val := range pinfo.regValues {
				if !val.NeedsReg() {
					continue
				}

				if len(phiIns) > i && phiIns[i][val] {
					// skip incoming phi values otherwise will get a conflict
					continue
				}

				// if !firstPass && info.regValues[r] != nil && r != reg.None && info.regValues[r] != val {
				// 	log.Panicf("conflict in reg %s from block %s:%s -> %s, has both %s and %s", r, blk.Func().Name, pred, blk, info.regValues[r].IDString(), val.IDString())
				// }

				if !firstPass && val.Reg == reg.None {
					log.Panicf("in block %s:%s, attempted to define an unassigned register! %s", blk.Func().Name, blk, val.ShortString())
				}

				if *debugVerify && !firstPass {
					log.Println("Putting", val.IDString(), "into", val.Reg)
				}
				info.regValues[r] = val
			}
		}

		firstPhiCopy := blk.NumInstrs()
		for i := 0; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			if val.Op == op.PhiCopy {
				firstPhiCopy = i
				break
			}

			if val.Op != op.Phi {
				for i := 0; i < val.NumArgs(); i++ {
					arg := val.Arg(i)
					if !firstPass && arg.Reg == reg.None && arg.NeedsReg() {
						log.Panicf("in block %s:%s, attempted to read %s from an unassigned register! %s", blk.Func().Name, blk, arg, val.ShortString())
					}
					if !firstPass && !val.Op.IsConst() && arg.Reg != reg.None && info.regValues[arg.Reg] != arg {
						log.Panicf("in block %s:%s, attempted to read %s from reg %s, but contained %s! %s", blk.Func().Name, blk, arg.IDString(), arg.Reg, info.regValues[arg.Reg].IDString(), val.ShortString())
					}
				}
			}

			for _, val := range info.kills[val] {
				if val.Reg != reg.GP && val.Reg != reg.SP {
					if *debugVerify && !firstPass {
						log.Println("Killing", val.IDString(), "from", val.Reg)
					}
					delete(info.regValues, val.Reg)
				}
			}

			if val.NeedsReg() {
				if *debugVerify && !firstPass {
					log.Println("Putting", val.IDString(), "into", val.Reg)
				}
				if !firstPass && val.Reg == reg.None {
					log.Panicf("in block %s:%s, attempted to define an unassigned register! %s", blk.Func().Name, blk, val.ShortString())
				}
				info.regValues[val.Reg] = val
			}
		}

		// process phi copies which are parallel copies
		// first process all reads
		for i := firstPhiCopy; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			for i := 0; i < val.NumArgs(); i++ {
				arg := val.Arg(i)
				if !firstPass && arg.Reg == reg.None && arg.NeedsReg() {
					log.Panicf("in block %s:%s, attempted to read %s from an unassigned register! %s", blk.Func().Name, blk, arg, val.ShortString())
				}
				if !firstPass && !val.Op.IsConst() && arg.Reg != reg.None && info.regValues[arg.Reg] != arg {
					log.Panicf("in block %s:%s, attempted to read %s from reg %s, but contained %s! %s", blk.Func().Name, blk, arg.IDString(), arg.Reg, info.regValues[arg.Reg].IDString(), val.ShortString())
				}
			}
		}

		// then kill the uses
		for i := firstPhiCopy; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)
			for _, val := range info.kills[val] {
				if val.Reg != reg.GP && val.Reg != reg.SP {
					if *debugVerify && !firstPass {
						log.Println("Killing", val.IDString(), "from", val.Reg)
					}
					delete(info.regValues, val.Reg)
				}
			}
		}

		// then do all writes
		for i := firstPhiCopy; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			if val.NeedsReg() {
				if *debugVerify && !firstPass {
					log.Println("Putting", val.IDString(), "into", val.Reg)
				}
				if !firstPass && val.Reg == reg.None {
					log.Panicf("in block %s:%s, attempted to define an unassigned register! %s", blk.Func().Name, blk, val.ShortString())
				}
				info.regValues[val.Reg] = val
			}
		}

		for i := 0; i < blk.NumControls(); i++ {
			val := blk.Control(i)
			if !firstPass && !val.Op.IsConst() && val.Reg != reg.None && info.regValues[val.Reg] != val {
				log.Panicf("in block %s:%s, attempted to read control %s from reg %s, but contained %s! %s", blk.Func().Name, blk, val.IDString(), val.Reg, info.regValues[val.Reg].IDString(), val.ShortString())
			}
		}

		for val := range info.blkKills {
			if val.Reg != reg.GP && val.Reg != reg.SP {
				if *debugVerify && !firstPass {
					log.Println("Killing", val.IDString(), "from", val.Reg)
				}
				delete(info.regValues, val.Reg)
			}
		}

		return true
	})
}
