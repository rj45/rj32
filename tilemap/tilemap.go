/*
 * Copyright (c) 2020-2021 rj45 and contributors
 * Licensed under the MIT license.
 * https://github.com/rj45/rj32
 */

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sort"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/muesli/clusters"
	"github.com/muesli/kmeans"
	"github.com/shabbyrobe/wu2quant"
)

// size of one side of square tiles (ie 8x8)
const tileSize = 8

// number of bits in integers by default
const intSize = 16

var ditherMatrix = map[string][][]float64{
	"floydsteinberg":    {{0, 0, 7.0 / 16.0}, {3.0 / 16.0, 5.0 / 16.0, 1.0 / 16.0}},
	"jarvisjudiceninke": {{0, 0, 0, 7.0 / 48.0, 5.0 / 48.0}, {3.0 / 48.0, 5.0 / 48.0, 7.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0}, {1.0 / 48.0, 3.0 / 48.0, 5.0 / 48.0, 3.0 / 48.0, 1.0 / 48.0}},
	"stucki":            {{0, 0, 0, 8.0 / 42.0, 4.0 / 42.0}, {2.0 / 42.0, 4.0 / 42.0, 8.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0}, {1.0 / 42.0, 2.0 / 42.0, 4.0 / 42.0, 2.0 / 42.0, 1.0 / 42.0}},
	"atkinson":          {{0, 0, 1.0 / 8.0, 1.0 / 8.0}, {1.0 / 8.0, 1.0 / 8.0, 1.0 / 8.0, 0}, {0, 1.0 / 8.0, 0, 0}},
	"burkes":            {{0, 0, 0, 8.0 / 32.0, 4.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 8.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}},
	"sierra":            {{0, 0, 0, 5.0 / 32.0, 3.0 / 32.0}, {2.0 / 32.0, 4.0 / 32.0, 5.0 / 32.0, 4.0 / 32.0, 2.0 / 32.0}, {0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 0}},
	"tworowsierra":      {{0, 0, 0, 4.0 / 16.0, 3.0 / 16.0}, {1.0 / 32.0, 2.0 / 32.0, 3.0 / 32.0, 2.0 / 32.0, 1.0 / 32.0}},
	"sierralite":        {{0, 0, 2.0 / 4.0}, {1.0 / 4.0, 1.0 / 4.0, 0}},
}

var colorConvFunc = map[string]func(color.RGBA) color.RGBA{
	"24": colorConv24bpp,
	"12": colorConv12bpp,
	"8":  colorConv8bpp,
}

var srcfile string
var outtileimg string
var outtilehex string
var outmaphex string
var outpalhex string
var idOffset int
var gridw int
var gridh int
var samethresh float64
var palettes int
var ditherer string
var outfile string
var clusterfile string
var colorsPerPal int
var recluster bool
var gentest string
var outputjson string
var tilebpp int
var colorConv string

func init() {
	flag.StringVar(&srcfile, "in", "",
		"Input png file name (required)")

	flag.StringVar(&outtileimg, "tilesimg", "",
		"Tile set output image")

	flag.StringVar(&outtilehex, "tiles", "",
		"Tile set output hex file")

	flag.StringVar(&outpalhex, "pal", "",
		"Palette set output hex file")

	flag.StringVar(&outmaphex, "map", "",
		"Tile map output hex file")

	flag.StringVar(&outputjson, "json", "",
		"Output JSON file for your own packing format")

	// used when arranging tiles as text characters or sprites
	flag.IntVar(&gridw, "gridw", 1, "x size of grid")
	flag.IntVar(&gridh, "gridh", 1, "y size of grid")

	flag.Float64Var(&samethresh, "samethresh", 0.001,
		"Threshold distance to consider colors the same")

	flag.IntVar(&palettes, "palettes", 16,
		"max palettes allowed")

	flag.IntVar(&colorsPerPal, "perpalette", 16,
		"colors per palette")

	flag.IntVar(&tilebpp, "tilebpp", 4,
		"bits per tile in tile bitmaps")

	flag.StringVar(&colorConv, "colorconv", "12",
		"Color conversion routine (12: 4:4:4 RGB, 8: 3:3:2 RGB")

	flag.StringVar(&ditherer, "dither", "floydsteinberg",
		"which dither matrix to use")

	flag.StringVar(&outfile, "outfile", "",
		"final image output file")

	flag.StringVar(&clusterfile, "clusterfile", "",
		"clustered image output file")

	flag.BoolVar(&recluster, "recluster", false,
		"dither and recluster dithered image")

	flag.StringVar(&gentest, "gentest", "",
		"Generate a test image")
}

