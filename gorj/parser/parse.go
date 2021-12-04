package parser

import (
	"bytes"
	"go/token"
	"log"
	"os"
	"strings"

	"github.com/rj45/rj32/gorj/ir"
	"golang.org/x/tools/go/ssa"
)

type Arch interface {
	Name() string
}

var arch Arch

func SetArch(a Arch) {
	arch = a
}

type Parser struct {
	pkg      *ir.Package
	members  []ssa.Member
	ssaFuncs map[*ir.Func]*ssa.Function
	parsed   map[*ir.Func]bool
}

func NewParser(dir string, patterns ...string) *Parser {
	members, err := parseProgram(dir, patterns...)
	if err != nil {
		log.Fatal(err)
	}

	return &Parser{pkg: &ir.Package{}, members: members}
}

func (p *Parser) Scan() {
	p.ssaFuncs = make(map[*ir.Func]*ssa.Function)
	for _, member := range p.members {
		switch member.Token() {
		case token.FUNC:
			fn := member.Package().Func(member.Name())
			name := genName(fn.Pkg.Pkg.Name(), fn.Name())
			referenced := name == "main__main" || name == "main__init"
			irFunc := &ir.Func{
				Name:       genName(fn.Pkg.Pkg.Name(), fn.Name()),
				Type:       fn.Signature,
				Pkg:        p.pkg,
				Referenced: referenced,
			}
			p.pkg.Funcs = append(p.pkg.Funcs, irFunc)
			p.ssaFuncs[irFunc] = fn

		case token.VAR:
			name := genName(member.Package().Pkg.Name(), member.Name())
			p.pkg.AddGlobal(name, member.Type())
		case token.TYPE:
		case token.CONST:
		default:
			log.Fatalln("unknown type", member.Token())
		}
	}
	p.parsed = make(map[*ir.Func]bool)
}

func (p *Parser) Package() *ir.Package {
	return p.pkg
}

func (p *Parser) NextUnparsedFunc() *ir.Func {
	for _, fn := range p.pkg.Funcs {
		if fn.Referenced && !p.parsed[fn] {
			return fn
		}
	}

	return nil
}

func (p *Parser) DumpOrignalSource(fn *ir.Func) (filename string, lines []string, startline int) {
	ssafn := p.ssaFuncs[fn]
	fset := ssafn.Prog.Fset

	if ssafn.Syntax() == nil {
		return
	}

	start := ssafn.Syntax().Pos()
	end := ssafn.Syntax().End()

	if start == token.NoPos || end == token.NoPos {
		return
	}

	startp := fset.PositionFor(start, true)
	filename = startp.Filename
	startline = startp.Line - 1

	endp := fset.PositionFor(end, true)
	buf, err := os.ReadFile(startp.Filename)
	if err != nil {
		log.Fatal(err)
	}
	lines = strings.Split(string(buf), "\n")
	lines = lines[startline:endp.Line]

	return
}

func (p *Parser) DumpOriginalSSA(fn *ir.Func) *bytes.Buffer {
	buf := &bytes.Buffer{}
	ssa.WriteFunction(buf, p.ssaFuncs[fn])
	return buf
}

func (p *Parser) ParseFunc(fn *ir.Func) {
	p.parsed[fn] = true
	walkFunc(fn, p.ssaFuncs[fn])
}
