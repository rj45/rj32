package rj32

//go:generate enumer -transform lower -trimprefix Fmt -type Op,Fmt .

type Fmt int

const (
	FmtRR  Fmt = 0b00
	FmtExt Fmt = 0b01
	FmtLS  Fmt = 0b10
	FmtRI6 Fmt = 0b11
	FmtRI8 Fmt = 0b100
	FmtI11 Fmt = 0b0101
	FmtI12 Fmt = 0b1101
)

type Op int

const (
	Nop Op = iota
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
)
