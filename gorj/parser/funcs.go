package parser

import (
	"bytes"
	"fmt"
	"go/constant"
	"go/token"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"golang.org/x/tools/go/ssa"
)

func genName(pkg, name string) string {
	sname := strings.Replace(name, "$", "_", -1)
	return fmt.Sprintf("%s__%s", pkg, sname)
}

func walkFunc(function *ir.Func, fn *ssa.Function) {
	valmap := make(map[ssa.Value]*ir.Value)
	storemap := make(map[*ssa.Store]*ir.Value)

	blockmap := make(map[*ssa.BasicBlock]*ir.Block)

	if fn.Blocks == nil {
		// extern function
		handleExternFunc(function, fn)
		return
	}

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
	var returns []*ir.Block

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

		if irBlock.Op == op.Return {
			returns = append(returns, irBlock)
		}

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

	if len(returns) > 1 {
		realRet := function.NewBlock(ir.Block{
			Op:      op.Return,
			Comment: "real.return",
		})
		function.InsertBlock(-1, realRet)

		var phis []*ir.Value

		for i := 0; i < returns[0].NumControls(); i++ {
			phi := function.NewValue(op.Phi, returns[0].Control(i).Type)
			realRet.InsertInstr(-1, phi)
			realRet.InsertControl(-1, phi)
			phis = append(phis, phi)
		}

		for _, ret := range returns {
			ret.AddSucc(realRet)
			realRet.AddPred(ret)
			ret.Op = op.Jump

			for i := 0; i < ret.NumControls(); i++ {
				phis[i].InsertArg(-1, ret.Control(i))
			}

			for ret.NumControls() > 0 {
				ret.RemoveControl(0)
			}
		}
	}
}

func handleExternFunc(function *ir.Func, fn *ssa.Function) {
	filename := fn.Prog.Fset.File(fn.Pos()).Name()
	folder, err := filepath.EvalSymlinks(filepath.Dir(filename))
	if err != nil {
		log.Fatalf("could not follow symlinks for folder %s", folder)
	}
	asm := ""
	filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		ext := filepath.Ext(d.Name())
		if ext == ".asm" || ext == ".s" || ext == ".S" {
			noext := strings.TrimSuffix(d.Name(), ext)
			parts := strings.Split(noext, "_")

			if len(parts) > 1 && parts[len(parts)-1] != arch.Name() {
				// skip files with an underscore and the last part of the name
				// does not match the arch.Name()
				return nil
			}

			buf, err := os.ReadFile(path)
			if err != nil {
				log.Fatalln(err)
			}

			// TODO: find build tags and ensure they match

			lines := bytes.Split(buf, []byte("\n"))
			startLine := -1
			label := []byte(fmt.Sprintf("%s:", fn.Name()))
			for i, line := range lines {
				if bytes.HasPrefix(bytes.TrimSpace(line), label) {
					startLine = i + 1
					break
				}
			}

			if startLine == -1 {
				return nil
			}

			endLine := -1
			for i := startLine; i < len(lines); i++ {
				trimmed := bytes.TrimSpace(lines[i])
				lines[i] = trimmed
				// if doesn't start with a dot, but does end in a colon
				if !bytes.HasPrefix(trimmed, []byte(".")) && bytes.HasSuffix(trimmed, []byte(":")) {
					endLine = i + 1
					break
				}
			}

			if endLine == -1 {
				endLine = len(lines)
			}

			if asm != "" {
				log.Fatalf("found duplicate of extern func %s in %s", fn.Name(), path)
			}

			asm = string(bytes.Join(lines[startLine:endLine], []byte("\n")))
		}
		return nil
	})
	if asm == "" {
		log.Fatalf("could not find assembly for extern func %s path %s", fn.Name(), folder)
	}
	genInlineAsmFunc(function, asm)
}

func genInlineAsmFunc(fn *ir.Func, asm string) {
	entry := fn.NewBlock(ir.Block{
		Comment: "entry",
		Op:      op.Jump,
	})
	body := fn.NewBlock(ir.Block{
		Comment: "inline.asm",
		Op:      op.Jump,
	})
	exit := fn.NewBlock(ir.Block{
		Comment: "exit",
		Op:      op.Return,
	})

	entry.AddSucc(body)
	body.AddSucc(exit)

	exit.AddPred(body)
	body.AddPred(entry)

	fn.InsertBlock(-1, entry)
	fn.InsertBlock(-1, body)
	fn.InsertBlock(-1, exit)
	blk := fn.Blocks()[1]

	val := fn.NewValue(op.InlineAsm, fn.Type.Results())
	val.Value = constant.MakeString(asm)
	blk.InsertInstr(-1, val)
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
				if con.Type().Underlying().String() == "string" {
					str := constant.StringVal(con.Value)
					pkg := block.Func().Pkg
					val, found := pkg.Strings[str]
					if !found {
						// generate a name
						name := fmt.Sprintf("%s_str%d", block.Func().Name, pkg.NextStringNum)
						pkg.NextStringNum++

						// create a global
						arg = pkg.AddGlobal(name, con.Type())

						// attach string to global as its value
						arg.InsertArg(-1, block.Func().Const(con.Type(), con.Value))

						// reuse string
						if pkg.Strings == nil {
							pkg.Strings = make(map[string]*ir.Value)
						}
						pkg.Strings[str] = arg
					} else {
						arg = val
					}

					block.Func().Globals = append(block.Func().Globals, arg)
				} else {

					arg = block.Func().Const(con.Type(), con.Value)
				}

			case *ssa.Function:
				name := genName(con.Pkg.Pkg.Name(), con.Name())
				otherFunc := block.Func().Pkg.LookupFunc(name)
				if otherFunc == nil {
					log.Fatalf("reference to unknown function %s in function %s", name, block.Func().Name)
				}
				// ensure it gets loaded
				otherFunc.Referenced = true
				arg = block.Func().Const(con.Type(), constant.MakeString(name))
				arg.Op = op.Func
				block.Func().NumCalls++

			case *ssa.Builtin:
				name := genName("builtin", con.Name())
				arg = block.Func().Const(con.Type(), constant.MakeString(name))
				arg.Op = op.Func
				block.Func().NumCalls++

			case *ssa.Global:
				name := fmt.Sprintf("\"%s\"", genName(con.Pkg.Pkg.Name(), con.Name()))
				ok = false
				for _, glob := range block.Func().Pkg.Globals {
					if glob.Value.String() == name {
						arg = glob
						ok = true
						break
					}
				}
				if !ok {
					name := genName(con.Pkg.Pkg.Name(), con.Name())
					arg = block.Func().Pkg.AddGlobal(name, con.Type())
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
