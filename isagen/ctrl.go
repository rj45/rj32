package main

import (
	"os"

	"github.com/rj45/rj32/isagen/bitfield"
)

type ExCtrl struct {
	Sub   bool `bitfield:"1"`
	AluOp int  `bitfield:"2"`
}

type WbCtrl struct {
	RegWen bool `bitfield:"1"`
	WbMux  int  `bitfield:"1"`
}

func genCtrl(filename string) {
	dfile, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer dfile.Close()

	err = bitfield.Gen(dfile, &ExCtrl{}, &bitfield.Config{
		Package: "ctrl",
	})
	if err != nil {
		panic(err)
	}
	dfile.WriteString("\n")

	err = bitfield.Gen(dfile, &WbCtrl{}, &bitfield.Config{})
	if err != nil {
		panic(err)
	}
	dfile.WriteString("\n")
}
