package ir

type Module struct {
	Funcs   []*Func
	Globals []*Value
}

func (mod *Module) LongString() string {
	str := ""

	for i, fn := range mod.Funcs {
		if i != 0 {
			str += "\n"
		}
		str += fn.LongString()
	}

	return str
}
