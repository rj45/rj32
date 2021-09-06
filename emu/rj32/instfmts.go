// Code generated by github.com/rj45/rj32/emu/bitfield. DO NOT EDIT.

package rj32

type Inst uint32

func (i Inst) Fmt() Fmt {
	return Fmt((i >> 29) & 0x7)
}

func (i Inst) SetFmt(v Fmt) Inst {
	i &= ^Inst(0x7 << 29)
	i |= (Inst(v) & 0x7) << 29
	return i
}

func (i Inst) Op() Op {
	return Op((i >> 24) & 0x1f)
}

func (i Inst) SetOp(v Op) Inst {
	i &= ^Inst(0x1f << 24)
	i |= (Inst(v) & 0x1f) << 24
	return i
}

func (i Inst) Rd() int {
	return int((i >> 20) & 0xf)
}

func (i Inst) SetRd(v int) Inst {
	i &= ^Inst(0xf << 20)
	i |= (Inst(v) & 0xf) << 20
	return i
}

func (i Inst) Rs() int {
	return int((i >> 16) & 0xf)
}

func (i Inst) SetRs(v int) Inst {
	i &= ^Inst(0xf << 16)
	i |= (Inst(v) & 0xf) << 16
	return i
}

func (i Inst) Imm() int {
	return int((i >> 4) & 0xfff)
}

func (i Inst) SetImm(v int) Inst {
	i &= ^Inst(0xfff << 4)
	i |= (Inst(v) & 0xfff) << 4
	return i
}
type InstRI6 uint16

func (i InstRI6) Rd() int {
	return int((i >> 12) & 0xf)
}

func (i InstRI6) SetRd(v int) InstRI6 {
	i &= ^InstRI6(0xf << 12)
	i |= (InstRI6(v) & 0xf) << 12
	return i
}

func (i InstRI6) Imm() int {
	return int((i >> 6) & 0x3f)
}

func (i InstRI6) SetImm(v int) InstRI6 {
	i &= ^InstRI6(0x3f << 6)
	i |= (InstRI6(v) & 0x3f) << 6
	return i
}

func (i InstRI6) Op() Op {
	return Op((i >> 2) & 0xf)
}

func (i InstRI6) SetOp(v Op) InstRI6 {
	i &= ^InstRI6(0xf << 2)
	i |= (InstRI6(v) & 0xf) << 2
	return i
}

func (i InstRI6) Fmt() Fmt {
	return Fmt((i >> 0) & 0x3)
}

func (i InstRI6) SetFmt(v Fmt) InstRI6 {
	i &= ^InstRI6(0x3 << 0)
	i |= (InstRI6(v) & 0x3) << 0
	return i
}
type InstRR uint16

func (i InstRR) Rd() int {
	return int((i >> 12) & 0xf)
}

func (i InstRR) SetRd(v int) InstRR {
	i &= ^InstRR(0xf << 12)
	i |= (InstRR(v) & 0xf) << 12
	return i
}

func (i InstRR) Rs() int {
	return int((i >> 8) & 0xf)
}

func (i InstRR) SetRs(v int) InstRR {
	i &= ^InstRR(0xf << 8)
	i |= (InstRR(v) & 0xf) << 8
	return i
}

func (i InstRR) NA() int {
	return int((i >> 7) & 0x1)
}

func (i InstRR) SetNA(v int) InstRR {
	i &= ^InstRR(0x1 << 7)
	i |= (InstRR(v) & 0x1) << 7
	return i
}

func (i InstRR) Op() Op {
	return Op((i >> 2) & 0x1f)
}

func (i InstRR) SetOp(v Op) InstRR {
	i &= ^InstRR(0x1f << 2)
	i |= (InstRR(v) & 0x1f) << 2
	return i
}

func (i InstRR) Fmt() Fmt {
	return Fmt((i >> 0) & 0x3)
}

func (i InstRR) SetFmt(v Fmt) InstRR {
	i &= ^InstRR(0x3 << 0)
	i |= (InstRR(v) & 0x3) << 0
	return i
}
type InstLS uint16

func (i InstLS) Rd() int {
	return int((i >> 12) & 0xf)
}

func (i InstLS) SetRd(v int) InstLS {
	i &= ^InstLS(0xf << 12)
	i |= (InstLS(v) & 0xf) << 12
	return i
}

func (i InstLS) Rs() int {
	return int((i >> 8) & 0xf)
}

func (i InstLS) SetRs(v int) InstLS {
	i &= ^InstLS(0xf << 8)
	i |= (InstLS(v) & 0xf) << 8
	return i
}

