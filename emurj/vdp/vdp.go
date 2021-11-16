package vdp

import (
	"image/color"

	"github.com/rj45/rj32/emurj/data"
)

const (
	NumColors    = 1024
	BlockRamSize = 256
	NumSprites   = BlockRamSize
	LineBufWidth = BlockRamSize
	ScreenWidth  = 640
	ScreenHeight = 360
)

const (
	BlockX = iota
	BlockY
	BlockDims
	BlockAddr
	BlockSPos
	BlockRG
	BlockBA
)

// MemMap is a map of where sheets and sets are stored
type MemMap struct {
	SheetAddr []int
	SetAddr   []int
	NextAddr  int

	PaletteAddr int
	SpriteAddr  int
}

// VDP is an emulated Video Display Processor
type VDP struct {
	Mem    *data.Memory
	MemMap MemMap

	// block rams
	Palette [NumColors]color.RGBA
	X       [NumSprites]int16
	Y       [NumSprites]int16
	Dims    [NumSprites]SpriteDims
	Addr    [NumSprites]SpriteAddr
	SPos    [NumSprites]SheetPos
	LineBuf [LineBufWidth][4]uint16

	// state
	NumRenderedSprites int
}

// NewVDP creats a new VDP
func NewVDP() *VDP {
	return &VDP{
		Mem: data.NewMemory(18),
	}
}

// DrawFrame draws a frame of video data onto a frame buffer
func (vdp *VDP) DrawFrame(framebuf []byte) {
	for y := 0; y < ScreenHeight; y++ {
		// both of these happen at the simultaneously in the circuit
		vdp.findScanLineSprites(y)
		vdp.drawLineBuffer(y, framebuf)
	}
}

func (vdp *VDP) HandleBus(bus data.Bus) data.Bus {
	if !bus.WE() {
		// let the shadowed RAM send the ack
		return bus
	}

	addr := bus.Address()
	index := addr & 0xff
	block := (addr >> 8) & 7
	data := bus.Data()

	switch block {
	case BlockX:
		vdp.X[index] = int16(data)
	case BlockY:
		vdp.Y[index] = int16(data)
	case BlockDims:
		vdp.Dims[index] = SpriteDims(data)
	case BlockAddr:
		vdp.Addr[index] = SpriteAddr(data)
	case BlockSPos:
		vdp.SPos[index] = SheetPos(data)
	case BlockRG:
		vdp.Palette[(index<<1)+0].R = uint8(data)
		vdp.Palette[(index<<1)+1].G = uint8(data >> 8)
	case BlockBA:
		vdp.Palette[(index<<1)+0].B = uint8(data)
		// not connected:
		vdp.Palette[(index<<1)+1].A = 0xff
	}

	// let the shadowed RAM send the ack
	return bus
}

// SetSpriteSheet sets the sheet and tile set address
// for a sprite from the memory map
func (vdp *VDP) SetSpriteSheetSet(sid, sheet, tileset int) {
	vdp.Addr[sid] = vdp.Addr[sid].
		SetSheetAddr(vdp.MemMap.SheetAddr[sheet] >> 10).
		SetTileSetAddr(vdp.MemMap.SetAddr[tileset] >> 12)
}

// ClearSprites zeros out the specified sprites
func (vdp *VDP) ClearSprites(start, len int) {
	for i := start; i < len+start; i++ {
		vdp.X[i] = 0
		vdp.Y[i] = 0
		vdp.Addr[i] = 0
		vdp.Dims[i] = 0
		vdp.SPos[i] = 0
	}
}

// SetSpritePos sets the x, y position of the sprite
func (vdp *VDP) SetSpritePos(sid, x, y int) {
	vdp.X[sid] = int16(x)
	vdp.Y[sid] = int16(y)
}

// IncSpritePos sets the x, y position of the sprite
func (vdp *VDP) IncSpritePos(sid, x, y int) {
	vdp.X[sid] += int16(x)
	vdp.Y[sid] += int16(y)
}

// SetSpriteSheetPos sets the x, y position of the sprite in
// the sprite sheet
func (vdp *VDP) SetSpriteSheetPos(sid, x, y int) {
	vdp.SPos[sid] = vdp.SPos[sid].SetSheetX(x).SetSheetY(y)
}

// SetSpriteDims sets the width and height of the sprite
func (vdp *VDP) SetSpriteDims(sid, w, h int) {
	vdp.Dims[sid] = vdp.Dims[sid].SetWidth(w).SetHeight(h)
}
