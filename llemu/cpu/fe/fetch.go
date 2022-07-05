package fe

import "github.com/rj45/rj32/llemu/mem"

type In struct {
	PC uint32

	DeStall bool

	Bus mem.Bus
}

type Out struct {
	PC uint32
	IR uint32

	Stall bool

	Bus mem.Bus
}

type Stage struct {
	In
	Out
}

func (s *Stage) Run() {
	s.Out.Bus.Mask = 0xffff_ffff
	s.Out.Bus.WE = false
	s.Out.Bus.Ack = false
	s.Out.Bus.Address = s.In.PC

	if s.In.Bus.Ack && !s.In.DeStall {
		s.Out.Stall = false

		s.Out.PC = s.In.PC + 4
		s.Out.IR = s.In.Bus.Data
	} else {
		s.Out.Stall = true
		s.Out.PC = s.In.PC
		s.Out.IR = 0 // reset
	}
}
