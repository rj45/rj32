package main

import "unsafe"

const consoleAddr = 0xFF00

func main() {
	console := (*rune)(unsafe.Pointer(uintptr(consoleAddr)))
	*console = 'H'
	*console = 'e'
	*console = 'l'
	*console = 'l'
	*console = 'o'
	*console = ','
	*console = ' '
	*console = 'W'
	*console = 'o'
	*console = 'r'
	*console = 'l'
	*console = 'd'
	*console = '!'
	*console = '\r'
	*console = '\n'
}
