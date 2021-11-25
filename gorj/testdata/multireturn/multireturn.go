package main

func testreturns(a, b, c int) (int, int, int) {
	return a + 1, b + 2, c + 3
}

func main() {
	a, b, c := testreturns(2, 3, 4)

	if a != 2+1 {
		panic(a)
	}

	if b != 3+2 {
		panic(b)
	}

	if c != 4+3 {
		panic(c)
	}
}