func (i InstLS) Imm() int {
	return int((i >> 4) & 0xf)
}

func (i InstLS) SetImm(v int) InstLS {
	i &= ^InstLS(0xf << 4)
	i |= (InstLS(v) & 0xf) << 4
	return i
}

func (i InstLS) Op() Op {
	return Op((i >> 2) & 0x3)
}

func (i InstLS) SetOp(v Op) InstLS {
	i &= ^InstLS(0x3 << 2)
	i |= (InstLS(v) & 0x3) << 2
	return i
}

func (i InstLS) Fmt() Fmt {
	return Fmt((i >> 0) & 0x3)
}

func (i InstLS) SetFmt(v Fmt) InstLS {
	i &= ^InstLS(0x3 << 0)
	i |= (InstLS(v) & 0x3) << 0
	return i
}
type InstI11 uint16

func (i InstI11) Imm() int {
	return int((i >> 5) & 0x7ff)
}

func (i InstI11) SetImm(v int) InstI11 {
	i &= ^InstI11(0x7ff << 5)
	i |= (InstI11(v) & 0x7ff) << 5
	return i
}

func (i InstI11) Op() Op {
	return Op((i >> 3) & 0x3)
}

func (i InstI11) SetOp(v Op) InstI11 {
	i &= ^InstI11(0x3 << 3)
	i |= (InstI11(v) & 0x3) << 3
	return i
}

func (i InstI11) Fmt() Fmt {
	return Fmt((i >> 0) & 0x7)
}

func (i InstI11) SetFmt(v Fmt) InstI11 {
	i &= ^InstI11(0x7 << 0)
	i |= (InstI11(v) & 0x7) << 0
	return i
}
type InstRI8 uint16

func (i InstRI8) Rd() int {
	return int((i >> 12) & 0xf)
}

func (i InstRI8) SetRd(v int) InstRI8 {
	i &= ^InstRI8(0xf << 12)
	i |= (InstRI8(v) & 0xf) << 12
	return i
}

func (i InstRI8) Imm() int {
	return int((i >> 4) & 0xff)
}

func (i InstRI8) SetImm(v int) InstRI8 {
	i &= ^InstRI8(0xff << 4)
	i |= (InstRI8(v) & 0xff) << 4
	return i
}

func (i InstRI8) Op() Op {
	return Op((i >> 3) & 0x1)
}

func (i InstRI8) SetOp(v Op) InstRI8 {
	i &= ^InstRI8(0x1 << 3)
	i |= (InstRI8(v) & 0x1) << 3
	return i
}

func (i InstRI8) Fmt() Fmt {
	return Fmt((i >> 0) & 0x7)
}

func (i InstRI8) SetFmt(v Fmt) InstRI8 {
	i &= ^InstRI8(0x7 << 0)
	i |= (InstRI8(v) & 0x7) << 0
	return i
}
type Bus uint64

func (b Bus) Data() int {
	return int((b >> 48) & 0xffff)
}

func (b Bus) SetData(v int) Bus {
	b &= ^Bus(0xffff << 48)
	b |= (Bus(v) & 0xffff) << 48
	return b
}

func (b Bus) Address() int {
	return int((b >> 27) & 0x1fffff)
}

func (b Bus) SetAddress(v int) Bus {
	b &= ^Bus(0x1fffff << 27)
	b |= (Bus(v) & 0x1fffff) << 27
	return b
}

func (b Bus) WE() bool {
	const bit = 1 << 26
	return b&bit == bit
}

func (b Bus) SetWE(v bool) Bus {
	const bit = 1 << 26
	if v {
		return b | bit
	}
	return b & ^Bus(bit)
}

func (b Bus) Req() bool {
	const bit = 1 << 25
	return b&bit == bit
}

func (b Bus) SetReq(v bool) Bus {
	const bit = 1 << 25
	if v {
		return b | bit
	}
	return b & ^Bus(bit)
}

func (b Bus) Ack() bool {
	const bit = 1 << 24
	return b&bit == bit
}

func (b Bus) SetAck(v bool) Bus {
	const bit = 1 << 24
	if v {
		return b | bit
	}
	return b & ^Bus(bit)
}

func (b Bus) Conflict() bool {
	const bit = 1 << 23
	return b&bit == bit
}

func (b Bus) SetConflict(v bool) Bus {
	const bit = 1 << 23
	if v {
		return b | bit
	}
	return b & ^Bus(bit)
}