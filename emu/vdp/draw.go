package vdp

// drawLineBuffer takes pixels from the line buffer, looks up their
// palette entry and draws them to the screen.
func (vdp *VDP) drawLineBuffer(y int, framebuf []byte) {
	fbYOffset := y * ScreenWidth * 4
	for x := 0; x < ScreenWidth/4; x++ {
		fbXOffset := fbYOffset + x*4*4

		// line buffer is 4 block RAMs
		for i := 0; i < 4; i++ {
			// look up color in the palette
			col := vdp.Palette[int(vdp.LineBuf[x][i])]

			// clear previously read pixel in next clock cycle
			vdp.LineBuf[x][i] = 0

			// this is done differently in the circuit -- this is specific
			// to the graphics library in use
			framebuf[fbXOffset+(i*4)+0] = col.R
			framebuf[fbXOffset+(i*4)+1] = col.G
			framebuf[fbXOffset+(i*4)+2] = col.B
			framebuf[fbXOffset+(i*4)+3] = 255
		}
	}
}

// findScanLineSprites scans all the sprites looking
// for ones intersecting the current scan line.
func (vdp *VDP) findScanLineSprites(y int) {
	for sid := 0; sid < vdp.NumRenderedSprites; sid++ {
		sy := (y - int(vdp.Y[sid]))
		dim := vdp.Dims[sid]

		sh := vdp.Dims[sid].Height() << 3

		// TODO: handle Y flipping

		if uint(sy) < uint(sh) {
			// in the circuit, this would go into a FIFO or block RAM so
			// that the scanner can work ahead while other circuits
			// are busy
			vdp.scanSpriteLine(sid, sy, dim.Width())
		}
	}
}

// scanSpriteLine scans through the width of the sprite looking
// for any tiles that are visible on the screen
func (vdp *VDP) scanSpriteLine(id, y, width int) {
	ty := y >> 3 // y divided by 8
	iy := y & 7

	ss := vdp.Addr[id]
	sp := vdp.SPos[id]

	startX := vdp.X[id]

	sheetY := ty + sp.SheetY()

	// sheet offset is which 1kw block we are in
	sheetOff := ss.SheetAddr() << 10

	// base address of the attribute for this line
	attrbase := sheetOff + ((sheetY & 0x7f) << 7)

	// base address of the tileset for this line
	addr := (ss.TileSetAddr() << 12) | iy<<1

	startSheetX := sp.SheetX()

	// startSheetX (7),

	for tx := 0; tx < width; tx++ {
		sheetX := (tx + startSheetX) & 0x7f
		lineX := (tx<<3 + int(startX))

		// tiles on the right wrap around and are partially
		// visible on the left, so adjust for that
		lineXm8 := (lineX + 8) & 0x3ff

		// skip drawing the tile if it's not visible
		if uint(lineXm8) > uint(ScreenWidth+8) {
			continue
		}

		attraddr := attrbase | sheetX

		// goes into a fifo here
		vdp.drawTile(attraddr, addr, lineX, ss.Transparent(), ss.PaletteSet())
	}
}

// drawTile takes a queued sprite tile, loads it from memory,
// and draws it 4 pixels at a time to the line buffer
func (vdp *VDP) drawTile(attraddr, addr, lineX int, transparent bool, palset int) {
	// read the tile attribute from the sprite sheet
	tile := Tile(vdp.Mem.Read(attraddr))
	pal := uint16((tile.Palette() << 4) | (palset << 9))

	addr += tile.TileID() << 4
	for ix := 0; ix < 2; ix++ {
		// load the pixels from the tile data
		px := Pixels(vdp.Mem.Read(addr + ix))

		// determine the line buffer addresses
		// this can wrap around back to zero mid-tile
		lineaddr0 := (lineX + (ix << 2) + 0) & 0x3ff
		lineaddr1 := (lineX + (ix << 2) + 1) & 0x3ff
		lineaddr2 := (lineX + (ix << 2) + 2) & 0x3ff
		lineaddr3 := (lineX + (ix << 2) + 3) & 0x3ff

		// in the circuit the following can be done simulateously to
		// 4 different block rams

		// pixel is written to line buffer if its not (transparent and zero)
		if !(transparent && px.Color0() == 0) {
			vdp.LineBuf[lineaddr0>>2][lineaddr0&3] = pal | uint16(px.Color0())
		}

		if !(transparent && px.Color1() == 0) {
			vdp.LineBuf[lineaddr1>>2][lineaddr1&3] = pal | uint16(px.Color1())
		}

		if !(transparent && px.Color2() == 0) {
			vdp.LineBuf[lineaddr2>>2][lineaddr2&3] = pal | uint16(px.Color2())
		}

		if !(transparent && px.Color3() == 0) {
			vdp.LineBuf[lineaddr3>>2][lineaddr3&3] = pal | uint16(px.Color3())
		}
	}
}
