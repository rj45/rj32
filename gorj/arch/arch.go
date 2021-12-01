package arch

import (
	"github.com/rj45/rj32/gorj/codegen"
	"github.com/rj45/rj32/gorj/ir/reg"
)

type Architecture interface {
	codegen.Arch
	reg.Arch
}

var arch Architecture

var arches map[string]Architecture

func Arch() Architecture {
	return arch
}

func Register(name string, a Architecture) int {
	if arches == nil {
		arches = make(map[string]Architecture)
	}
	arches[name] = a
	SetArch(name)
	return 0
}

func SetArch(name string) {
	arch = arches[name]
	reg.SetArch(arch)
	codegen.SetArch(arch)
}
