package parser

import (
	"go/constant"
	"go/token"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"golang.org/x/tools/go/ssa"
)

func walkInstrs(block *ir.Block, instrs []ssa.Instruction, valmap map[ssa.Value]*ir.Value, storemap map[*ssa.Store]*ir.Value) {
	for _, instr := range instrs {
		irInstr := ir.Value{}
		var store *ssa.Store

		// ops = instr.Operands(ops[:0])
		switch ins := instr.(type) {
		case *ssa.DebugRef:
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
			store = ins
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
				log.Fatalf("unsupported binop: %#v", ins)
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
				log.Fatalf("unsupported unop: %#v", ins)
			}

		case *ssa.RunDefers:
			// ignore
		default:
			log.Fatalf("unknown instruction type %#v", instr)
		}

		if irInstr.Op != op.Invalid {
			ins := block.Func().NewValue(irInstr.Op, irInstr.Type)
			ins.Value = irInstr.Value

			if ins.Type == nil {
				if typed, ok := instr.(interface{ Type() types.Type }); ok {
					ins.Type = typed.Type()
				}
			}

			block.InsertInstr(-1, ins)
			if store != nil {
				storemap[store] = ins
			}

			if vin, ok := instr.(ssa.Value); ok {
				valmap[vin] = ins
			}
		}
	}
}
