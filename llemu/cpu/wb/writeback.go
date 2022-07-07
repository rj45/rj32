package wb

import "github.com/rj45/rj32/llemu/rj32"

type In struct {
	PC uint32
	IR rj32.Inst

	Result uint32
}

type Out struct {
	PC uint32

	Result uint32
	Wen    bool
	Rd     rj32.Reg
}

type Stage struct {
	In
	Out
}

func (s *Stage) Run() {
	s.Out.PC = s.In.PC

	s.Out.Rd = s.In.IR.Rd()
}
