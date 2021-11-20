package parser

import (
	"fmt"
	"go/constant"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"golang.org/x/tools/go/ssa"
)

func genName(pkg, name string) string {
	sname := strings.Replace(name, "$", "_", -1)
	return fmt.Sprintf("%s__%s", pkg, sname)
}

func walkFunc(pkg *ir.Package, fn *ssa.Function) {
	function := &ir.Func{
		Name: genName(fn.Pkg.Pkg.Name(), fn.Name()),
		Type: fn.Signature,
		Mod:  pkg,
	}

	valmap := make(map[ssa.Value]*ir.Value)
	storemap := make(map[*ssa.Store]*ir.Value)

	blockmap := make(map[*ssa.BasicBlock]*ir.Block)

	for _, param := range fn.Params {
		irp := function.NewValue(op.Parameter, param.Type())
		irp.Value = constant.MakeString(param.Name())
		function.Params = append(function.Params, irp)
		valmap[param] = irp
	}

	// order blocks by reverse succession
	blockList := reverseSSASuccessorSort(fn.Blocks[0], nil, make(map[*ssa.BasicBlock]bool))

	// reverse it to get succession ordering
	for i, j := 0, len(blockList)-1; i < j; i, j = i+1, j-1 {
		blockList[i], blockList[j] = blockList[j], blockList[i]
	}

	type critical struct {
		pred *ssa.BasicBlock
		succ *ssa.BasicBlock
		blk  *ir.Block
	}
	var criticals []critical

	for bn, block := range blockList {
		irBlock := function.NewBlock(ir.Block{
			Comment: block.Comment,
		})
		function.InsertBlock(-1, irBlock)

		if bn == 0 {
			for _, param := range function.Params {
				irBlock.InsertInstr(-1, param)
			}
		}

		blockmap[block] = irBlock

		walkInstrs(irBlock, block.Instrs, valmap, storemap)

		for _, succ := range block.Succs {
			if len(block.Succs) > 1 && len(succ.Preds) > 1 {
				irBlock := function.NewBlock(ir.Block{
					Op:      op.Jump,
					Comment: block.Comment + "." + succ.Comment,
				})
				function.InsertBlock(-1, irBlock)

				criticals = append(criticals, critical{
					pred: block,
					succ: succ,
					blk:  irBlock,
				})
			}
		}
	}

	for _, block := range blockList {
		irBlock := blockmap[block]
		for _, succ := range block.Succs {
			found := false
			for _, crit := range criticals {
				if crit.pred == block && crit.succ == succ {
					irBlock.AddSucc(crit.blk)
					crit.blk.AddPred(irBlock)
					found = true
					break
				}
			}

			if !found {
				irBlock.AddSucc(blockmap[succ])
			}
		}
		for _, pred := range block.Preds {
			found := false
			for _, crit := range criticals {
				if crit.pred == pred && crit.succ == block {
					irBlock.AddPred(crit.blk)
					crit.blk.AddSucc(irBlock)
					found = true
					break
				}
			}

			if !found {
				irBlock.AddPred(blockmap[pred])
			}
		}

		irBlock.Idom = blockmap[block.Idom()]

		for _, dom := range block.Dominees() {
			irBlock.Dominees = append(irBlock.Dominees, blockmap[dom])
		}

		if irBlock.Op != op.Jump {
			irBlock.SetControls(getArgs(irBlock, block.Instrs[len(block.Instrs)-1], valmap))
		}

		var linelist []token.Pos

		// do a pass to resolve args
		for i, instr := range block.Instrs {
			pos := getPos(instr)
			linelist = append(linelist, pos)

			// skip the last op if the block has a op other than jump
			if i == (len(block.Instrs)-1) && irBlock.Op != op.Jump {
				continue
			}

			args := getArgs(irBlock, instr, valmap)
			var irVal *ir.Value
			if len(args) > 0 {
				if val, ok := instr.(ssa.Value); ok {
					irVal = valmap[val]
				} else if val, ok := instr.(*ssa.Store); ok {
					irVal = storemap[val]
				} else if _, ok := instr.(*ssa.DebugRef); ok {
					continue
				} else {
					log.Fatalf("can't look up args for %#v", instr)
				}
				for _, arg := range args {
					irVal.InsertArg(-1, arg)
				}

				// double check everything was wired up correctly
				var foundVal *ir.Value
				for j := 0; j < irBlock.NumInstrs(); j++ {
					val := irBlock.Instr(j)
					if val == irVal {
						foundVal = val
					}
				}
				if foundVal == nil {
					log.Fatalf("val not found! %s", irVal.LongString())
				}
			}
		}

		fset := fn.Prog.Fset
		var lines []string
		filename := ""
		src := ""
		lastline := 0
		for _, pos := range linelist {
			if pos != token.NoPos {
				position := fset.PositionFor(pos, true)
				if filename != position.Filename {
					filename = position.Filename
					buf, err := os.ReadFile(position.Filename)
					if err != nil {
						log.Fatal(err)
					}
					lines = strings.Split(string(buf), "\n")
				}

				if position.Line != lastline {
					lastline = position.Line
					src += strings.TrimSpace(lines[position.Line-1]) + "\n"
				}
			}
		}
		if filename != "" {
			irBlock.Source = src
		}
	}

	pkg.Funcs = append(pkg.Funcs, function)
}

func reverseSSASuccessorSort(block *ssa.BasicBlock, list []*ssa.BasicBlock, visited map[*ssa.BasicBlock]bool) []*ssa.BasicBlock {
	visited[block] = true

	for i := len(block.Succs) - 1; i >= 0; i-- {
		succ := block.Succs[i]
		if !visited[succ] {
			list = reverseSSASuccessorSort(succ, list, visited)
		}
	}

	return append(list, block)
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
				arg = block.Func().Const(con.Type(), con.Value)

			case *ssa.Function:
				name := genName(con.Pkg.Pkg.Name(), con.Name())
				arg = block.Func().Const(con.Type(), constant.MakeString(name))
				arg.Op = op.Func
				block.Func().NumCalls++

			case *ssa.Global:
				name := fmt.Sprintf("\"%s\"", genName(con.Pkg.Pkg.Name(), con.Name()))
				ok = false
				for _, glob := range block.Func().Mod.Globals {
					if glob.Value.String() == name {
						arg = glob
						ok = true
						break
					}
				}
				if !ok {
					name := genName(con.Pkg.Pkg.Name(), con.Name())
					arg = block.Func().Mod.AddGlobal(name, con.Type())
					ok = true
				}
				block.Func().Globals = append(block.Func().Globals, arg)

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
