package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"strconv"

	"github.com/lucasb-eyer/go-colorful"
)

// size of one side of square tiles (ie 8x8)
const tileSize = 8

// number of bits in integers by default
const intSize = 16

// a colour palette which will be used to generate
// various shades of each colour
// var palette = []color.RGBA{
// 	{0x10, 0xB0, 0xE0, 0xff}, // teal
// 	{0xE0, 0x00, 0x70, 0xff}, // red
// 	{0xE0, 0x80, 0x10, 0xff}, // orange
// 	{0x80, 0x70, 0xf0, 0xff}, // purple
// 	{0x10, 0xE0, 0x30, 0xff}, // green
// 	{0xff, 0xff, 0xff, 0xff}, // white
// }

// tron-like palette
var palette = []color.RGBA{
	{0x00, 0xE0, 0xE0, 0xff}, // teal
	{0x00, 0x70, 0xE0, 0xff}, // blue
	{0xE0, 0x70, 0x00, 0xff}, // orange
	{0xE0, 0x00, 0x70, 0xff}, // red
	{0x70, 0x00, 0xe0, 0xff}, // purple
	{0x70, 0xe0, 0x00, 0xff}, // green
}

var outtileimg string
var palimages string

func init() {
	flag.StringVar(&outtileimg, "tiles", "", "Tile set output image")
	flag.StringVar(&palimages, "palimgs", "", "Dump images in each palette colour")
}