func main() {
	flag.Parse()

	var rawimg SubableImage
	if gentest != "" {
		rawimg = genTestImage()

		outfp, err := os.Create(gentest)
		if err != nil {
			panic(err)
		}
		png.Encode(outfp, rawimg)
		outfp.Close()
	} else if srcfile != "" {
		rawimg = readImage(srcfile).(SubableImage)
	} else {
		flag.Usage()
		fmt.Println("Missing -in file; nothing to do!")
		os.Exit(1)
	}

	// Problem:
	// - with N palettes of 15 or 16 colors, how do you generate
	//   the optimal set of palettes so as to maximize the
	//   overall total colors and quality of an encoded photo
	// Algorithm:
	// - split image into tiles and "grids" (groups of tiles)
	// - for each "grid"
	//   - pull out an array of colors into Luv color space
	//   - sort colors by hue (removing pixel location relevance)
	//   - flatten colors into one long array
	// - run k-means clustering to determine which grids have
	//   common colors
	// - take each cluster and build new images with the tiles
	// - quantize the colors in each cluster tile image to colorsPerPal colors
	// - build a full palette with all the colors
	// - do a global (whole image) dither with this full palette
	// - re-cluster the grids with the dithered image
	//   - because dithering moves pixels into other tiles, this
	//     helps reduce dithering glitches
	// - re-quantize the colors for each cluster based on the
	//   original undithered tiles
	// - re-dither the original image using only palettes for
	//   each cluster instead of a global dither over full palette
	// - emit tile images, map and palette

	// TODO: tile reduction:
	// - try doing k-means on the tile patterns somehow

	quant := wu2quant.New()
	var img *image.RGBA

	clustimg, gimgs, clustergrids := clusterImage(rawimg)

	if clusterfile != "" {
		fmt.Println("writing cluster result to", clusterfile)
		outfp, err := os.Create(clusterfile)
		if err != nil {
			panic(err)
		}
		png.Encode(outfp, clustimg)
		outfp.Close()
	}

	clustpals := clusterPalettes(quant, gimgs, clustergrids)
	fullpal, fullpalrgb, palsets := genFullPalette(clustpals)

	if recluster && ditherer != "none" {
		fullset := make(intset)
		for i := range fullpal {
			fullset.add(i)
		}
		img, _ := dither(rawimg, fullpalrgb, func(w, x, y int) intset {
			return fullset
		})

		// re-cluster based on dithered image to reduce dither artifacts
		_, _, clustergrids = clusterImage(img)

		// determine cluster palettes based on the original images
		// not on the dithered image
		clustpals = clusterPalettes(quant, gimgs, clustergrids)
		fullpal, fullpalrgb, palsets = genFullPalette(clustpals)
	}

	gridclusters := genGridClusters(clustergrids)
	findPal := gridPal(palsets, gridclusters)
	img, pimg := dither(rawimg, fullpalrgb, findPal)

	if outfile != "" {
		fmt.Println("writing intermediate result to", outfile)
		outfp, err := os.Create(outfile)
		if err != nil {
			panic(err)
		}
		if pimg != nil {
			png.Encode(outfp, pimg)
		} else {
			png.Encode(outfp, img)
		}
		outfp.Close()
	}

	fmt.Println("full palette size", len(fullpal))
	fmt.Println("num palettes", len(clustpals))
	totalcols := 0
	for _, pal := range clustpals {
		totalcols += len(pal)
	}
	fmt.Println("total colours", totalcols)

	tileimgs, tilemap, _, _ := splitTiles(img)
	clustpals, tilepals := genTilePals(fullpalrgb, palsets, gridclusters, tilemap)
	dumpJSON(clustpals, tilepals, tileimgs, tilemap)
	dumpPalettes(clustpals)
	dumpTiles(clustpals, tilepals, tileimgs)
	dumpMap(tilemap, tilepals)
}

