package rj32

import (
	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir/op"
)

type Opcode int

//go:generate enumer -type=Opcode -transform snake

const (
	// Natively implemented instructions
	Nop Opcode = iota
	Rets
	Error
	Halt
	Rcsr
	Wcsr
	Move
	Loadc
	Jump
	Imm
	Call
	Imm2
	Load
	Store
	Loadb
	Storeb
	Add
	Sub
	Addc
	Subc
	Xor
	And
	Or
	Shl
	Shr
	Asr
	IfEq
	IfNe
	IfLt
	IfGe
	IfUlt
	IfUge

	// Psuedoinstructions
	Not
	Neg
	Swap
	IfGt
	IfLe
	IfUgt
	IfUle
	Return
)

func (op Opcode) Fmt() asm.Fmt {
	return opDefs[op].fmt
}

func (op Opcode) IsMove() bool {
	return op == Move || op == Swap
}

func (op Opcode) IsCall() bool {
	return op == Call
}

type def struct {
	code Opcode
	fmt  Fmt
	op   op.Op
}

var opDefs = []def{
	{code: Nop, fmt: NoFmt},
	{code: Rets, fmt: NoFmt},
	{code: Error, fmt: NoFmt},
	{code: Halt, fmt: NoFmt},
	{code: Rcsr, fmt: BinaryFmt},
	{code: Wcsr, fmt: BinaryFmt},
	{code: Move, fmt: MoveFmt, op: op.Copy},
	{code: Loadc, fmt: BinaryFmt},
	{code: Jump, fmt: CallFmt},
	{code: Imm, fmt: UnaryFmt},
	{code: Call, fmt: CallFmt, op: op.Call},
	{code: Imm2, fmt: BinaryFmt},
	{code: Load, fmt: LoadFmt, op: op.Load},
	{code: Store, fmt: StoreFmt, op: op.Store},
	{code: Loadb, fmt: LoadFmt},
	{code: Storeb, fmt: StoreFmt},
	{code: Add, fmt: BinaryFmt, op: op.Add},
	{code: Sub, fmt: BinaryFmt, op: op.Sub},
	{code: Addc, fmt: BinaryFmt},
	{code: Subc, fmt: BinaryFmt},
	{code: Xor, fmt: BinaryFmt, op: op.Xor},
	{code: And, fmt: BinaryFmt, op: op.And},
	{code: Or, fmt: BinaryFmt, op: op.Or},
	{code: Shl, fmt: BinaryFmt, op: op.ShiftLeft},
	{code: Shr, fmt: BinaryFmt, op: op.ShiftRight},
	{code: Asr, fmt: BinaryFmt},
	{code: IfEq, fmt: CompareFmt},
	{code: IfNe, fmt: CompareFmt},
	{code: IfLt, fmt: CompareFmt},
	{code: IfGe, fmt: CompareFmt},
	{code: IfUlt, fmt: CompareFmt},
	{code: IfUge, fmt: CompareFmt},
	{code: Not, fmt: UnaryFmt, op: op.Invert},
	{code: Neg, fmt: UnaryFmt, op: op.Negate},
	{code: Swap, fmt: CompareFmt, op: op.SwapIn},
	{code: IfGt, fmt: CompareFmt},
	{code: IfLe, fmt: CompareFmt},
	{code: IfUgt, fmt: CompareFmt},
	{code: IfUle, fmt: CompareFmt},
	{code: Return, fmt: NoFmt},
}

var translations []Opcode

// sort opDefs so we don't have to worry about that
func init() {
	var newdefs []def
	maxCode := Nop
	for _, op := range opDefs {
		if op.code > maxCode {
			maxCode = op.code
		}
	}
	newdefs = make([]def, maxCode+1)
	translations = make([]Opcode, op.NumOps)
	for _, op := range opDefs {
		newdefs[op.code] = op
		translations[op.op] = op.code
	}
	opDefs = newdefs

	for _, op := range OpcodeValues() {
		if newdefs[op].code != op {
			panic("missing OpDef for " + op.String())
		}
	}
}
