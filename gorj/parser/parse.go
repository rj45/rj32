package parser

import (
	"fmt"
	"go/constant"
	"go/token"
	"log"
	"path/filepath"

	"github.com/rj45/rj32/gorj/ir"
	"github.com/rj45/rj32/gorj/ir/op"
)

func ParseModule(filename string) *ir.Module {
	log.SetFlags(log.Lshortfile)

	dir := filepath.Dir(filename)
	basename := filepath.Base(filename)

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
			name := fmt.Sprintf("%s.%s", member.Package().Pkg.Name(), member.Name())
			mod.Globals = append(mod.Globals, &ir.Value{
				Op:    op.Global,
				Value: constant.MakeString(name),
				Type:  member.Type(),
			})
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