func genTestImage() *image.RGBA {
	width := 640
	height := 400
	gw := tileSize * gridw
	gh := tileSize * gridh

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	draw.Draw(img, img.Bounds(), image.Black, image.Point{}, draw.Src)

	for yo := 0; yo < height; yo += gh {
		for xo := 0; xo < width; xo += gw {

			rand.Intn(9)

			for y := 0; y < gh; y++ {
				setTestPatternPixel(img, xo+0, yo+y)
				setTestPatternPixel(img, xo+gw-1, yo+y)
			}
			for x := 0; x < gw; x++ {
				setTestPatternPixel(img, xo+x, yo+0)
				setTestPatternPixel(img, xo+x, yo+gh-1)
			}

			mx := (gw / 2) - 1
			my := (gh / 2) - 1

			m := ((xo >> 3) + (yo >> 3)) & 0x7

			if (m & 4) == 4 {
				setTestPatternPixel(img, xo+mx+0, yo+my+0)
				setTestPatternPixel(img, xo+mx+1, yo+my+0)
				setTestPatternPixel(img, xo+mx+0, yo+my+1)
				setTestPatternPixel(img, xo+mx+1, yo+my+1)
			}

			m &= 0x3

			setTestPatternPixel(img, xo+mx-1, yo+my-1)
			if m > 0 {
				setTestPatternPixel(img, xo+mx-1, yo+my+2)
			}
			if m > 1 {
				setTestPatternPixel(img, xo+mx+2, yo+my+2)
			}
			if m > 2 {
				setTestPatternPixel(img, xo+mx+2, yo+my-1)
			}
		}
	}

	return img
}

func setTestPatternPixel(img *image.RGBA, x, y int) {
	i := uint8(x) & 0x7
	j := uint8(y) & 0x7

	m := (x>>3)&3 | ((y>>3)&1)<<2
	m += y >> 3
	m &= 0x7

	switch m {
	case 0:
		img.SetRGBA(x, y, color.RGBA{
			R: (i + j) << 4,
			G: (j + 8) << 4,
			B: (i + 8) << 4,
			A: 255,
		})
	case 1:
		img.SetRGBA(x, y, color.RGBA{
			R: (i + 8) << 4,
			G: (i + j) << 4,
			B: (j + 8) << 4,
			A: 255,
		})
	case 2:
		img.SetRGBA(x, y, color.RGBA{
			R: (i + 8) << 4,
			G: (j + 8) << 4,
			B: (i + j) << 4,
			A: 255,
		})
	case 3:
		img.SetRGBA(x, y, color.RGBA{
			R: (i + 8) << 4,
			G: (j + 8) << 4,
			B: (i + 8) << 4,
			A: 255,
		})
	case 4:
		img.SetRGBA(x, y, color.RGBA{
			R: (i + j) << 4,
			G: (i + 8) << 4,
			B: (j + 8) << 4,
			A: 255,
		})
	case 5:
		img.SetRGBA(x, y, color.RGBA{
			R: (j + 8) << 4,
			G: (i + j) << 4,
			B: (i + 8) << 4,
			A: 255,
		})
	case 6:
		img.SetRGBA(x, y, color.RGBA{
			R: (j + 8) << 4,
			G: (i + 8) << 4,
			B: (i + j) << 4,
			A: 255,
		})
	case 7:
		img.SetRGBA(x, y, color.RGBA{
			R: (j + 8) << 4,
			G: (i + 8) << 4,
			B: (j + 8) << 4,
			A: 255,
		})
	}
}

func genTilePals(fullpal color.Palette, clustpals []intset, gridclusters []int, tilemap [][]uint) ([]color.Palette, []int) {
	var pals []color.Palette

	for _, set := range clustpals {
		var pal color.Palette
		for i := range set {
			pal = append(pal, fullpal[i])
		}
		pals = append(pals, pal)
	}

	var tilepals []int
	for y := 0; y < len(tilemap); y++ {
		for x := 0; x < len(tilemap[y]); x++ {
			gx := x / gridw
			gy := y / gridh
			gw := len(tilemap[y]) / gridw
			gid := (gy * gw) + gx

			palid := gridclusters[gid]

			tid := tilemap[y][x]

			for int(tid) >= len(tilepals) {
				tilepals = append(tilepals, 0)
			}

			tilepals[tid] = palid
		}
	}

	return pals, tilepals
}

