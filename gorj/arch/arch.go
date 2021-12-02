package arch

import (
	"strings"

	"github.com/rj45/rj32/gorj/codegen"
	"github.com/rj45/rj32/gorj/ir/reg"
	"github.com/rj45/rj32/gorj/sizes"
)

const defaultArch = "rj32"

type Architecture interface {
	codegen.Arch
	reg.Arch
	sizes.Arch
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
	if name == defaultArch {
		SetArch(name)
	}
	return 0
}

func SetArch(name string) {
	arch = arches[strings.ToLower(name)]
	reg.SetArch(arch)
	codegen.SetArch(arch)
	sizes.SetArch(arch)
}
