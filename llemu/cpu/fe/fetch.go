package fe

import (
	"bytes"

	"github.com/rj45/rj32/llemu/mem"
)

type In struct {
	Pc  uint32
	Npc uint32

	Bus mem.Bus

	Log *bytes.Buffer
}

type Out struct {
	Pc  uint32
	Npc uint32
	IR  uint32

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
	s.Out.Bus.Address = s.In.Npc

	if s.In.Bus.Ack {
		s.Out.Pc = s.In.Npc
		s.Out.IR = s.In.Bus.Data
	} else {
		s.Out.Pc = s.In.Pc
		s.Out.IR = 0 // nop
	}

	s.Out.Npc = s.In.Pc + 4
}
