package main

func mulu(a, b uint) (res uint) {
	for b > 0 {
		if b&1 != 0 {
			res += a
		}
		a <<= 1
		b >>= 1
	}
	return
}

func divu(dividend, divisor uint) (quotient uint) {
	var remainder uint

	if divisor == 0 {
		panic(-255)
	}

	for i := int(15); i >= 0; i-- {
		quotient <<= 1
		remainder <<= 1

		remainder |= (dividend & (1 << uint(i))) >> uint(i)

		if remainder >= divisor {
			remainder -= divisor
			quotient |= 1
		}
	}

	return
}

func remu(dividend, divisor uint) (remainder uint) {
	var quotient uint
	if divisor == 0 {
		panic(-255)
	}

	for i := int(15); i >= 0; i-- {
		quotient <<= 1
		remainder <<= 1

		remainder |= (dividend & (1 << uint(i))) >> uint(i)

		if remainder >= divisor {
			remainder -= divisor
			quotient |= 1
		}
	}

	return
}

func main() {
	res := mulu(2, 5)
	if res != 2*5 {
		panic(res)
	}

	res = mulu(49, 1234)
	if res != 49*1234 {
		panic(res)
	}

	res = divu(10, 5)
	if res != 10/5 {
		panic(res)
	}

	res = remu(7, 2)
	if res != 7%2 {
		panic(res)
	}
}
