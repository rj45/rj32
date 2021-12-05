package a32

import (
	"fmt"
	"go/constant"
	"go/types"
	"log"

	"github.com/rj45/rj32/gorj/arch"
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/ir/reg"
	"github.com/rj45/rj32/gorj/sizes"
	"github.com/rj45/rj32/gorj/xform"
)

type cpuArch struct{}

var _ = arch.Register(cpuArch{})

func (cpuArch) Name() string {
	return "a32"
}

func (cpuArch) XformTags() []xform.Tag {
	return []xform.Tag{xform.HasFramePointer}
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
		asmGlob.Section = asm.Code
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

			strs := []string{
				"#d le(($+8)`32)",
				fmt.Sprintf("#d le(%d`32)", len(str)),
				fmt.Sprintf("#d %q", str),
				"#align 32",
			}

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
	if opcode == NOP {
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
				Op:   MOV,
				Args: []*asm.Var{varFor(val), {String: "0"}},
			}, &asm.Instr{
				Op:   opcode,
				Args: []*asm.Var{varFor(val.Arg(0)), varFor(val.Arg(1))},
			}, &asm.Instr{
				Op:     MOV,
				Args:   []*asm.Var{varFor(val), {String: "1"}},
				Indent: true,
			})
		}

		switch val.Op {
		case op.Load:
			switch sizes.Sizeof(val.Type) {
			case 1:
				opcode = LD8
			case 2:
				opcode = LD16
			case 4:
				opcode = LD
			default:
				log.Fatalln("load of unsupported size", sizes.Sizeof(val.Type), "type:", val.Type, val.ShortString(), "in func", val.Func().Name)
			}
		case op.Store:
			switch sizes.Sizeof(val.Type) {
			case 1:
				opcode = ST8
			case 2:
				opcode = ST16
			case 4:
				opcode = ST
			default:
				log.Fatalln("store of unsupported size", sizes.Sizeof(val.Type), "type:", val.Type, val.ShortString(), "in func", val.Func().Name)
			}
		case op.Extract:
			// ignore
		case op.SwapOut:
			// ignore
		case op.Phi:
			// ignore
		case op.PhiCopy, op.Copy:
			if val.Arg(0).Reg != reg.None {
				opcode = MOV
			} else {
				opcode = LDI
			}
		default:
			log.Panicf("unable to assemble %s", val.ShortString())
		}
	}

	if opcode == NOP {
		return list
	}

	list = append(list, &asm.Instr{
		Op:   opcode,
		Args: opcode.Fmt().Vars(val),
	})

	return list
}

var signedCompareOps = map[op.Op]Opcode{
	op.Equal:        BR_EQ,
	op.NotEqual:     BR_NEQ,
	op.Less:         BR_S_L,
	op.LessEqual:    BR_S_LE,
	op.Greater:      BR_S_G,
	op.GreaterEqual: BR_S_GE,
}

var unsignedCompareOps = map[op.Op]Opcode{
	op.Equal:        BR_EQ,
	op.NotEqual:     BR_NEQ,
	op.Less:         BR_U_L,
	op.LessEqual:    BR_U_LE,
	op.Greater:      BR_U_G,
	op.GreaterEqual: BR_U_GE,
}

func (cpuArch) AssembleBlockOp(list []*asm.Instr, blk *ir.Block, flip bool) []*asm.Instr {
	switch blk.Op {
	case op.Jump:
		list = append(list, &asm.Instr{
			Op:   BRA,
			Args: []*asm.Var{blockVar(blk.Succ(0))},
		})

	case op.Return:
		paramSlots := len(blk.Func().Params) - len(reg.ArgRegs)
		if paramSlots < 0 {
			paramSlots = 0
		}

		list = append(list, &asm.Instr{
			Op:   RET,
			Args: []*asm.Var{{String: fmt.Sprintf("%d", paramSlots)}},
		})

	case op.Panic:
		list = append(list, &asm.Instr{
			Op:   MOV,
			Args: []*asm.Var{{String: "a0"}, {Value: blk.Control(0)}},
		}, &asm.Instr{
			Op: BRK,
		}, &asm.Instr{
			Op: ERR,
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
		Op:   CMP,
		Args: controls,
	}, &asm.Instr{
		Op:   opcode,
		Args: []*asm.Var{succ[0]},
	}, &asm.Instr{
		Op:   BRA,
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
