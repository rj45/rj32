package main

import "flag"

func main() {
	cpudef := flag.String("cpudef", "", "Path of cpudef.asm to overwrite")
	decoder := flag.String("decoder", "", "Path of decoder.go to overwrite")
	flag.Parse()

	if *cpudef != "" {
		genCpudef(*cpudef)
	}

	if *decoder != "" {
		genDecoder(*decoder)
	}
}
