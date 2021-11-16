package vdp

import (
	"fmt"
	"image/color"
	"os"

	"github.com/rj45/rj32/emurj/data"
)

// ResetMemMap resets the memory map so loading starts at zero
func (vdp *VDP) ResetMemMap() {
	mm := &vdp.MemMap
	mm.NextAddr = 0x2000
	mm.SheetAddr = mm.SheetAddr[:0]
	mm.SetAddr = mm.SetAddr[:0]
}

// LoadSheetSets loads a set of sprite sheets and tile sets
func (vdp *VDP) LoadSheetSets(sheets, tiles string) error {
	mm := &vdp.MemMap

	for i := 0; i < 1000; i++ {
		fn := fmt.Sprintf(sheets, i)
		addr, err := vdp.LoadFile(fn, 1024)
		if err != nil {
			if i == 0 {
				return err
			}
			break
		}

		fmt.Println("Loaded", fn, "at", addr)
		mm.SheetAddr = append(mm.SheetAddr, addr)
	}

	for i := 0; i < 1000; i++ {
		fn := fmt.Sprintf(tiles, i)
		addr, err := vdp.LoadFile(fn, 4096)
		if err != nil {
			if i == 0 {
				return err
			}
			break
		}

		fmt.Println("Loaded", fn, "at", addr)
		mm.SetAddr = append(mm.SetAddr, addr)
	}

	return nil
}

// LoadFile loads a file into memory at NextAddr with
// the given address granularity (will round NextAddr to
// the nearest even multiple of addrround)
func (vdp *VDP) LoadFile(filename string, addrround int) (int, error) {
	mm := &vdp.MemMap

	// round the start address to the next valid address
	granules := mm.NextAddr / addrround
	if (mm.NextAddr % addrround) > 0 {
		granules++
	}
	start := granules * addrround

	bytes, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}
	words := vdp.Mem.Load(start, bytes)
	mm.NextAddr = start + words

	// return the start address
	return start, nil
}

// LoadPalette loads a palette from a file
func (vdp *VDP) LoadPalette(filename string) error {
	bytes, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	data.Load(24, bytes, func(i int, val uint64) {
		vdp.Palette[i] = color.RGBA{
			R: uint8(val >> 16),
			G: uint8(val >> 8),
			B: uint8(val),
		}
	})
	return nil
}
