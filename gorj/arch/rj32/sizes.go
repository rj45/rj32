package rj32

import "go/types"

var basicSizes = [...]byte{
	types.Bool:       1,
	types.Int:        1,
	types.Int8:       1,
	types.Int16:      1,
	types.Int32:      2,
	types.Int64:      4,
	types.Uint:       1,
	types.Uint8:      1,
	types.Uint16:     1,
	types.Uint32:     2,
	types.Uint64:     4,
	types.Uintptr:    1,
	types.Float32:    2,
	types.Float64:    4,
	types.Complex64:  4,
	types.Complex128: 8,
}

var runeSize = 1

func (cpuArch) BasicSizes() [17]byte {
	return basicSizes
}
func (cpuArch) RuneSize() int {
	return runeSize
}
