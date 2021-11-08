package xform

import (
	"github.com/rj45/rj32/gorj/ir"
)

var passes [][]func(*ir.Value) int

func addToPass(pass int, fn func(*ir.Value) int) int {
	for pass >= len(passes) {
		passes = append(passes, nil)
	}

	passes[pass] = append(passes[pass], fn)
	return 0
}

func Transform(fn *ir.Func) {
	for pass := range passes {
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
}
