package data

type Memory struct {
	mask int
	mem  []uint16
}

func NewMemory(addrBits int) *Memory {
	size := 1 << addrBits
	mask := size - 1
	return &Memory{
		mask: mask,
		mem:  make([]uint16, size),
	}
}

func (mem *Memory) Read(address int) uint16 {
	return mem.mem[address&mem.mask]
}

func (mem *Memory) Write(address int, data uint16) {
	mem.mem[address&mem.mask] = data
}

func (mem *Memory) WriteField(address, bits, shift int, data uint16) {
	mask := uint16((1 << bits) - 1)
	val := mem.mem[address&mem.mask]
	val &= ^(mask << shift)
	val |= (data & mask) << shift
	mem.mem[address&mem.mask] = val
}

func (mem *Memory) Load(address int, buf []byte) int {
	return Load(16, buf, func(i int, d uint64) {
		mem.Write(address+i, uint16(d))
	})
}

func (mem *Memory) Clear(address int, len int) {
	for i := address; i < (address + len); i++ {
		mem.Write(i, 0)
	}
}

func (mem *Memory) HandleBus(bus Bus) Bus {
	addr := bus.Address()

	if bus.WE() {
		mem.Write(addr, uint16(bus.Data()))
		return bus.SetAck(true)
	}

	return bus.
		SetAck(true).
		SetData(int(mem.Read(addr)))
}
