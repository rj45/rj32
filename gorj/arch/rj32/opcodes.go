package rj32

import (
	"fmt"

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

	NumOps
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
	fmt Fmt
	op  op.Op
}

var opDefs = [...]def{
	Nop:    {fmt: NoFmt},
	Rets:   {fmt: NoFmt},
	Error:  {fmt: NoFmt},
	Halt:   {fmt: NoFmt},
	Rcsr:   {fmt: BinaryFmt},
	Wcsr:   {fmt: BinaryFmt},
	Move:   {fmt: MoveFmt, op: op.Copy},
	Loadc:  {fmt: BinaryFmt},
	Jump:   {fmt: CallFmt},
	Imm:    {fmt: UnaryFmt},
	Call:   {fmt: CallFmt, op: op.Call},
	Imm2:   {fmt: BinaryFmt},
	Load:   {fmt: LoadFmt, op: op.Load},
	Store:  {fmt: StoreFmt, op: op.Store},
	Loadb:  {fmt: LoadFmt},
	Storeb: {fmt: StoreFmt},
	Add:    {fmt: BinaryFmt, op: op.Add},
	Sub:    {fmt: BinaryFmt, op: op.Sub},
	Addc:   {fmt: BinaryFmt},
	Subc:   {fmt: BinaryFmt},
	Xor:    {fmt: BinaryFmt, op: op.Xor},
	And:    {fmt: BinaryFmt, op: op.And},
	Or:     {fmt: BinaryFmt, op: op.Or},
	Shl:    {fmt: BinaryFmt, op: op.ShiftLeft},
	Shr:    {fmt: BinaryFmt, op: op.ShiftRight},
	Asr:    {fmt: BinaryFmt},
	IfEq:   {fmt: CompareFmt},
	IfNe:   {fmt: CompareFmt},
	IfLt:   {fmt: CompareFmt},
	IfGe:   {fmt: CompareFmt},
	IfUlt:  {fmt: CompareFmt},
	IfUge:  {fmt: CompareFmt},
	Not:    {fmt: UnaryFmt, op: op.Invert},
	Neg:    {fmt: UnaryFmt, op: op.Negate},
	Swap:   {fmt: CompareFmt, op: op.SwapIn},
	IfGt:   {fmt: CompareFmt},
	IfLe:   {fmt: CompareFmt},
	IfUgt:  {fmt: CompareFmt},
	IfUle:  {fmt: CompareFmt},
	Return: {fmt: NoFmt},
}

var translations []Opcode

func init() {
	translations = make([]Opcode, op.NumOps)
	for i := Nop; i < NumOps; i++ {
		if opDefs[i].fmt == BadFmt {
			panic(fmt.Sprintf("missing opDef for %s", i))
		}
		translations[opDefs[i].op] = i
	}
}

func (cpuArch) IsTwoOperand() bool {
	return true
}
