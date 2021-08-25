// Code generated by github.com/rj45/rj32/emu/bitfield. DO NOT EDIT.

package vdp

type SpriteDims uint16

func (s SpriteDims) FlipY() bool {
	const bit = 1 << 15
	return s&bit == bit
}

func (s SpriteDims) SetFlipY(v bool) SpriteDims {
	const bit = 1 << 15
	if v {
		return s | bit
	}
	return s & ^SpriteDims(bit)
}

func (s SpriteDims) Height() int {
	return int((s >> 8) & 0x7f)
}

func (s SpriteDims) SetHeight(v int) SpriteDims {
	s &= ^SpriteDims(0x7f << 8)
	s |= (SpriteDims(v) & 0x7f) << 8
	return s
}

func (s SpriteDims) FlipX() bool {
	const bit = 1 << 7
	return s&bit == bit
}

func (s SpriteDims) SetFlipX(v bool) SpriteDims {
	const bit = 1 << 7
	if v {
		return s | bit
	}
	return s & ^SpriteDims(bit)
}

func (s SpriteDims) Width() int {
	return int((s >> 0) & 0x7f)
}

func (s SpriteDims) SetWidth(v int) SpriteDims {
	s &= ^SpriteDims(0x7f << 0)
	s |= (SpriteDims(v) & 0x7f) << 0
	return s
}

type SpriteAddr uint16

func (s SpriteAddr) Transparent() bool {
	const bit = 1 << 15
	return s&bit == bit
}

func (s SpriteAddr) SetTransparent(v bool) SpriteAddr {
	const bit = 1 << 15
	if v {
		return s | bit
	}
	return s & ^SpriteAddr(bit)
}

func (s SpriteAddr) PaletteSet() int {
	return int((s >> 14) & 0x1)
}

func (s SpriteAddr) SetPaletteSet(v int) SpriteAddr {
	s &= ^SpriteAddr(0x1 << 14)
	s |= (SpriteAddr(v) & 0x1) << 14
	return s
}

func (s SpriteAddr) TileSetAddr() int {
	return int((s >> 8) & 0x3f)
}

func (s SpriteAddr) SetTileSetAddr(v int) SpriteAddr {
	s &= ^SpriteAddr(0x3f << 8)
	s |= (SpriteAddr(v) & 0x3f) << 8
	return s
}

func (s SpriteAddr) SheetAddr() int {
	return int((s >> 0) & 0xff)
}

func (s SpriteAddr) SetSheetAddr(v int) SpriteAddr {
	s &= ^SpriteAddr(0xff << 0)
	s |= (SpriteAddr(v) & 0xff) << 0
	return s
}

type SheetPos uint16

func (s SheetPos) WrapY() bool {
	const bit = 1 << 15
	return s&bit == bit
}

func (s SheetPos) SetWrapY(v bool) SheetPos {
	const bit = 1 << 15
	if v {
		return s | bit
	}
	return s & ^SheetPos(bit)
}

func (s SheetPos) SheetY() int {
	return int((s >> 8) & 0x7f)
}

func (s SheetPos) SetSheetY(v int) SheetPos {
	s &= ^SheetPos(0x7f << 8)
	s |= (SheetPos(v) & 0x7f) << 8
	return s
}

func (s SheetPos) WrapX() bool {
	const bit = 1 << 7
	return s&bit == bit
}

func (s SheetPos) SetWrapX(v bool) SheetPos {
	const bit = 1 << 7
	if v {
		return s | bit
	}
	return s & ^SheetPos(bit)
}

func (s SheetPos) SheetX() int {
	return int((s >> 0) & 0x7f)
}

func (s SheetPos) SetSheetX(v int) SheetPos {
	s &= ^SheetPos(0x7f << 0)
	s |= (SheetPos(v) & 0x7f) << 0
	return s
}
