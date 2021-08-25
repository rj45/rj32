package main

import (
	"os"

	"github.com/rj45/rj32/emu/bitfield"
)

type SpriteSheetID struct {
	Sheet  int `bitfield:"4"`
	SheetY int `bitfield:"5"`
	SheetX int `bitfield:"7"`
}

type SpriteY struct {
	Transparent bool `bitfield:"1"`
	FlipY       bool `bitfield:"1"`
	Height      int  `bitfield:"3"`
	Y           int  `bitfield:"11"`
}

type SpriteX struct {
	FlipX bool `bitfield:"1"`
	Width int  `bitfield:"3"`
	X     int  `bitfield:"12"`
}

type Tile struct {
	Palette int `bitfield:"5"`
	TileID  int `bitfield:"11"`
}

type SpriteDims struct {
	FlipY  bool `bitfield:"1"`
	Height int  `bitfield:"7"`
	FlipX  bool `bitfield:"1"`
	Width  int  `bitfield:"7"`
}

type SpriteAddr struct {
	Transparent bool `bitfield:"1"`
	PaletteSet  int  `bitfield:"1"`
	TileSetAddr int  `bitfield:"6"`
	SheetAddr   int  `bitfield:"8"`
}

type SheetPos struct {
	WrapY  bool `bitfield:"1"`
	SheetY int  `bitfield:"7"`
	WrapX  bool `bitfield:"1"`
	SheetX int  `bitfield:"7"`
}

type Pixels struct {
	Color3 int `bitfield:"4"`
	Color2 int `bitfield:"4"`
	Color1 int `bitfield:"4"`
	Color0 int `bitfield:"4"`
}

type Color struct {
	A int `bitfield:"4"`
	R int `bitfield:"4"`
	G int `bitfield:"4"`
	B int `bitfield:"4"`
}

func main() {
	file, err := os.Create("data/fields.go")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	err = bitfield.Gen(file, &SpriteSheetID{}, &bitfield.Config{
		Package: "data",
	})
	if err != nil {
		panic(err)
	}
	file.WriteString("\n")

	err = bitfield.Gen(file, &SpriteY{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	file.WriteString("\n")

	err = bitfield.Gen(file, &SpriteX{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	file.WriteString("\n")

	err = bitfield.Gen(file, &Tile{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	file.WriteString("\n")

	err = bitfield.Gen(file, &Pixels{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	file.WriteString("\n")

	err = bitfield.Gen(file, &Color{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	file.WriteString("\n")

	ofile, err := os.Create("vdp/fields.go")
	if err != nil {
		panic(err)
	}
	defer ofile.Close()

	err = bitfield.Gen(ofile, &SpriteDims{}, &bitfield.Config{
		Package: "vdp",
	})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(ofile, &SpriteAddr{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(ofile, &SheetPos{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")
}
