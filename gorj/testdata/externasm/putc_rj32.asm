; func(c rune)
putc:
  consoleAddr = 0xFF00
  store [gp, consoleAddr], a0
