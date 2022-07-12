package ex

import (
	"bytes"
	"fmt"

	"github.com/rj45/rj32/llemu/rj32"
)

type In struct {
	PC uint32
	IR rj32.Inst

	L uint32
	R uint32

	AluOp rj32.Opcode
	Sub   bool

	RegWen bool
	Rd     rj32.Reg

	Log *bytes.Buffer
}

type Out struct {
	PC uint32
	IR rj32.Inst

	Result uint32
	RegWen bool
	Rd     rj32.Reg
}

type Stage struct {
	In
	Out
}

func (s *Stage) Run() {
	s.Out.PC = s.In.PC
	s.Out.IR = s.In.IR

	s.Out.RegWen = s.In.RegWen
	s.Out.Rd = s.In.Rd

	switch s.AluOp {
	case rj32.Add:
		if s.In.Sub {
			s.Out.Result = s.In.L - s.In.R
		} else {
			s.Out.Result = s.In.L + s.In.R
		}
	case rj32.Slt:
		if int32(s.In.L) < int32(s.In.R) {
			s.Out.Result = 1
		} else {
			s.Out.Result = 0
		}
	case rj32.Sltu:
		if uint32(s.In.L) < uint32(s.In.R) {
			s.Out.Result = 1
		} else {
			s.Out.Result = 0
		}
	case rj32.And:
		s.Out.Result = s.In.L & s.In.R
	case rj32.Or:
		s.Out.Result = s.In.L | s.In.R
	case rj32.Xor:
		s.Out.Result = s.In.L ^ s.In.R
	case rj32.Srl:
		if s.In.Sub {
			s.Out.Result = uint32(int32(s.In.L) >> int32(s.In.R))
		} else {
			s.Out.Result = s.In.L >> s.In.R
		}
	case rj32.Sll:
		s.Out.Result = s.In.L << s.In.R
	}

	fmt.Fprintf(s.In.Log, "  ex: %d %s %d = %#08x\n", s.In.L, s.In.AluOp, s.In.R, s.Out.Result)
}
