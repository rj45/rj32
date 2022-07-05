package mem

// Memory is a 32 bit memory. It is byte addressed, but the
// last two bits are ignored and reads and writes are only at the
// word level.
type Memory struct {
	mask uint32
	mem  []uint32
}

// NewMemory creates a new memory with the specified number of
// address bits. Note that the last two bits are ignored.
func NewMemory(addrBits int) *Memory {
	size := 1 << (addrBits - 2)
	mask := uint32(size - 1)
	return &Memory{
		mask: mask,
		mem:  make([]uint32, size),
	}
}

// Read a word at the specified byte address. Last
// two bits are ignored.
func (mem *Memory) Read(address uint32) uint32 {
	return mem.mem[(address>>2)&mem.mask]
}

// Write a word at the specified byte address. Last
// two bits of the address are ignored.
func (mem *Memory) Write(address, data uint32) {
	mem.mem[(address>>2)&mem.mask] = data
}

// WriteMasked writes a masked word. Only the bits in
// the mask will be altered by the write, the rest of
// the bits will be kept the same.
func (mem *Memory) WriteMasked(address, data, mask uint32) {
	mem.mem[(address>>2)&mem.mask] &^= mask
	mem.mem[(address>>2)&mem.mask] |= data & mask
}

// WriteField is similar to WriteMasked, except the mask is
// built from a `bits` number of bits, shifted `shift` amount left.
func (mem *Memory) WriteField(address, bits, shift, data uint32) {
	mask := uint32((1 << bits) - 1)
	val := mem.mem[(address>>2)&mem.mask]
	val &= ^(mask << shift)
	val |= (data & mask) << shift
	mem.mem[(address>>2)&mem.mask] = val
}

// Load the memory from a byte buffer
func (mem *Memory) Load(address uint32, buf []byte) int {
	return Load(32, buf, func(i uint64, d uint64) {
		mem.Write(address+uint32(i), uint32(d))
	})
}

// Clear the memory
func (mem *Memory) Clear(address, len uint32) {
	for i := address; i < (address + len); i++ {
		mem.Write(i, 0)
	}
}

// HandleBus interactions with the memory
func (mem *Memory) HandleBus(bus *Bus) {
	bus.Ack = true

	if bus.WE {
		mem.WriteMasked(bus.Address, bus.Data, bus.Mask)
		return
	}

	bus.Data = mem.Read(bus.Address) & bus.Mask
}
