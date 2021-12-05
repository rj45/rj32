package a32

func (cpuArch) AssemblerFormat() string {
	return "binary"
}

func (cpuArch) EmulatorCmd() string {
	return "a32emu"
}

func (cpuArch) EmulatorArgs() []string {
	return []string{"--headless", "--max-cycles", "1000000000", "--rom"}
}
