package main

// from: http://caseylmanus.github.io/golang/eightqueenspuzzle/2014/11/18/Eight-Queens-in-Go-(GoLang).html

const N = 8

type point struct {
	x int
	y int
}

var board [N]point
var num int

func main() {
	for col := 0; col < N; col++ {
		recurse(col, 0, 0)
	}

	println(num)

	if num != 92 {
		panic(num)
	}
}

func recurse(x, y, n int) {
	board[n].x = x
	board[n].y = y
	n++
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

func printBoard() {
	print(num, ": ")
	for n := 0; n < len(board); n++ {
		print("(", board[n].x, ",", board[n].y, ") ")
	}
	println(" ")
}
