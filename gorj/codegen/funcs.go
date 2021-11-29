package codegen

import (
	"fmt"
	"go/constant"
	"go/types"
	"io"
	"log"
	"unicode/utf16"

	"github.com/rj45/rj32/gorj/codegen/asm"
	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
	"github.com/rj45/rj32/gorj/sizes"
)

func (gen *Generator) Func(fn *ir.Func, out io.Writer) *asm.Func {
	gen.out = out
	gen.fn = &asm.Func{
		Comment: fn.Type.String(),
		Label:   fn.Name,
	}

	for _, glob := range fn.Globals {
		if gen.emittedGlobals[glob] {
			continue
		}
		gen.emittedGlobals[glob] = true

		gen.fn.Globals = append(gen.fn.Globals, gen.arch.AssembleGlobal(glob))

		typ := glob.Type.Underlying()
		if ptr, ok := typ.(*types.Pointer); ok {
			typ = ptr.Elem()
		}

		size := sizes.Sizeof(typ)

		if glob.NumArgs() > 0 {
			if gen.section != "data" {
				gen.emit("\n#bank data")
				gen.section = "data"
			}
		} else {
			if gen.section != "bss" {
				gen.emit("\n#bank bss")
				gen.section = "bss"
			}
		}

		name := constant.StringVal(glob.Value)
		gen.emit("%s:  ; %s", name, typ)

		if glob.NumArgs() > 0 {
			data := glob.Arg(0).Value
			if data.Kind() == constant.String {
				str := constant.StringVal(data)

				runes := []rune(str)
				utf16 := utf16.Encode(runes)

				hex := ""
				for i, v := range utf16 {
					if i != 0 && i%8 != 0 {
						hex += ", "
					} else if i != 0 {
						hex += "\n  #d16 "
					}
					hex += fmt.Sprintf("0x%04x", v)
				}

				gen.emit("  #d16 $+2")
				gen.emit("  #d16 %d", len(utf16))
				gen.emit("  ; %q", str)
				gen.emit("  #d16 %s ", hex)
			}
		} else {
			gen.emit("  #res %d", size)
		}
	}

	if gen.section != "code" {
		gen.emit("\n#bank code")
		gen.section = "code"
	}

	gen.emit("\n; %s", fn.Type)
	gen.emit("%s:", fn.Name)

	var retblock *ir.Block

	// order blocks by reverse succession
	blockList := reverseIRSuccessorSort(fn.Blocks()[0], nil, make(map[*ir.Block]bool))

	// reverse it to get succession ordering
	for i, j := 0, len(blockList)-1; i < j; i, j = i+1, j-1 {
		blockList[i], blockList[j] = blockList[j], blockList[i]
	}

	for i, blk := range blockList {
		if blk.Op == op.Return {
			if retblock != nil {
				log.Fatalf("two return blocks! %s", fn.LongString())
			}

			retblock = blk
			continue
		}

		var next *ir.Block
		if (i + 1) < len(blockList) {
			next = blockList[i+1]
		}

		gen.genBlock(blk, next)
	}

	if retblock != nil {
		gen.genBlock(retblock, nil)
	}

	return gen.fn
}

func reverseIRSuccessorSort(block *ir.Block, list []*ir.Block, visited map[*ir.Block]bool) []*ir.Block {
	visited[block] = true

	for i := block.NumSuccs() - 1; i >= 0; i-- {
		succ := block.Succ(i)
		if !visited[succ] {
			list = reverseIRSuccessorSort(succ, list, visited)
		}
	}

	return append(list, block)
}
