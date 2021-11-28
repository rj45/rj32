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

func muls(a, b int) (res int) {
	sign := false
	if a < 0 {
		a = -a
		sign = true
	}

	if b < 0 {
		b = -b
		sign = !sign
	}

	for b > 0 {
		if b&1 != 0 {
			res += a
		}
		a <<= 1
		b >>= 1
	}

	if sign {
		res = -res
	}

	return
}

func divu(dividend, divisor uint) (quotient, remainder uint) {
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

	// sres := muls(-3, 9)
	// if sres != -3*9 {
	// 	panic(sres)
	// }

	quo, rem := divu(10, 5)
	if quo != 10/5 {
		panic(quo)
	}

	if rem != 10%5 {
		panic(rem)
	}

	quo, rem = divu(1234, 13)
	if quo != 1234/13 {
		panic(quo)
	}

	if rem != 1234%13 {
		panic(rem)
	}
}
