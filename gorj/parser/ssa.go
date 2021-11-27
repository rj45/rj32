package parser

import (
	"fmt"
	"os"
	"sort"

	"github.com/rj45/rj32/gorj/goenv"
	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type members []ssa.Member

func (m members) Len() int           { return len(m) }
func (m members) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m members) Less(i, j int) bool { return m[i].Pos() < m[j].Pos() }

func parseProgram(dir string, patterns ...string) ([]ssa.Member, error) {
	goroot, err := goenv.GetCachedGoroot()
	if err != nil {
		return nil, err
	}

	// Load, parse, and type-check the whole program.
	cfg := packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  dir,
		Env:  append(os.Environ(), "GOROOT="+goroot),
	}

	initial, err := packages.Load(&cfg, patterns...)
	if err != nil {
		return nil, err
	}

	hasRuntime := false
	var main *packages.Package
	for _, pkg := range initial {
		if pkg.Name == "main" {
			main = pkg
		}
		if hasRuntimePackage(pkg) {
			hasRuntime = true
			break
		}
	}

	if !hasRuntime && main != nil {
		cfg.Fset = main.Fset

		rt, err := packages.Load(&cfg, "runtime")
		if err != nil {
			return nil, err
		}

		if packages.PrintErrors(rt) > 0 {
			return nil, fmt.Errorf("runtime loading had errors")
		}

		main.Imports["runtime"] = rt[0]
	}

	// Print any errors that happened in the build process
	if packages.PrintErrors(initial) > 0 {
		return nil, fmt.Errorf("initial package loading had errors")
	}

	// Create SSA packages for well-typed packages and their dependencies.
	prog, pkgs := ssautil.AllPackages(initial, ssa.BuildSerially)

	// for _, pkg := range pkgs {
	// 	pkg.SetDebugMode(true)
	// }
	_ = pkgs

	// Build SSA code for the whole program.
	prog.Build()

	members := members([]ssa.Member{})
	for _, pkg := range prog.AllPackages() {
		for _, member := range pkg.Members {
			members = append(members, member)
		}
	}

	// Sort by Pos()
	sort.Sort(members)

	return members, nil
}

func hasRuntimePackage(pkg *packages.Package) bool {
	if pkg.Name == "runtime" {
		return true
	}
	for _, pkg := range pkg.Imports {
		if hasRuntimePackage(pkg) {
			return true
		}
	}
	return false
}
