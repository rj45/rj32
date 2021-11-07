package main

import (
	"bytes"
	"fmt"
	"go/build"
	"go/constant"
	"go/token"
	"go/types"
	"io"
	"log"
	"os"
	"sort"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type members []ssa.Member

func (m members) Len() int           { return len(m) }
func (m members) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m members) Less(i, j int) bool { return m[i].Pos() < m[j].Pos() }

// toSSA converts go source to SSA
func toSSA(source io.Reader, fileName, packageName string, debug bool) ([]byte, members, error) {
	// adopted from saa package example

	conf := loader.Config{
		Build: &build.Default,
	}

	file, err := conf.ParseFile(fileName, source)
	if err != nil {
		return nil, nil, err
	}

	conf.CreateFromFiles("main.go", file)

	prog, err := conf.Load()
	if err != nil {
		return nil, nil, err
	}

	ssaProg := ssautil.CreateProgram(prog, ssa.BuildSerially)
	ssaProg.Build()
	mainPkg := ssaProg.Package(prog.InitialPackages()[0].Pkg)

	out := new(bytes.Buffer)
	mainPkg.SetDebugMode(debug)
	mainPkg.WriteTo(out)
	mainPkg.Build()

	// grab just the functions
	funcs := members([]ssa.Member{})
	all := members([]ssa.Member{})
	for _, obj := range mainPkg.Members {
		if obj.Token() == token.FUNC {
			funcs = append(funcs, obj)
		}
		all = append(all, obj)
	}
	// sort by Pos()
	sort.Sort(funcs)
	sort.Sort(all)
	for _, f := range funcs {
		mainPkg.Func(f.Name()).WriteTo(out)
	}
	return out.Bytes(), all, nil
}

func walk(mod *ir.Module, all members) {
	for _, member := range all {
		switch member.Token() {
		case token.VAR:
			// walkGlobal(member.Package().Var(member.Name()))
			mod.Globals = append(mod.Globals, &ir.Value{
				Op:    op.Global,
				Value: constant.MakeString(member.Name()),
				Type:  member.Type(),
			})
		// case Token.CONST:
		// case Token.TYPE:
		default:
		}
	}

	for _, member := range all {
		switch member.Token() {
		case token.FUNC:
			walkFunc(mod, member.Package().Func(member.Name()))
		case token.VAR:
		default:
			log.Fatalln("unknown type", member.Token())
		}
	}
}

func walkFunc(mod *ir.Module, fn *ssa.Function) {
	function := &ir.Func{
		Name: fn.Name(),
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
		// TODO: this is bugged: it does not look up all the args, not sure why
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

func walkInstrs(block *ir.Block, instrs []ssa.Instruction, valmap map[ssa.Value]*ir.Value, storemap map[*ssa.Store]*ir.Value) {
	for _, instr := range instrs {
		irInstr := &ir.Value{
			Block: block,
		}

		// ops = instr.Operands(ops[:0])
		switch ins := instr.(type) {
		case *ssa.If:
			block.Op = op.If
		case *ssa.Jump:
			block.Op = op.Jump
		case *ssa.Return:
			block.Op = op.Return
		case *ssa.Phi:
			irInstr.Op = op.Phi
		case *ssa.Store:
			irInstr.Type = ins.Val.Type()
			irInstr.Op = op.Store
			storemap[ins] = irInstr
		case *ssa.Alloc:
			irInstr.Value = constant.MakeString(ins.Comment)
			if ins.Heap {
				irInstr.Op = op.New
			} else {
				irInstr.Op = op.Local
			}
		case *ssa.Call:
			irInstr.Op = op.Call
			switch call := ins.Call.Value.(type) {
			case *ssa.Function:
				// irInstr.Value = constant.MakeString(call.Name())
				irInstr.Type = call.Signature
			case *ssa.Builtin:
				irInstr.Value = constant.MakeString(call.Name())
				irInstr.Type = call.Type()
			default:
				log.Fatalf("unsupported call type: %#v", ins.Call.Value)
			}

		case *ssa.Convert:
			irInstr.Op = op.Convert
		case *ssa.IndexAddr:
			irInstr.Op = op.IndexAddr
		case *ssa.FieldAddr:
			irInstr.Op = op.FieldAddr
			irInstr.Value = constant.MakeInt64(int64(ins.Field))
		case *ssa.BinOp:
			switch ins.Op {
			case token.ADD:
				irInstr.Op = op.Add
			case token.SUB:
				irInstr.Op = op.Sub
			case token.MUL:
				irInstr.Op = op.Mul
			case token.QUO:
				irInstr.Op = op.Div
			case token.REM:
				irInstr.Op = op.Rem
			case token.AND:
				irInstr.Op = op.And
			case token.OR:
				irInstr.Op = op.Or
			case token.XOR:
				irInstr.Op = op.Xor
			case token.SHL:
				irInstr.Op = op.ShiftLeft
			case token.SHR:
				irInstr.Op = op.ShiftRight
			case token.AND_NOT:
				irInstr.Op = op.AndNot
			case token.EQL:
				irInstr.Op = op.Equal
			case token.NEQ:
				irInstr.Op = op.NotEqual
			case token.LSS:
				irInstr.Op = op.Less
			case token.LEQ:
				irInstr.Op = op.LessEqual
			case token.GTR:
				irInstr.Op = op.Greater
			case token.GEQ:
				irInstr.Op = op.GreaterEqual
			default:
				panic("not handled")
			}
		case *ssa.UnOp:
			switch ins.Op {
			case token.NOT:
				irInstr.Op = op.Not
			case token.SUB:
				irInstr.Op = op.Negate
			case token.MUL:
				irInstr.Op = op.Load
			case token.XOR:
				irInstr.Op = op.Invert
			default:
				panic("not handled")
			}

		case *ssa.RunDefers:
			// ignore
		default:
			log.Fatalf("unknown type %#v", instr)
		}

		if irInstr.Op != op.Invalid {
			irInstr.ID = block.NextInstrID()

			if irInstr.Type == nil {
				if typed, ok := instr.(interface{ Type() types.Type }); ok {
					irInstr.Type = typed.Type()
				}
			}

			block.Instrs = append(block.Instrs, irInstr)

			if vin, ok := instr.(ssa.Value); ok {
				valmap[vin] = irInstr
			}
		}
	}
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
				name := fmt.Sprintf("\"%s\"", con.Name())
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

func main() {
	log.SetFlags(log.Lshortfile)

	file, err := os.Open("./testfiles/seive/seive.go")
	if err != nil {
		log.Fatal(err)
	}

	dump, members, err := toSSA(file, "seive.go", "main", true)
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(dump)

	fmt.Println("\n-------------------")

	var mod ir.Module

	walk(&mod, members)

	// pretty.Println(&mod)
	fmt.Println(mod.LongString())
}
