package main

func opcodeNames() []string {
	list := make([]string, len(opmap)*len(opmap[0]))
	for quad := range opmap {
		for i := range opmap[quad] {
			list[quad*len(opmap[quad])+i] = opmap[quad][i]
		}
	}
	return list
}

func opcodeToFmt(opcode uint) string {
	for _, matcher := range opcodeFmtMatchers {
		if (opcode & matcher.mask) == matcher.match {
			return matcher.fmt
		}
	}
	panic("missing catch all")
}
