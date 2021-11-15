package parser

import (
	"go/token"
	"log"

	"github.com/rj45/rj32/gorj/ir"
)

func ParseModule(dir, basename string) *ir.Module {
	log.SetFlags(log.Lshortfile)

	members, err := parseProgram(dir, basename)
	if err != nil {
		log.Fatal(err)
	}

	mod := &ir.Module{}

	walk(mod, members)

	return mod
}

func walk(mod *ir.Module, all members) {
	for _, member := range all {
		if member.Token() == token.VAR {
			name := genName(member.Package().Pkg.Name(), member.Name())
			mod.AddGlobal(name, member.Type())
		}
	}

	for _, member := range all {
		switch member.Token() {
		case token.FUNC:
			walkFunc(mod, member.Package().Func(member.Name()))
		case token.VAR:
		case token.TYPE:
		case token.CONST:
		default:
			log.Fatalln("unknown type", member.Token())
		}
	}
}
