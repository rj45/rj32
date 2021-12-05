package main

func printstr(s string) int {
	num := 0
	for _, c := range s {
		print(c)
		num++
	}

	return num
}

func main() {
	num := printstr("Hello, World!\r\n")
	if num != 15 {
		panic(num)
	}
	println(num)
}
