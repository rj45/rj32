package parser

import (
	"fmt"
	"sort"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

type members []ssa.Member

func (m members) Len() int           { return len(m) }
func (m members) Swap(i, j int)      { m[i], m[j] = m[j], m[i] }
func (m members) Less(i, j int) bool { return m[i].Pos() < m[j].Pos() }

func parseProgram(dir string, patterns ...string) ([]ssa.Member, error) {
	// Load, parse, and type-check the whole program.
	cfg := packages.Config{
		Mode: packages.LoadAllSyntax,
		Dir:  dir,
	}
	initial, err := packages.Load(&cfg, patterns...)
	if err != nil {
		return nil, err
	}

	// Print any errors that happened in the build process
	if packages.PrintErrors(initial) > 0 {
		return nil, fmt.Errorf("initial package loading had errors")
	}

	// Create SSA packages for well-typed packages and their dependencies.
	prog, _ := ssautil.AllPackages(initial, ssa.BuildSerially)

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
