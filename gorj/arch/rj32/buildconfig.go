package rj32

func (cpuArch) AssemblerFormat() string {
	return "logisim16"
}

func (cpuArch) EmulatorCmd() string {
	return "emurj"
}

func (cpuArch) EmulatorArgs() []string {
	return []string{"-run"}
}