func dumpPalettes(pals []color.Palette) {
	hexfile := os.Stdout
	if outpalhex == "" {
		fmt.Println("skipping palette output")
		return
	}
	if outpalhex != "-" {
		fmt.Println("writing palettes hex to", outpalhex)
		var err error
		hexfile, err = os.Create(outpalhex)
		if err != nil {
			panic(err)
		}
		defer hexfile.Close()
	}

	fmt.Fprintln(hexfile, "v2.0 raw")

	for _, pal := range pals {
		for _, col := range pal {
			rgba := col.(color.RGBA)
			fmt.Fprintf(hexfile, "%s ", hex([]uint{
				uint(rgba.B) >> 4,
				uint(rgba.G) >> 4,
				uint(rgba.R) >> 4,
				uint(rgba.A) >> 4,
			}))
		}
		fmt.Fprintf(hexfile, "\n")
	}
}

func findPalColor(pal color.Palette, col color.Color) int {
	r1, g1, b1, _ := col.RGBA()
	for index := range pal {
		other := pal[index]

		r2, g2, b2, _ := other.RGBA()

		if r1 == r2 && g1 == g2 && b1 == b2 {
			return index
		}
	}

	panic(fmt.Errorf("color should be found %v %v", col, pal))
}

func dumpTiles(pals []color.Palette, tilepals []int, tiles []image.Image) {
	hexfile := os.Stdout
	if outtilehex == "" {
		fmt.Println("Skipping tiles output")
		return
	}

	if outtilehex != "-" {
		fmt.Println("writing tiles hex to", outtilehex)
		var err error
		hexfile, err = os.Create(outtilehex)
		if err != nil {
			panic(err)
		}
		defer hexfile.Close()
	}

	fmt.Fprintln(hexfile, "v2.0 raw")

	tileimgW := len(tiles)
	if tileimgW > 16 {
		tileimgW = 16
	}
	tileimgH := len(tiles) / tileimgW
	if (tileimgW * tileimgH) < len(tiles) {
		tileimgH++
	}
	tileimg := image.NewRGBA(image.Rect(0, 0, tileimgW*tileSize, tileimgH*tileSize))
	tileimgX := 0
	tileimgY := 0

	// dump hex of the tiles
	fmt.Println("# Tiles:", len(tiles),
		"size:", len(tiles)*tileSize*tileSize*tilebpp/intSize)
	values := make([]uint, intSize/tilebpp)
	for tid, tile := range tiles {
		i := 0
		tilesz := tile.Bounds()
		for y := tilesz.Min.Y; y < tilesz.Max.Y; y++ {
			for x := tilesz.Min.X; x < tilesz.Max.X; x++ {
				col := colorAt(tile, x, y)
				values[i] = uint(findPalColor(pals[tilepals[tid]], col))
				i++
				if i >= len(values) {
					fmt.Fprintf(hexfile, "%s ", hex(values))
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
			tile, tilesz.Min, draw.Src)
		tileimgX++
		if tileimgX > 16 {
			tileimgX = 0
			tileimgY++
		}
		fmt.Fprintf(hexfile, "\n")
	}

	if outtileimg != "" {
		fmt.Println("writing tiles to", outtileimg)
		outfp, err := os.Create(outtileimg)
		if err != nil {
			panic(err)
		}
		png.Encode(outfp, tileimg)
		outfp.Close()
	}
}

func dumpMap(tilemap [][]uint, tilepals []int) {
	hexfile := os.Stdout
	if outmaphex == "" {
		fmt.Println("Skipping tilemap output")
		return
	}
	if outmaphex != "-" {
		fmt.Println("writing tilemap hex to", outmaphex)
		var err error
		hexfile, err = os.Create(outmaphex)
		if err != nil {
			panic(err)
		}
		defer hexfile.Close()
	}

	fmt.Fprintln(hexfile, "v2.0 raw")

	// dump the tile map as hex
	block := nextPowerOfTwo(len(tilemap)*len(tilemap[0])) / 2
	fmt.Println("# Tilemap:",
		len(tilemap[0]), "x", len(tilemap),
		"size:", len(tilemap)*len(tilemap[0]),
		"block size:", block,
	)
	i := 0
	for y := 0; y < len(tilemap); y++ {
		for x := 0; x < len(tilemap[y]); x += 2 {
			id1 := tilemap[y][x]
			id2 := tilemap[y][x+1]
			id1 += uint(idOffset)
			id2 += uint(idOffset)

			val := (id1 & 0xff) | ((id2 & 0xff) << 8)

			fmt.Fprintf(hexfile, "%04X ", val)
			i++
			if (i % 16) == 0 {
				fmt.Fprintf(hexfile, "\n")
			}
		}
	}

	for i < block {
		fmt.Fprintf(hexfile, "%04X ", 0)
		i++
		if (i % 16) == 0 {
			fmt.Fprintf(hexfile, "\n")
		}
	}
	fmt.Fprintf(hexfile, "\n\n")

	for y := 0; y < len(tilemap); y++ {
		for x := 0; x < len(tilemap[y]); x += 2 {
			id1 := tilemap[y][x]
			id2 := tilemap[y][x+1]
			pal1 := tilepals[id1]
			pal2 := tilepals[id2]
			id1 += uint(idOffset)
			id2 += uint(idOffset)

			val := ((uint(pal1) & 0xf) << 4) | ((id1 >> 8) & 0xf) |
				((uint(pal2) & 0xf) << 12) | (((id2 >> 8) & 0xf) << 8)

			fmt.Fprintf(hexfile, "%04X ", val)
			i++
			if (i % 16) == 0 {
				fmt.Fprintf(hexfile, "\n")
			}
		}
	}
	for i < (block * 2) {
		fmt.Fprintf(hexfile, "%04X ", 0)
		i++
		if (i % 16) == 0 {
			fmt.Fprintf(hexfile, "\n")
		}
	}
}

func dumpJSON(clustpals []color.Palette, tilepals []int, tileimgs []image.Image, tilemap [][]uint) {
	if outputjson == "" {
		fmt.Println("Skipping JSON output")
		return
	}

	data := struct {
		Palettes []color.Palette `json:"palettes"`
		TilePals []int           `json:"tilepals"`
		Tiles    [][]int         `json:"tiles"`
		TileMap  [][]uint        `json:"tilemap"`
	}{
		Palettes: clustpals,
		TilePals: tilepals,
		TileMap:  tilemap,
	}
	for tid, tile := range tileimgs {
		i := 0
		tilesz := tile.Bounds()
		values := make([]int, tilesz.Dx()*tilesz.Dy())
		for y := tilesz.Min.Y; y < tilesz.Max.Y; y++ {
			for x := tilesz.Min.X; x < tilesz.Max.X; x++ {
				col := colorAt(tile, x, y)
				values[i] = findPalColor(clustpals[tilepals[tid]], col)
				i++
			}
		}
		data.Tiles = append(data.Tiles, values)
	}

	buf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		panic(err)
	}

	fmt.Println("writing json to", outputjson)
	os.WriteFile(outputjson, buf, 0666)
}

func nextPowerOfTwo(v int) int {
	power := 1
	for power < v {
		power <<= 1
	}
	return power
}

func genFullPalette(gridpals []color.Palette) ([]colorful.Color, color.Palette, []intset) {
	fmt.Println("Generating full palette...")
	var fullpal []colorful.Color
	var fullpalrgb color.Palette
	var gridpali []intset
	for _, pal := range gridpals {
		ipal := make(intset)
		for _, col := range pal {
			ccol, _ := colorful.MakeColor(col)
			index := -1
			for fi, fcol := range fullpal {
				if fcol.DistanceLuv(ccol) < samethresh {
					index = fi
					break
				}
			}
			if index < 0 {
				index = len(fullpal)
				fullpal = append(fullpal, ccol)
				fullpalrgb = append(fullpalrgb, col)
			}
			ipal.add(index)
		}
		gridpali = append(gridpali, ipal)
	}

	return fullpal, fullpalrgb, gridpali
}

func genGridClusters(clustergrids [][]int) []int {
	var ret []int

	for clust, clustgrid := range clustergrids {
		for _, grid := range clustgrid {
			for len(ret) <= grid {
				ret = append(ret, -1)
			}
			ret[grid] = clust
		}
	}

	return ret
}

// generate new images with the grids in each cluster, then
// quantize a palette for each
func clusterPalettes(quant *wu2quant.Quantizer, gridimgs []image.Image, clustergrids [][]int) []color.Palette {
	var pals []color.Palette
	for _, grids := range clustergrids {
		img := image.NewRGBA(image.Rect(0, 0, len(grids)*tileSize*gridw, tileSize*gridh))
		for i, id := range grids {
			area := image.Rect(
				i*tileSize*gridw, 0,
				(i+1)*tileSize*gridw, tileSize*gridh,
			)
			draw.Draw(img,
				area,
				gridimgs[id], gridimgs[id].Bounds().Min, draw.Src)
		}

		pal := quantizeColors(quant, colorsPerPal, img)
		pals = append(pals, pal)
	}
	return pals
}

func clusterImage(img SubableImage) (image.Image, []image.Image, [][]int) {
	var gridimgs []image.Image
	var gridobs clusters.Observations

	fmt.Println("Prepping tiles for clustering...")
	sz := img.Bounds()
	for y := 0; y < sz.Dy(); y += (gridh * tileSize) {
		for x := 0; x < sz.Dx(); x += (gridw * tileSize) {
			gridimg := img.SubImage(image.Rect(
				x, y,
				x+(tileSize*gridw),
				y+(tileSize*gridh),
			).Intersect(sz))

			gridimgs = append(gridimgs, gridimg)

			var colors []lab
			gsz := gridimg.Bounds()
			for j := gsz.Min.Y; j < gsz.Max.Y; j++ {
				for i := gsz.Min.X; i < gsz.Max.X; i++ {
					col := colorAt(gridimg, x, y)
					colors = append(colors, toLuv(col))
				}
			}

			// sort colors by hue
			// this is to remove the spatial element from the clustering
			sort.Slice(colors, func(a, b int) bool {
				_, _, ah := colorful.LuvToLuvLCh(colors[a][0], colors[a][1], colors[a][2])
				_, _, bh := colorful.LuvToLuvLCh(colors[b][0], colors[b][1], colors[b][2])
				return ah < bh
			})

			var coords clusters.Coordinates
			for _, col := range colors {
				for _, v := range col {
					coords = append(coords, v)
				}
			}

			gridobs = append(gridobs, coords)
		}
	}

	fmt.Println("Clustering", len(gridobs), "tiles...")
	clusterer := kmeans.New()
	clusters, err := clusterer.Partition(gridobs, palettes)
	if err != nil {
		panic(err)
	}

	clusteredimg := image.NewRGBA(sz)
	var clustergrids [][]int
	for _, cluster := range clusters {
		// figure out the representitive color
		var replab lab
		for i, v := range cluster.Center {
			replab[i%3] += v
		}
		replab[0] /= float64(len(cluster.Center) / 3)
		replab[1] /= float64(len(cluster.Center) / 3)
		replab[2] /= float64(len(cluster.Center) / 3)
		repcol := colorful.Luv(replab[0], replab[1], replab[2])

		var grids []int
		for _, clustobs := range cluster.Observations {
			index := -1
			for i, gridob := range gridobs {
				a := clustobs.Coordinates()
				b := gridob.Coordinates()
				// need to compare pointers to get a good identification
				// of the original grid
				if &a[0] == &b[0] {
					index = i
					break
				}
			}
			if index < 0 {
				panic("could not find grid!")
			}
			grids = append(grids, index)

			gridimg := gridimgs[index].(draw.Image)
			gsz := gridimg.Bounds()
			for y := gsz.Min.Y; y < gsz.Max.Y; y++ {
				for x := gsz.Min.X; x < gsz.Max.X; x++ {
					clusteredimg.Set(x, y, repcol)
				}
			}
		}
		clustergrids = append(clustergrids, grids)
	}

	return clusteredimg, gridimgs, clustergrids
}

func quantizeColors(quant *wu2quant.Quantizer, num int, img image.Image) color.Palette {
	// find the optimal num colors
	rawpal := quant.Quantize(make(color.Palette, 0, num), img)

	// squash those colors to 12 bpp
	// this can produce duplicates so we generate a new palette
	var pal color.Palette
	for _, rcol := range rawpal {
		col := colorConvFunc[colorConv](rcol.(color.RGBA))
		found := false
		for _, pcol := range pal {
			if pcol == col {
				found = true
				break
			}
		}
		if !found {
			pal = append(pal, col)
		}
	}

	return pal
}

type SubableImage interface {
	image.Image
	SubImage(r image.Rectangle) image.Image
}

func readImage(filename string) image.Image {
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
	rawimg = imgConv12bpp(rawimg)

	return rawimg
}

func splitTiles(rawimg image.Image) ([]image.Image, [][]uint, [][]uint, []image.Image) {
	img := rawimg.(SubableImage)
	imgrect := img.Bounds()

	// split image into tiles, deduplicate them and make a tilemap
	var tiles []image.Image
	var tilemap [][]uint
	for y := 0; y < imgrect.Dy(); y += tileSize {
		tilemap = append(tilemap, make([]uint, imgrect.Dx()/tileSize))
		ty := y / tileSize
		for x := 0; x < imgrect.Dx(); x += tileSize {
			tx := x / tileSize

			index := len(tiles)
			tile := img.SubImage(image.Rect(x, y, x+tileSize, y+tileSize))
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
	var gridimgs []image.Image
	for y := 0; y < len(tilemap); y += gridh {
		for x := 0; x < len(tilemap[y]); x += gridw {
			var cell []uint
			for gy := 0; gy < gridh; gy++ {
				for gx := 0; gx < gridw; gx++ {
					if y+gy < len(tilemap) && x+gx < len(tilemap[y+gy]) {
						cell = append(cell, tilemap[y+gy][x+gx])
					}
				}
			}
			gridmap = append(gridmap, cell)
			ty := y * tileSize
			tx := x * tileSize
			gridimg := img.SubImage(image.Rect(
				tx, ty,
				tx+(tileSize*gridw),
				ty+(tileSize*gridh),
			))
			gridimgs = append(gridimgs, gridimg)
		}
	}

	return tiles, tilemap, gridmap, gridimgs
}

// imgConv12bpp converts the image to 4 bits per r, g, b
func imgConv12bpp(src image.Image) image.Image {
	imgrect := src.Bounds()
	dest := image.NewRGBA(imgrect)
	for y := 0; y < imgrect.Dy(); y++ {
		for x := 0; x < imgrect.Dx(); x++ {
			col := colorAt(src, x, y)

			col = colorConvFunc[colorConv](col)

			dest.SetRGBA(x, y, col)
		}
	}

	return dest
}

// colorConv24bpp does nothing to the color
func colorConv24bpp(col color.RGBA) color.RGBA {
	return col
}

// colorConv12bpp converts an RGB color to 12bpp
// the color is rounded to the nearest 0x10th for
// each of r, g, b, a
func colorConv12bpp(col color.RGBA) color.RGBA {
	var n int

	n = int(col.R) + 7
	if n > 255 {
		n = 255
	}
	col.R = uint8(n & 0xf0)

	n = int(col.G) + 7
	if n > 255 {
		n = 255
	}
	col.G = uint8(n & 0xf0)

	n = int(col.B) + 7
	if n > 255 {
		n = 255
	}
	col.B = uint8(n & 0xf0)

	n = int(col.A) + 7
	if n > 255 {
		n = 255
	}
	col.A = uint8(n & 0xf0)

	return col
}

// 3:3:2 RGB
func colorConv8bpp(col color.RGBA) color.RGBA {
	var n int

	n = int(col.R) + 15
	if n > 255 {
		n = 255
	}
	col.R = uint8(n & 0xe0)

	n = int(col.G) + 15
	if n > 255 {
		n = 255
	}
	col.G = uint8(n & 0xe0)

	n = int(col.B) + 31
	if n > 255 {
		n = 255
	}
	col.B = uint8(n & 0xc0)

	return col
}

// compare tiles to see if they are the same
func sameTile(orig image.Image, tile image.Image) bool {
	ob := orig.Bounds()
	tb := tile.Bounds()
	for y := 0; y < tileSize; y++ {
		for x := 0; x < tileSize; x++ {
			ox := x + ob.Min.X
			oy := y + ob.Min.Y
			tx := x + tb.Min.X
			ty := y + tb.Min.Y
			r1, g1, b1, a1 := orig.At(ox, oy).RGBA()
			r2, g2, b2, a2 := tile.At(tx, ty).RGBA()
			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
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

// from https://github.com/ericpauley/go-quantize/blob/master/quantize/mediancut.go
func colorAt(m image.Image, x int, y int) color.RGBA {
	switch i := m.(type) {
	case *image.YCbCr:
		yi := i.YOffset(x, y)
		ci := i.COffset(x, y)
		c := color.YCbCr{
			i.Y[yi],
			i.Cb[ci],
			i.Cr[ci],
		}
		return color.RGBA{c.Y, c.Cb, c.Cr, 255}
	case *image.RGBA:
		ci := i.PixOffset(x, y)
		return color.RGBA{i.Pix[ci+0], i.Pix[ci+1], i.Pix[ci+2], i.Pix[ci+3]}
	default:
		return color.RGBAModel.Convert(i.At(x, y)).(color.RGBA)
	}
}

type intset map[int]struct{}

func (set intset) add(i int) {
	set[i] = struct{}{}
}

type lab [3]float64

type labimg [][]lab

func toLuv(col color.Color) lab {
	ret, ok := colorful.MakeColor(col)
	if !ok {
		panic("Bad alpha")
	}
	l, a, b := ret.Luv()
	return lab{l, a, b}
}

func initLabImg(img image.Image) labimg {
	sz := img.Bounds()
	ret := make(labimg, sz.Dy())
	for y := 0; y < sz.Dy(); y++ {
		ret[y] = make([]lab, sz.Dx())
	}
	return ret
}

type findPal func(w, x, y int) intset

func gridPal(pals []intset, gridpals []int) findPal {
	return func(w, x, y int) intset {
		gy := y / tileSize / gridh
		goff := (w / tileSize / gridw) * gy
		gx := x / tileSize / gridw
		return pals[gridpals[goff+gx]]
	}
}

func dither(img image.Image, fullpal color.Palette, findPal findPal) (*image.RGBA, *image.Paletted) {
	errimg := initLabImg(img)
	mat := ditherMatrix[ditherer]
	shift := findShift(mat)

	sz := img.Bounds()
	result := image.NewRGBA(image.Rect(0, 0, sz.Dx(), sz.Dy()))
	var presult *image.Paletted

	if len(fullpal) < 256 {
		var pal color.Palette
		for _, col := range fullpal {
			pal = append(pal, col)
		}
		presult = image.NewPaletted(image.Rect(0, 0, sz.Dx(), sz.Dy()), pal)
	}

	fulllab := make([]lab, len(fullpal))
	for i, col := range fullpal {
		fulllab[i] = toLuv(col)
	}

	for y := 0; y < sz.Dy(); y++ {
		for x := 0; x < sz.Dx(); x++ {
			pal := findPal(sz.Dx(), x, y)
			col := toLuv(colorAt(img, x+sz.Min.X, y+sz.Min.Y))
			err := errimg[y][x]

			// determine an error corrected color
			col[0] = col[0] + err[0]*0.75
			col[1] = col[1] + err[1]*0.75
			col[2] = col[2] + err[2]*0.75

			index := findColor(col, x, y, fulllab, pal)
			foundcol := fulllab[index]

			result.Set(x, y, fullpal[index])

			if presult != nil {
				presult.SetColorIndex(x, y, uint8(index))
			}

			err = lab{
				col[0] - foundcol[0],
				col[1] - foundcol[1],
				col[2] - foundcol[2],
			}
			errimg[y][x] = err

			// diffusing the error using the diffusion matrix
			for i, v1 := range mat {
				for j, v2 := range v1 {
					ey := y + i
					ex := x + j + shift
					if ey < len(errimg) && ex < len(errimg[ey]) && ey >= 0 && ex >= 0 {
						errimg[ey][ex][0] += err[0] * v2
						errimg[ey][ex][1] += err[0] * v2
						errimg[ey][ex][2] += err[0] * v2
					}
				}
			}
		}
	}

	return result, presult
}

func findShift(matrix [][]float64) int {
	if matrix == nil {
		return 0
	}

	for _, v1 := range matrix {
		for j, v2 := range v1 {
			if v2 > 0.0 {
				return -j + 1
			}
		}
	}
	return 0
}

func sq(x float64) float64 {
	return x * x
}

func findColor(col lab, x, y int, colors []lab, pal intset) int {
	index := -1
	mindist := math.Inf(1)
	for pi := range pal {
		other := colors[pi]

		dist := math.Sqrt((math.Abs(col[0] - other[0])) +
			sq(math.Abs(col[1]-other[1])) +
			sq(math.Abs(col[2]-other[2])))

		if dist < mindist {
			mindist = dist
			index = pi
		}
	}

	return index
}
