package main

import "flag"

func main() {
	cpudef := flag.String("cpudef", "", "Path of cpudef.asm to overwrite")
	flag.Parse()

	if *cpudef != "" {
		genCpudef(*cpudef)
	}

}
