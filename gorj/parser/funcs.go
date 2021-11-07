package parser

import (
	"fmt"
	"go/constant"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"golang.org/x/tools/go/ssa"
)

func walkFunc(mod *ir.Module, fn *ssa.Function) {
	function := &ir.Func{
		Name: fmt.Sprintf("%s.%s", fn.Pkg.Pkg.Name(), fn.Name()),
		Type: fn.Signature,
		Mod:  mod,
	}

	valmap := make(map[ssa.Value]*ir.Value)
	storemap := make(map[*ssa.Store]*ir.Value)

	blockmap := make(map[*ssa.BasicBlock]*ir.Block)

	for _, param := range fn.Params {
		irp := &ir.Value{
			ID:    function.NextInstrID(),
			Op:    op.Parameter,
			Type:  param.Type(),
			Value: constant.MakeString(param.Name()),
		}
		function.Params = append(function.Params, irp)
		valmap[param] = irp
	}

	for _, block := range fn.Blocks {
		irBlock := &ir.Block{
			ID:      function.NextBlockID(),
			Comment: block.Comment,
			Func:    function,
		}
		function.Blocks = append(function.Blocks, irBlock)

		blockmap[block] = irBlock

		walkInstrs(irBlock, block.Instrs, valmap, storemap)
	}

	for _, block := range fn.Blocks {
		irBlock := blockmap[block]
		for i, succ := range block.Succs {
			irBlock.Succs = append(irBlock.Succs, ir.BlockRef{Index: i, Block: blockmap[succ]})
		}
		for i, pred := range block.Preds {
			irBlock.Preds = append(irBlock.Preds, ir.BlockRef{Index: i, Block: blockmap[pred]})
		}

		irBlock.Idom = blockmap[block.Idom()]

		for _, dom := range block.Dominees() {
			irBlock.Dominees = append(irBlock.Dominees, blockmap[dom])
		}

		if irBlock.Op == op.If || irBlock.Op == op.Return {
			irBlock.Controls = getArgs(irBlock, block.Instrs[len(block.Instrs)-1], valmap)
		}

		// do a pass to resolve args
		for i, instr := range block.Instrs {
			if i == (len(block.Instrs)-1) && (irBlock.Op == op.If || irBlock.Op == op.Return) {
				continue
			}
			args := getArgs(irBlock, instr, valmap)
			var irVal *ir.Value
			if len(args) > 0 {
				if val, ok := instr.(ssa.Value); ok {
					irVal = valmap[val]
				} else if val, ok := instr.(*ssa.Store); ok {
					irVal = storemap[val]
				} else {
					log.Fatalf("can't look up args for %#v", instr)
				}
				irVal.Args = args

				// double check everything was wired up correctly
				var foundVal *ir.Value
				for _, val := range irBlock.Instrs {
					if val == irVal {
						foundVal = val
					}
				}
				if foundVal == nil {
					log.Fatalf("val not found! %s", irVal.LongString())
				}
			}
		}
	}

	mod.Funcs = append(mod.Funcs, function)
}

func getArgs(block *ir.Block, instr ssa.Instruction, valmap map[ssa.Value]*ir.Value) []*ir.Value {
	var args []*ir.Value

	var valarr [5]*ssa.Value
	vals := instr.Operands(valarr[:])

	for _, val := range vals {
		if val == nil {
			continue
		}
		arg, ok := valmap[*val]
		if !ok {
			ok = true
			switch con := (*val).(type) {
			case *ssa.Const:
				arg = &ir.Value{
					ID:    block.NextInstrID(),
					Op:    op.Const,
					Type:  con.Type(),
					Value: con.Value,
				}
				block.Func.Consts = append(block.Func.Consts, arg)

			case *ssa.Function:
				arg = &ir.Value{
					ID:    block.NextInstrID(),
					Op:    op.Func,
					Type:  con.Type(),
					Value: constant.MakeString(con.Name()),
				}
				block.Func.Calls = append(block.Func.Calls, arg)

			case *ssa.Global:
				name := fmt.Sprintf("\"%s.%s\"", con.Pkg.Pkg.Name(), con.Name())
				ok = false
				for _, glob := range block.Func.Mod.Globals {
					if glob.Value.String() == name {
						arg = glob
						ok = true
						break
					}
				}

			default:
				ok = false
			}
			if ok {
				valmap[*val] = arg
			}
		}
		if ok && arg != nil {
			args = append(args, arg)
		} else {
			log.Printf("Unmapped value: %#v\n", *val)
		}
	}

	return args
}
