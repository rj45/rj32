package rj32

import (
	"fmt"
	"go/constant"
	"go/types"
	"log"
	"unicode/utf16"

	"github.com/rj45/rj32/gorj/arch"
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/sizes"
)

type cpuArch struct{}

var _ = arch.Register(cpuArch{})

func (cpuArch) Name() string {
	return "rj32"
}

func (cpuArch) AssembleGlobal(glob *ir.Value) *asm.Global {
	asmGlob := &asm.Global{
		Value: glob,
	}

	typ := glob.Type.Underlying()
	if ptr, ok := typ.(*types.Pointer); ok {
		typ = ptr.Elem()
	}

	if glob.NumArgs() > 0 {
		asmGlob.Section = asm.Data
	} else {
		asmGlob.Section = asm.Bss
	}

	name := constant.StringVal(glob.Value)

	asmGlob.Label = name
	asmGlob.Comment = typ.String()

	if glob.NumArgs() > 0 {
		data := glob.Arg(0).Value
		if data.Kind() == constant.String {
			str := constant.StringVal(data)

			runes := []rune(str)
			utf16 := utf16.Encode(runes)

			strs := []string{
				"#d16 $+2",
				fmt.Sprintf("#d16 %d", len(utf16)),
				fmt.Sprintf("; %q", str),
			}

			hex := "#d16 "
			for i, v := range utf16 {
				if i != 0 && i%8 != 0 {
					hex += ", "
				} else if i != 0 {
					strs = append(strs, hex)
					hex = "#d16 "
				}
				hex += fmt.Sprintf("0x%04x", v)
			}
			strs = append(strs, hex)
			asmGlob.Strings = strs
		}
	} else {
		size := sizes.Sizeof(typ)
		asmGlob.Strings = []string{fmt.Sprintf("#res %d", size)}
	}

	return asmGlob
}

func (cpuArch) AssembleInstr(list []*asm.Instr, val *ir.Value) []*asm.Instr {
	opcode := translations[val.Op]
	if opcode == Nop {
		// special cases

		// storing a comparison into a bool
		if val.Op.IsCompare() {
			opcode = signedCompareOps[val.Op]
			if isUnsigned(val.Arg(0).Type) {
				opcode = unsignedCompareOps[val.Op]
			}

			// move ${val}, 0
			// if_cc ${arg0}, ${arg1}
			//   move ${val}, 1
			return append(list, &asm.Instr{
				Op:   Move,
				Args: []*asm.Var{varFor(val), {String: "0"}},
			}, &asm.Instr{
				Op:   opcode,
				Args: []*asm.Var{varFor(val.Arg(0)), varFor(val.Arg(1))},
			}, &asm.Instr{
				Op:     Move,
				Args:   []*asm.Var{varFor(val), {String: "1"}},
				Indent: true,
			})
		}

		switch val.Op {
		case op.Extract:
			// ignore
		case op.SwapOut:
			// ignore
		case op.Phi:
			// ignore
		case op.PhiCopy:
			opcode = Move
		default:
			log.Panicf("unable to assemble %s", val.ShortString())
		}
	}

	if opcode == Nop {
		return list
	}

	list = append(list, &asm.Instr{
		Op:   opcode,
		Args: opcode.Fmt().Vars(val),
	})

	return list
}

var signedCompareOps = map[op.Op]Opcode{
	op.Equal:        IfEq,
	op.NotEqual:     IfNe,
	op.Less:         IfLt,
	op.LessEqual:    IfLe,
	op.Greater:      IfGt,
	op.GreaterEqual: IfGe,
}

var unsignedCompareOps = map[op.Op]Opcode{
	op.Equal:        IfEq,
	op.NotEqual:     IfNe,
	op.Less:         IfUlt,
	op.LessEqual:    IfUle,
	op.Greater:      IfUgt,
	op.GreaterEqual: IfUge,
}

func (cpuArch) AssembleBlockOp(list []*asm.Instr, blk *ir.Block, flip bool) []*asm.Instr {
	switch blk.Op {
	case op.Jump:
		list = append(list, &asm.Instr{
			Op:   Jump,
			Args: []*asm.Var{blockVar(blk.Succ(0))},
		})

	case op.Return:
		list = append(list, &asm.Instr{
			Op: Return,
		})

	case op.Panic:
		list = append(list, &asm.Instr{
			Op:   Move,
			Args: []*asm.Var{{String: "a0"}, {Value: blk.Control(0)}},
		}, &asm.Instr{
			Op: Error,
		})

	case op.If:
		list = asmIf(list, op.NotEqual,
			[]*asm.Var{{Value: blk.Control(0)}, {String: "0"}},
			[]*asm.Var{blockVar(blk.Succ(0)), blockVar(blk.Succ(1))}, flip)

	case op.IfEqual, op.IfNotEqual, op.IfLess, op.IfLessEqual, op.IfGreater, op.IfGreaterEqual:
		list = asmIf(list, blk.Op.Compare(),
			[]*asm.Var{varFor(blk.Control(0)), varFor(blk.Control(1))},
			[]*asm.Var{blockVar(blk.Succ(0)), blockVar(blk.Succ(1))}, flip)

	default:
		log.Panicln("unimplemented block op:", blk.Op)
	}

	return list
}

func asmIf(list []*asm.Instr, op op.Op, controls []*asm.Var, succ []*asm.Var, flip bool) []*asm.Instr {
	var opcode Opcode

	if flip {
		op = op.Opposite()
		succ[0], succ[1] = succ[1], succ[0]
	}

	if isUnsigned(controls[0].Value.Type) {
		opcode = unsignedCompareOps[op]
	} else {
		opcode = signedCompareOps[op]
	}

	return append(list, &asm.Instr{
		Op:   opcode,
		Args: controls,
	}, &asm.Instr{
		Op:     Jump,
		Args:   []*asm.Var{succ[0]},
		Indent: true,
	}, &asm.Instr{
		Op:   Jump,
		Args: []*asm.Var{succ[1]},
	})
}

func blockVar(blk *ir.Block) *asm.Var {
	return &asm.Var{String: "." + blk.String(), Block: blk}
}

func isUnsigned(typ types.Type) bool {
	basic, ok := typ.(*types.Basic)
	if !ok {
		return false
	}
	switch basic.Kind() {
	case types.Uint, types.Uint8, types.Uint16, types.Uint32, types.Uint64, types.Uintptr:
		return true
	}
	return false
}
