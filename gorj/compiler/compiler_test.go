package compiler_test

import (
	"testing"

	"github.com/rj45/rj32/gorj/arch"
	"github.com/rj45/rj32/gorj/compiler"

	// load the supported architectures so they register with the arch package
	_ "github.com/rj45/rj32/gorj/arch/a32"
	_ "github.com/rj45/rj32/gorj/arch/rj32"
)

var testCases = []struct {
	desc     string
	filename string
}{
	{
		desc:     "simple test",
		filename: "./simple/simple.go",
	},
	{
		desc:     "hello world",
		filename: "./hello/hello.go",
	},
	{
		desc:     "fibonacci",
		filename: "./fib/fib.go",
	},
	{
		desc:     "multiply and divide",
		filename: "./muldiv/muldiv.go",
	},
	{
		desc:     "incremental seive of eratosthenes",
		filename: "./seive/seive.go",
	},
	{
		desc:     "n queens problem",
		filename: "./nqueens/nqueens.go",
	},
	{
		desc:     "multiple return values",
		filename: "./multireturn/multireturn.go",
	},
	{
		desc:     "iterating and storing strings",
		filename: "./print/print.go",
	},
	{
		desc:     "external assembly",
		filename: "./externasm/externasm.go",
	},
}

func TestCompilerForRj32(t *testing.T) {
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			arch.SetArch("rj32")
			result := compiler.Compile("-", "../testdata/", []string{tC.filename}, true, true)
			if result != 0 {
				t.Errorf("test %s failed with code %d", tC.filename, result)
			}
		})
	}
}

func TestCompilerForA32(t *testing.T) {
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			arch.SetArch("a32")
			result := compiler.Compile("-", "../testdata/", []string{tC.filename}, true, true)
			if result != 0 {
				t.Errorf("test %s failed with code %d", tC.filename, result)
			}
		})
	}
}
