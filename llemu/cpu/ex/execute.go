package ex

import "github.com/rj45/rj32/llemu/rj32"

type In struct {
	PC uint32
	IR rj32.Inst

	L uint32
	R uint32
}

type Out struct {
	PC uint32
	IR rj32.Inst

	Result uint32
}

type Stage struct {
	In
	Out
}

func (s *Stage) Run() {
	s.Out.PC = s.In.PC
	s.Out.IR = s.In.IR
}
