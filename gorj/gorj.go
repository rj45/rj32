// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/rj45/rj32/gorj/arch"
	"github.com/rj45/rj32/gorj/compiler"

	// load the supported architectures so they register with the arch package
	_ "github.com/rj45/rj32/gorj/arch/a32"
	_ "github.com/rj45/rj32/gorj/arch/rj32"
)

type dumper interface {
	WritePhase(string, string)
	WriteAsm(string, *bytes.Buffer)
	Close()
}

var output = flag.String("o", "", "output file for the result")
var dir = flag.String("c", "", "set working dir (default current dir)")
var theArch = flag.String("arch", "", "architecture to compile for")

func main() {
	log.SetFlags(log.Lshortfile)

	flag.Parse()

	command := flag.Arg(0)
	printUsage := flag.NArg() < 2

	assemble := false
	run := false
	_ = run

	switch command {
	case "b", "build":
		assemble = true
	case "r", "run":
		assemble = true
		run = true
	case "s", "asm":
	default:
		printUsage = true
	}

	if printUsage {
		fmt.Fprintln(os.Stderr, "gorj - a go compiler for rj32")
		fmt.Fprintln(os.Stderr, "https://github.com/rj45/rj32/gorj")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Usage: gorj <flags> <command> <packages...>")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Commands:")
		fmt.Fprintln(os.Stderr, "  build: compile and assemble with customasm")
		fmt.Fprintln(os.Stderr, "  asm: compile and write assembly to file")
		fmt.Fprintln(os.Stderr, "  run: compile, assemble and run emulator")
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "Flags:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *dir == "" {
		cwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		*dir = cwd
	}

	outname := "-"
	if *output != "" {
		outname = *output
	}

	if *theArch != "" {
		arch.SetArch(*theArch)
	}

	result := compiler.Compile(outname, *dir, flag.Args()[1:], assemble, run)

	os.Exit(result)
}
