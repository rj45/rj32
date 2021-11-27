package parser

import (
	"go/token"
	"log"

	"github.com/rj45/rj32/gorj/ir"
	"golang.org/x/tools/go/ssa"
)

func ParseProgram(dir string, patterns ...string) *ir.Package {
	members, err := parseProgram(dir, patterns...)
	if err != nil {
		log.Fatal(err)
	}

	pkg := &ir.Package{}

	walk(pkg, members)

	return pkg
}

func walk(pkg *ir.Package, all members) {
	ssaFuncs := make(map[*ir.Func]*ssa.Function)
	for _, member := range all {
		switch member.Token() {
		case token.FUNC:
			fn := member.Package().Func(member.Name())
			name := genName(fn.Pkg.Pkg.Name(), fn.Name())
			referenced := name == "main__main" || name == "main__init"
			irFunc := &ir.Func{
				Name:       genName(fn.Pkg.Pkg.Name(), fn.Name()),
				Type:       fn.Signature,
				Pkg:        pkg,
				Referenced: referenced,
			}
			pkg.Funcs = append(pkg.Funcs, irFunc)
			ssaFuncs[irFunc] = fn

		case token.VAR:
			name := genName(member.Package().Pkg.Name(), member.Name())
			pkg.AddGlobal(name, member.Type())
		case token.TYPE:
		case token.CONST:
		default:
			log.Fatalln("unknown type", member.Token())
		}
	}

	parsed := make(map[*ir.Func]bool)
	changes := 1

	for changes > 0 {
		changes = 0
		for _, fn := range pkg.Funcs {
			if fn.Referenced && !parsed[fn] {
				parsed[fn] = true
				changes++
				walkFunc(fn, ssaFuncs[fn])
			}
		}
	}
}
