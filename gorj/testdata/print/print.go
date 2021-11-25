package main

import "unsafe"

const consoleAddr = 0xFF00

func putc(c rune) {
	*(*rune)(unsafe.Pointer(uintptr(consoleAddr))) = c
}

func print(s string) {
	for _, c := range s {
		putc(c)
	}
}

func main() {
	print("Hello, World!\r\n")
}
