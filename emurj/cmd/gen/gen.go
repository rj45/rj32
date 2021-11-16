package main

import (
	"os"

	"github.com/rj45/rj32/emurj/bitfield"
)

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

type Fmt int
type Op int

type Inst struct {
	Fmt Fmt `bitfield:"4"`
	Op  Op  `bitfield:"5"`
	Rd  int `bitfield:"4"`
	Rs  int `bitfield:"4"`
	Imm int `bitfield:"13"`
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

type InstI12 struct {
	Imm int `bitfield:"12"`
	Op  Op  `bitfield:"1"`
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
	dfile, err := os.Create("data/fields.go")
	if err != nil {
		panic(err)
	}
	defer dfile.Close()

	err = bitfield.Gen(dfile, &Bus{}, &bitfield.Config{
		Package: "data",
	})
	if err != nil {
		panic(err)
	}
	dfile.WriteString("\n")

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

	err = bitfield.Gen(ofile, &Tile{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(ofile, &Pixels{}, &bitfield.Config{})
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

	err = bitfield.Gen(cfile, &InstI12{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")

	err = bitfield.Gen(cfile, &InstRI8{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	ofile.WriteString("\n")
}
