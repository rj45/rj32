package main

import "unsafe"

const consoleAddr = 0xFF00

func putc(c rune) {
	*(*int)(unsafe.Pointer(uintptr(consoleAddr))) = int(c)
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

var num int
var primes [0x3ffe]struct {
	prime    int
	multiple int
}

// algorithm from:
// https://www.codevamping.com/2019/01/incremental-sieve-of-eratosthenes/
// modified to avoid a square root and division
func main() {
	limit := 3
	limsquared := 9
	nextodd := 7

	for x := 5; x < 32750 && num < len(primes); x += 2 {
		// approximate sqrt(x)
		for limsquared < x && limsquared > 0 {
			limit++

			// weird trick to avoid calculating limit*limit,
			// the square of `limit` is the sum of `limit` odd numbers
			limsquared += nextodd
			nextodd += 2
		}

		isprime := true
		for i := 0; i < num && primes[i].prime <= limit; i++ {
			// note: be careful of overflow into negative numbers
			for primes[i].multiple < x && primes[i].multiple > 0 {
				primes[i].multiple += primes[i].prime
			}

			if primes[i].multiple == x {
				isprime = false
				break
			}
		}

		if isprime {
			primes[num].prime = x

			// this could be faster if it was x*x instead of x+x, but * is slow.
			primes[num].multiple = x + x

			num++

			putdec(x)
			putc('\r')
			putc('\n')
		}

		putdec(num)
		putc('\r')
		putc('\n')
	}
}
