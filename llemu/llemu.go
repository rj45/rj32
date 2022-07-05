package main

import (
	"flag"
	"os"
	"time"

	"github.com/rj45/rj32/llemu/cpu"
	"github.com/rj45/rj32/llemu/mem"
	"github.com/rj45/rj32/llemu/vcd"
)

func main() {
	vcd := flag.String("vcd", "", "output a vcd file for gtkwave viewing")

	flag.Parse()

	filename := flag.Arg(0)

	run(filename, *vcd)
}

func run(prog, vcdFile string) {

	ram := mem.NewMemory(24)

	{
		buf, err := os.ReadFile(prog)
		if err != nil {
			panic(err)
		}
		ram.Load(0, buf)
	}

	mmap := mem.MemMap{
		mem.Device{Address: 0x80000000, Size: 1 << 24, Handler: ram},
	}

	cpu := cpu.CPU{BusHandler: mmap}

	var vcdOut *vcd.VCD

	if vcdFile != "" {
		vcdf, err := os.Create(vcdFile)
		if err != nil {
			panic(err)
		}
		defer vcdf.Close()
		vcdOut = vcd.NewVCD(vcdf)

		cpu.AddVCDVariables(vcdOut)
		vcdOut.WriteHeader(1*time.Nanosecond, "v1")
	}

	cpu.Jump(0x80000000)

	for i := 0; i < 100; i++ {
		cpu.Run()
		cpu.ClockRegisters()

		if vcdOut != nil {
			vcdOut.Step(1)
		}
	}
}
