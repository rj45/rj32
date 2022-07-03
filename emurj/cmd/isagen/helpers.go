package main

func opcodeToFmt(opcode uint) string {
	for _, matcher := range opcodeFmtMatchers {
		if (opcode & matcher.mask) == matcher.match {
			return matcher.fmt
		}
	}
	panic("missing catch all")
}
