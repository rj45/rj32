package main

import (
	"fmt"
	"os"
	"strings"
)

func genCpudef(filename string) {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	p := func(fmtstr string, args ...interface{}) {
		fmt.Fprintf(f, fmtstr+"\n", args...)
	}

	p("; Code generated by github.com/rj45/rj32/emurj/cmd/isagen. DO NOT EDIT.")
	p("#bits 8")
	p("")

	p("#subruledef reg {")
	for i := 0; i < registers; i++ {
		p("  x%-3d =>% 3d", i, i)
	}
	p("")
	p("  ; ABI names")
	names := strings.Fields(abiNames)
	for i, name := range names {
		p("  %-4s =>% 3d", name, i)
	}
	p("}")
	p("")

	p("#subruledef opcode {")
	for quad := range opmap {
		if quad != 0 {
			p("")
		}
		p("  ; Quadrant %03b", quad)
		for i, name := range opmap[quad] {
			if name == "" {
				continue
			}
			p("  %-7s => 0b%03b_%03b", name, quad, i)
		}
	}
	p("}")
	p("")

	p("#ruledef {")
	for _, name := range fmtOrder {
		args := ""
		argMap := map[string]bool{}
		argMap["opcode"] = true
		for _, operand := range operands[name] {
			args += fmt.Sprintf(", {%s:%s}", operand[0], operand[1])
			argMap[operand[0]] = true
		}
		bits := ""
		for i, field := range fmtFields[name] {
			size := fmt.Sprintf("`%d", 1+fieldBits[field][0]-fieldBits[field][1])
			arg := fieldOperands[field]
			if !argMap[arg] {
				arg = "0"
			} else if fieldOperands[field] == "val" {
				size = fmt.Sprintf("[%d:%d]", immBits[name][field][0], immBits[name][field][1])
			}
			if i != 0 {
				bits += " @ "
			}
			bits += arg + size
		}

		p("  %-5s {opcode:opcode}%s => {\n    le(%s)\n  }", fmtInstr(name), args, bits)
	}

	for quad := range opmap {
		p("")
		p("  ; Quadrant %03b", quad)
		for i, name := range opmap[quad] {
			if name == "" {
				continue
			}
			fm := opcodeToFmt(uint(quad*len(opmap[quad]) + i))

			args := ""
			params := ""
			for j, operand := range operands[fm] {
				if j != 0 {
					args += ", "
				}
				args += fmt.Sprintf("{%s:%s}", operand[0], operand[1])
				if operand[0] != "val" {
					params += fmt.Sprintf(", {%s}", operand[0])
				} else if opcodeIsPcRelative(quad, i) {
					params += ", val - pc - 1"
				} else {
					params += ", val"
				}
			}
			p("  %-7s %s => asm { %s %s%s }", name, args, fmtInstr(fm), name, params)
		}
	}

	p("}")
	p("")
}

func fmtInstr(name string) string {
	return fmt.Sprintf("fmt_%s", strings.ToLower(name))
}
