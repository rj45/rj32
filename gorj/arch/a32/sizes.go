package a32

import "go/types"

var basicSizes = [...]byte{
	types.Bool:       4,
	types.Int:        4,
	types.Int8:       1,
	types.Int16:      2,
	types.Int32:      4,
	types.Int64:      8,
	types.Uint:       4,
	types.Uint8:      1,
	types.Uint16:     2,
	types.Uint32:     4,
	types.Uint64:     8,
	types.Uintptr:    4,
	types.Float32:    4,
	types.Float64:    8,
	types.Complex64:  8,
	types.Complex128: 16,
}

var runeSize = 4

func (cpuArch) BasicSizes() [17]byte {
	return basicSizes
}

func (cpuArch) RuneSize() int {
	return runeSize
}
