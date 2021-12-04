package xform

import (
	"github.com/rj45/rj32/gorj/ir"
)

type Tag uint8

const (
	Invalid Tag = iota
	HasFramePointer

	// ...

	NumTags
)

var activeTags []bool

type Arch interface {
	XformTags() []Tag
}

func SetArch(a Arch) {
	activeTags = make([]bool, NumTags)
	for _, tag := range a.XformTags() {
		activeTags[tag] = true
	}
}

//go:generate enumer -type=Pass

type Pass int

const (
	Elaboration Pass = iota
	Simplification
	Lowering
	Legalize
	CleanUp
)

type xformerDef struct {
	tags []Tag
	fn   func(*ir.Value) int
}

var passes [][]xformerDef

func addToPass(pass Pass, fn func(*ir.Value) int, tags ...Tag) int {
	for int(pass) >= len(passes) {
		passes = append(passes, nil)
	}

	passes[pass] = append(passes[pass], xformerDef{tags: tags, fn: fn})
	return 0
}

func Transform(pass Pass, fn *ir.Func) {
	var xforms []func(*ir.Value) int

next:
	for _, xf := range passes[pass] {
		for _, tag := range xf.tags {
			if !activeTags[tag] {
				continue next
			}
		}

		xforms = append(xforms, xf.fn)
	}

	changes := 1
	tries := 0
nextchange:
	for changes > 0 {
		changes = 0
		tries++
		if tries > 1000 {
			panic("too many tries")
		}
		for _, blk := range fn.Blocks() {
			for i := 0; i < blk.NumInstrs(); i++ {
				instr := blk.Instr(i)

				for _, xform := range xforms {
					changes += xform(instr)
					if changes > 0 {
						continue nextchange
					}
				}
			}
		}
	}
}
