package wb

import (
	"bytes"
	"fmt"

	"github.com/rj45/rj32/llemu/rj32"
)

type In struct {
	PC uint32
	IR rj32.Inst

	Result uint32
	RegWen bool
	Rd     rj32.Reg

	Log *bytes.Buffer
}

type Out struct {
	PC uint32
}

type RegWriter interface {
	WriteReg(reg rj32.Reg, val uint32)
}

type Stage struct {
	In
	Out

	Writer RegWriter
}

func (s *Stage) Run() {
	s.Out.PC = s.In.PC

	if s.In.RegWen {
		s.Writer.WriteReg(s.In.Rd, s.In.Result)
		fmt.Fprintf(s.In.Log, "  wb: %s <- %#08x (%d)\n", s.In.Rd, s.In.Result, s.In.Result)
	}
}
