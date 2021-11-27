package main

import (
	"unsafe"
)

const consoleAddr = 0xFF00

func putc(c rune) {
	*(*rune)(unsafe.Pointer(uintptr(consoleAddr))) = c
}

func print(s string) int {
	num := 0
	for _, c := range s {
		putc(c)
		num++
	}

	return num
}

func main() {
	num := print("Hello, World!\r\n")
	if num != 16 {
		panic(num)
	}
}
