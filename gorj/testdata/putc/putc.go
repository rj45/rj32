package main

import (
	"unsafe"
)

const consoleAddr = 0xFF00

func putc(c rune) {
	*(*rune)(unsafe.Pointer(uintptr(consoleAddr))) = c
}

func main() {
	putc('H')
	putc('e')
	putc('l')
	putc('l')
	putc('o')
	putc(',')
	putc(' ')
	putc('W')
	putc('o')
	putc('r')
	putc('l')
	putc('d')
	putc('!')
	putc('\r')
	putc('\n')
}
