package cpu

type VCD interface {
	AddBit(module, name string, bit *bool)
	AddWire(module, name string, width int, wire *uint32)
	AddStr(module, name string, str *string)
}

func (c *CPU) AddVCDVariables(v VCD) {
	v.AddWire("fe.Out", "PC", 32, &c.fe.Out.Pc)
	v.AddWire("fe.Out", "IR", 32, &c.fe.Out.IR)
	v.AddBit("fe.Out", "Bus.Ack", &c.fe.Out.Bus.Ack)
	v.AddWire("fe.Out", "Bus.Address", 32, &c.fe.Out.Bus.Address)
}
