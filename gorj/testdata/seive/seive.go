package main

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

nextnum:
	for x := 3; x < 32750 && num < len(primes); x += 2 {
		// approximate sqrt(x)
		for limsquared < x && limsquared > 0 {
			limit++

			// weird trick to avoid calculating limit*limit,
			// the square of `limit` is the sum of `limit` odd numbers
			limsquared += nextodd
			nextodd += 2
		}

		for i := 0; i < num; i++ {
			if primes[i].prime > limit {
				break
			}

			// note: be careful of overflow into negative numbers
			for primes[i].multiple < x && primes[i].multiple > 0 {
				primes[i].multiple += primes[i].prime
			}

			if primes[i].multiple == x {
				continue nextnum
			}
		}

		primes[num].prime = x

		// this could be faster if it was x*x instead of x+x, but * is slow.
		primes[num].multiple = x + x

		num++

		println(x)
	}
	println(num)
}
