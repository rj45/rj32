package mem

import (
	"bytes"
	"debug/elf"
	"fmt"
	"strconv"
)

func Load(width int, buf []byte, write func(uint64, uint64)) int {
	if bytes.HasPrefix(buf, []byte(elf.ELFMAG)) {
		return loadElf(width, buf, write)
	}

	if bytes.HasPrefix(bytes.TrimSpace(buf), []byte("v2.0 raw")) {
		return loadHex(width, buf, write)
	}

	nbytes := width / 8
	for i := 0; i < len(buf); i += nbytes {
		var val uint64
		shift := 0
		for j := 0; j < nbytes; j++ {
			val |= (uint64(buf[i+j]) & 0xff) << shift
			shift += 8
		}
		write(uint64(i/nbytes), val)
	}

	return len(buf) / nbytes
}

func loadHex(width int, buf []byte, write func(uint64, uint64)) int {
	trimmed := bytes.TrimPrefix(bytes.TrimSpace(buf), []byte("v2.0 raw"))
	data := bytes.Fields(trimmed)

	nbytes := width / 8
	str := ""
	for i := 0; i < len(data); i++ {
		str = string(data[i]) + str
		if len(str) < nbytes*2 {
			continue
		}
		val, err := strconv.ParseUint(str, 16, width)
		str = ""
		if err != nil {
			panic(err)
		}

		write(uint64(i), val)
	}
	return len(data)
}

func loadElf(width int, buf []byte, write func(uint64, uint64)) int {
	rd := bytes.NewReader(buf)
	file, err := elf.NewFile(rd)
	if err != nil {
		panic(err)
	}

	byteWidth := width / 8
	totalLen := 0

	for _, sec := range file.Sections {
		if sec.Type == elf.SHT_PROGBITS && (sec.Flags&elf.SHF_ALLOC) != 0 {
			addr := sec.Addr

			data, err := sec.Data()
			if err != nil {
				panic(err)
			}

			fmt.Printf("Loading %d bytes from %s at 0x%08x\n", len(data), sec.Name, addr)

			for i := 0; i < len(data)&^3; i += byteWidth {
				var val uint64
				shift := 0
				for j := 0; j < byteWidth; j++ {
					val |= (uint64(data[i+j]) & 0xff) << shift
					shift += 8
				}
				write(addr, val)
				addr += uint64(byteWidth)
				totalLen++
			}
		}
	}

	return 0
}
