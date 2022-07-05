package wb

type In struct {
	PC uint32
	IR uint32
}

type Out struct {
	PC     uint32
	Result uint32
}

type Stage struct {
	In
	Out
}

func (s *Stage) Run() {
	s.Out.PC = s.In.PC
}
