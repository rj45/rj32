package data

import (
	"bytes"
	"strconv"
)

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
