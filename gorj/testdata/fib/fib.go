package main

func Fibonacci(n uint16) uint16 {
	if n <= 1 {
		return n
	}

	var n2, n1 uint16 = 0, 1

	for i := uint16(2); i < n; i++ {
		n2, n1 = n1, n1+n2
	}

	return n2 + n1
}

func main() {
}
