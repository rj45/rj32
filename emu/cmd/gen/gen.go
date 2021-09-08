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

type Fmt int
type Op int

type Inst struct {
	Fmt Fmt `bitfield:"3"`
	Op  Op  `bitfield:"5"`
	Rd  int `bitfield:"4"`
	Rs  int `bitfield:"4"`
	Imm int `bitfield:"12"`
}

type InstRI6 struct {
	Rd  int `bitfield:"4"`
	Imm int `bitfield:"6"`
	Op  Op  `bitfield:"4"`
	Fmt Fmt `bitfield:"2"`
}

type InstRR struct {
	Rd  int `bitfield:"4"`
	Rs  int `bitfield:"4"`
	NA  int `bitfield:"1"`
	Op  Op  `bitfield:"5"`
	Fmt Fmt `bitfield:"2"`
}

type InstLS struct {
	Rd  int `bitfield:"4"`
	Rs  int `bitfield:"4"`
	Imm int `bitfield:"4"`
	Op  Op  `bitfield:"2"`
	Fmt Fmt `bitfield:"2"`
}

type InstI11 struct {
	Imm int `bitfield:"11"`
	Op  Op  `bitfield:"2"`
	Fmt Fmt `bitfield:"3"`
}

type InstRI8 struct {
	Rd  int `bitfield:"4"`
	Imm int `bitfield:"8"`
	Op  Op  `bitfield:"1"`
	Fmt Fmt `bitfield:"3"`
}

type Bus struct {
	Data    int  `bitfield:"16"`
	Address int  `bitfield:"21"`
	WE      bool `bitfield:"1"`
	Ack     bool `bitfield:"1"`
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

	cfile, err := os.Create("rj32/instfmts.go")
	if err != nil {
		panic(err)
	}
	defer cfile.Close()

	err = bitfield.Gen(cfile, &Inst{}, &bitfield.Config{
		Package: "rj32",
	})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &InstRI6{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &InstRR{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &InstLS{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &InstI11{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &InstRI8{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &Bus{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")
}
