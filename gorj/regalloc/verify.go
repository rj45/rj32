package regalloc

import (
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
)

func (ra *regAlloc) verify(firstPass bool) {
	ra.Func.Blocks()[0].VisitSuccessors(func(blk *ir.Block) bool {
		info := &ra.blockInfo[blk.ID()]
		info.regValues = make(map[reg.Reg]*ir.Value)

		info.regValues[reg.GP] = ra.Func.FixedReg(reg.GP)
		info.regValues[reg.SP] = ra.Func.FixedReg(reg.SP)

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
				if len(phiIns) > i && phiIns[i][val] {
					// skip incoming phi values otherwise will get a conflict
					continue
				}

				if !firstPass && info.regValues[r] != nil && r != reg.None && info.regValues[r] != val {
					log.Panicf("conflict in reg %s from block %s:%s -> %s, has both %s and %s", r, blk.Func().Name, pred, blk, info.regValues[r].IDString(), val.IDString())
				}

				info.regValues[r] = val
			}
		}

		for i := 0; i < blk.NumInstrs(); i++ {
			val := blk.Instr(i)

			if val.Op != op.Phi {
				for i := 0; i < val.NumArgs(); i++ {
					arg := val.Arg(i)
					if !firstPass && !val.Op.IsConst() && arg.Reg != reg.None && info.regValues[arg.Reg] != arg {
						log.Panicf("in block %s:%s, attempted to read %s from reg %s, but contained %s! %s", blk.Func().Name, blk, arg.IDString(), arg.Reg, info.regValues[arg.Reg], val.LongString())
					}
				}
			}

			for _, val := range info.kills[val] {
				if val.Reg != reg.GP && val.Reg != reg.SP {
					delete(info.regValues, val.Reg)
				}
			}

			info.regValues[val.Reg] = val
		}

		for val := range info.blkKills {
			if val.Reg != reg.GP && val.Reg != reg.SP {
				delete(info.regValues, val.Reg)
			}
		}

		return true
	})
}
