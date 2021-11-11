package xform

import (
	"github.com/rj45/rj32/gorj/ir"
)

type Pass int

const (
	FirstPass Pass = iota
	LastPass
)

var passes [][]func(*ir.Value) int

func addToPass(pass Pass, fn func(*ir.Value) int) int {
	for int(pass) >= len(passes) {
		passes = append(passes, nil)
	}

	passes[pass] = append(passes[pass], fn)
	return 0
}

func Transform(pass Pass, fn *ir.Func) {
	changes := 1
	tries := 0
	for changes > 0 {
		changes = 0
		tries++
		if tries > 10 {
			panic("too many tries")
		}
		for _, blk := range fn.Blocks {
			for i := 0; i < len(blk.Instrs); i++ {
				for _, xform := range passes[pass] {
					changes += xform(blk.Instrs[i])
				}
			}
		}
	}
}