func main() {
	flag.Parse()

	var palettes []color.Palette

	// take first arg as the png file to open
	filename := flag.Arg(0)

	gridx, err := strconv.Atoi(flag.Arg(1))
	if err != nil {
		panic(err)
	}

	gridy, err := strconv.Atoi(flag.Arg(2))
	if err != nil {
		panic(err)
	}

	// open it
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// decode it
	rawimg, err := png.Decode(file)
	if err != nil {
		panic(err)
	}

	// make sure it's a indexed/paletted image
	img, ok := rawimg.(*image.Paletted)
	if !ok {
		panic("Expected Paletted image!")
	}

	bpp := 0
	switch len(img.Palette) {
	case 2:
		bpp = 1
	case 4:
		bpp = 2
	case 16:
		bpp = 4
	case 256:
		bpp = 8
	default:
		panic("Expected 2, 4, 8, 16 or 256 colors!")
	}

	// split image into tiles, deduplicate them and make a tilemap
	var tiles []*image.Paletted
	var tilemap [][]uint
	for y := 0; y < img.Rect.Dy(); y += tileSize {
		tilemap = append(tilemap, make([]uint, img.Rect.Dx()/tileSize))
		ty := y / 8
		for x := 0; x < img.Rect.Dx(); x += tileSize {
			tx := x / 8
			index := len(tiles)
			tile := img.SubImage(image.Rect(x, y, x+tileSize, y+tileSize)).(*image.Paletted)
			found := false
			for i, orig := range tiles {
				if sameTile(orig, tile) {
					index = i
					found = true
				}
			}
			if !found {
				tiles = append(tiles, tile)
			}
			tilemap[ty][tx] = uint(index)
		}
	}

	// split tilemap into a series of grid cells
	// useful for sprites / text characters
	var gridmap [][]uint
	for y := 0; y < len(tilemap); y += gridy {
		for x := 0; x < len(tilemap[y]); x += gridx {
			cell := make([]uint, gridy*gridx)
			for gy := 0; gy < gridy; gy++ {
				for gx := 0; gx < gridx; gx++ {
					cell[gy*gridx+gx] = tilemap[y+gy][x+gx]
				}
			}
			gridmap = append(gridmap, cell)
		}
	}

	// dump a colour palette. This assumes the original
	// image is greyscale and you want to colour it with
	// multiple different colours. It takes the colours
	// in the input image's palette and generates a new
	// one with the same luminosity for each colour in
	// palette at the top of this file.
	fmt.Println("# Colors:", len(palette)+2, "BPP:", bpp)
	// dump an all black / transparent first palette
	for range img.Palette {
		fmt.Printf("%s ", hex([]uint{0, 0, 0, 0}))
	}
	fmt.Printf("\n")
	for _, tcol := range palette {
		tc, _ := colorful.MakeColor(tcol)
		var pal color.Palette
		for _, col := range img.Palette {
			// get the luminosity of the color in the image palette
			ic, _ := colorful.MakeColor(col)
			il, _, _ := ic.Lab()

			// get the luminosity of black
			black := colorful.Color{0, 0, 0}
			bl, _, _ := black.Lab()

			// get the luminosity of white
			white := colorful.Color{1, 1, 1}
			wl, _, _ := white.Lab()

			// figure out the fractional brightness of the
			// img.Palette color compared to black and white
			fract := (il - bl) / wl

			// get the blend of black and the palette colour
			rc := black.BlendLab(tc, fract)

			// convert this to clamped RGB values
			rr, gg, bb := rc.Clamped().RGB255()

			// set black to transparent, all else opaque
			aa := 0xff
			if rr == 0 && gg == 0 && bb == 0 {
				aa = 0
			}

			pal = append(pal, color.RGBA{uint8(rr) & 0xf0, uint8(gg) & 0xf0, uint8(bb) & 0xf0, uint8(aa) & 0xf0})

			fmt.Printf("%s ", hex([]uint{uint(bb) >> 4, uint(gg) >> 4, uint(rr) >> 4, uint(aa) >> 4}))
		}
		palettes = append(palettes, pal)
		fmt.Printf("\n")
	}
	// dump the original image palette last
	for _, col := range img.Palette {
		r, g, b, a := col.RGBA()
		if r == 0 && g == 0 && b == 0 {
			a = 0
		}
		fmt.Printf("%s ", hex([]uint{uint(b) >> 4, uint(g) >> 4, uint(r) >> 4, uint(a) >> 4}))
	}
	fmt.Printf("\n\n")

	tileimgW := len(tiles)
	if tileimgW > 16 {
		tileimgW = 16
	}
	tileimgH := len(tiles) / tileimgW
	if (tileimgW * tileimgH) < len(tiles) {
		tileimgH++
	}
	tileimg := image.NewPaletted(image.Rect(0, 0, tileimgW*tileSize, tileimgH*tileSize), img.Palette)
	tileimgX := 0
	tileimgY := 0

	// dump hex of the tiles
	fmt.Println("# Tiles:", len(tiles),
		"size:", len(tiles)*tileSize*tileSize*bpp/intSize)
	values := make([]uint, intSize/bpp)
	for _, tile := range tiles {
		i := 0
		for y := tile.Rect.Min.Y; y < tile.Rect.Max.Y; y++ {
			for x := tile.Rect.Min.X; x < tile.Rect.Max.X; x++ {
				values[i] = uint(tile.ColorIndexAt(x, y))
				i++
				if i >= len(values) {
					fmt.Printf("%s ", hex(values))
					i = 0
				}
			}
		}
		draw.Draw(tileimg,
			image.Rect(
				tileimgX*tileSize,
				tileimgY*tileSize,
				(tileimgX+1)*tileSize,
				(tileimgY+1)*tileSize),
			tile, tile.Rect.Min, draw.Over)
		tileimgX++
		if tileimgX > 16 {
			tileimgX = 0
			tileimgY++
		}
		fmt.Printf("\n")
	}

	if outtileimg != "" {
		fmt.Println("writing tiles to", outtileimg)
		outfp, err := os.Create(outtileimg)
		defer outfp.Close()
		if err != nil {
			panic(err)
		}
		png.Encode(outfp, tileimg)
	}

	fmt.Println("")

	// figure out if tile ids need to be 8 or 16 bits
	bpt := 8
	if len(tiles) > 255 {
		bpt = 16
	}

	// dump the tile map as hex
	fmt.Println("# Tilemap:",
		len(tilemap), "x", len(tilemap[0]),
		"size:", len(tilemap)*len(tilemap[0]),
		"bpt:", bpt)
	for y := 0; y < len(tilemap); y++ {
		for x := 0; x < len(tilemap[y]); x++ {
			format := fmt.Sprintf("%%0%dX ", bpt/4)
			if bpt == 8 {
				fmt.Printf(format, tilemap[y][x])
			}
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	// figure out how to align the grid to a
	// power of 2 (2, 4, 8, 16, 32...)
	// This is so that addresses of grid cells
	// are a whole number of bits.
	align := 1
	for align < len(gridmap[0]) {
		align <<= 1
	}

	// Dump the "gridmap" as hex. This could
	// be a sprite sheet, or text characters
	// or whatever.
	fmt.Println("# Gridmap:",
		len(gridmap),
		"size:", len(gridmap)*align,
		"bpt:", bpt)
	for y := 0; y < len(gridmap); y++ {
		x := 0
		format := fmt.Sprintf("%%0%dX ", bpt/4)
		for ; x < len(gridmap[y]); x++ {
			fmt.Printf(format, gridmap[y][x])
		}
		for ; x < align; x++ {
			fmt.Printf(format, 0)
		}
		fmt.Printf("\n")
	}
	fmt.Printf("\n")

	if palimages != "" {
		for i, pal := range palettes {
			img.Palette = pal
			outfp, err := os.Create(fmt.Sprintf("%s_%d.png", palimages, i))
			defer outfp.Close()
			if err != nil {
				panic(err)
			}
			png.Encode(outfp, img)
		}
	}
}

// compare tiles to see if they are the same
func sameTile(orig *image.Paletted, tile *image.Paletted) bool {
	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			ox := x + orig.Rect.Min.X
			oy := y + orig.Rect.Min.Y
			tx := x + tile.Rect.Min.X
			ty := y + tile.Rect.Min.Y
			if orig.ColorIndexAt(ox, oy) != tile.ColorIndexAt(tx, ty) {
				return false
			}
		}
	}
	return true
}

// convert a series of small ints into a packed
// hex number intSize big
func hex(values []uint) string {
	bits := intSize / len(values)
	mask := uint(0)
	for i := 0; i < bits; i++ {
		mask <<= 1
		mask |= 1
	}
	num := uint(0)
	for i, val := range values {
		num |= (val & mask) << uint(i*bits)
	}
	return fmt.Sprintf("%08X", num)[8-(intSize/4):]
}
