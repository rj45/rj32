package parser

import (
	"go/token"
	"log"

	"github.com/rj45/rj32/gorj/ir"
)

func ParseProgram(dir, basename string) *ir.Package {
	log.SetFlags(log.Lshortfile)

	members, err := parseProgram(dir, basename)
	if err != nil {
		log.Fatal(err)
	}

	pkg := &ir.Package{}

	walk(pkg, members)

	return pkg
}

func walk(pkg *ir.Package, all members) {
	for _, member := range all {
		if member.Token() == token.VAR {
			name := genName(member.Package().Pkg.Name(), member.Name())
			pkg.AddGlobal(name, member.Type())
		}
	}

	for _, member := range all {
		switch member.Token() {
		case token.FUNC:
			walkFunc(pkg, member.Package().Func(member.Name()))
		case token.VAR:
		case token.TYPE:
		case token.CONST:
		default:
			log.Fatalln("unknown type", member.Token())
		}
	}
}
