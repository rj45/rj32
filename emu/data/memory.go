package data

import (
	"bytes"
	"strconv"
)

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

func Load(width int, buf []byte, write func(int, uint64)) int {
	if bytes.HasPrefix(bytes.TrimSpace(buf), []byte("v2.0 raw")) {
		return loadHex(width, buf, write)
	}
	byteWidth := width / 8
	for i := 0; i < len(buf); i += byteWidth {
		var val uint64
		for j := 0; j < byteWidth; j++ {
			val <<= 8
			val |= uint64(buf[i+j]) & 0xff
		}
		write(i/byteWidth, val)
	}

	return len(buf) / byteWidth
}

func loadHex(width int, buf []byte, write func(int, uint64)) int {
	trimmed := bytes.TrimPrefix(bytes.TrimSpace(buf), []byte("v2.0 raw"))
	data := bytes.Fields(trimmed)
	for i := 0; i < len(data); i++ {
		val, err := strconv.ParseUint(string(data[i]), 16, width)
		if err != nil {
			panic(err)
		}

		write(i, val)
	}
	return len(data)
}
