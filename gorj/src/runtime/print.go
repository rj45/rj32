package runtime

func putc(c byte)

//go:nobounds
func printstring(s string) {
	for i := 0; i < len(s); i++ {
		putc(s[i])
	}
}

func printnl() {
	putc('\r')
	putc('\n')
}

func printrune(x rune) {
	putc(byte(x))
}

func printint(x int) {
	i := x

	// todo: simplify this when the relavent ops are implemented

	if x >= 10000 {
		j := 0
		for i >= 10000 {
			i -= 10000
			j++
		}
		putc('0' + byte(j))
	}

	if x >= 1000 {
		j := 0
		for i >= 1000 {
			i -= 1000
			j++
		}
		putc('0' + byte(j))
	}

	if x >= 100 {
		j := 0
		for i >= 100 {
			i -= 100
			j++
		}
		putc('0' + byte(j))
	}

	if x >= 10 {
		j := 0
		for i >= 10 {
			i -= 10
			j++
		}
		putc('0' + byte(j))
	}

	putc('0' + byte(i))
}
