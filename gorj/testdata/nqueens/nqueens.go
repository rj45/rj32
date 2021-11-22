package main

import "unsafe"

// from: http://caseylmanus.github.io/golang/eightqueenspuzzle/2014/11/18/Eight-Queens-in-Go-(GoLang).html

const N = 4

type point struct {
	x int
	y int
}

var board [N]point
var num int

func main() {
	for col := 0; col < N; col++ {
		putc('c')
		putdec(col)
		putc('\r')
		putc('\n')
		recurse(col, 0, 0)
	}

	putdec(num)
	putc('\r')
	putc('\n')

	if num != 92 {
		panic(num)
	}
}

func recurse(x, y, n int) {
	board[n].x = x
	board[n].y = y
	n++
	putc('p')
	putdec(n)
	putc(' ')
	putdec(x)
	putc(' ')
	putdec(y)
	putc('\r')
	putc('\n')
	// if n > 2 {
	// 	putc('n')
	// 	putdec(n)
	// 	putc('\r')
	// 	putc('\n')
	// }
	if n == N {
		num++
		printBoard()
	}

	for col := 0; col < N; col++ {
	nexttry:
		for row := y; row < N; row++ {
			for i := 0; i < n; i++ {
				if col == board[i].x {
					continue nexttry
				}

				if row == board[i].y {
					continue nexttry
				}

				ay := board[i].y - row
				if ay < 0 {
					ay = -ay
				}

				ax := board[i].x - col
				if ax < 0 {
					ax = -ax
				}

				if ay == ax {
					continue nexttry
				}
			}

			recurse(col, row, n)
		}
	}
}

// func canPlace(x, y, n int) bool {
// 		if canAttack(x, y, i) {
// 			return false
// 		}
// 	}
// 	return true
// }

// func canAttack(x, y, i int) bool {
// 	return
// }

func printBoard() {
	putdec(num)
	putc(':')
	putc(' ')
	for n := 0; n < len(board); n++ {
		putc('(')
		putdec(board[n].x)
		putc(',')
		putdec(board[n].y)
		putc(')')
		putc(' ')
	}
	putc('\r')
	putc('\n')
}

const consoleAddr = 0xFF00

func putc(c rune) {
	*(*int)(unsafe.Pointer(uintptr(consoleAddr))) = int(c)
	// switch c {
	// case '\r':
	// case '\n':
	// 	fmt.Println()
	// default:
	// 	fmt.Printf("%c", c)
	// }
}

func putdec(x int) {
	i := x

	if x >= 10000 {
		j := 0
		for i >= 10000 {
			i -= 10000
			j++
		}
		putc('0' + rune(j))
	}

	if x >= 1000 {
		j := 0
		for i >= 1000 {
			i -= 1000
			j++
		}
		putc('0' + rune(j))
	}

	if x >= 100 {
		j := 0
		for i >= 100 {
			i -= 100
			j++
		}
		putc('0' + rune(j))
	}

	if x >= 10 {
		j := 0
		for i >= 10 {
			i -= 10
			j++
		}
		putc('0' + rune(j))
	}

	putc('0' + rune(i))
}
