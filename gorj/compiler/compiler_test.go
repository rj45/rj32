package compiler_test

import (
	"testing"

	"github.com/rj45/rj32/gorj/compiler"
)

func TestCompiler(t *testing.T) {
	testCases := []struct {
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
			desc:     "putc",
			filename: "./putc/putc.go",
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
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			// Compile(outname, dir string, patterns []string, assemble, run bool) int
			result := compiler.Compile("-", "../testdata/", []string{tC.filename}, true, true)
			if result != 0 {
				t.Errorf("test %s failed with code %d", tC.filename, result)
			}
		})
	}
}
