# gorj

Gorge yourself on `gorj`, a Go compiler for `rj32`, a 16-bit homebrew CPU.

This code is meant to be forked and converted to whatever homebrew/hobby CPU you want to run it on. Please, have at it, even change the name if you like. The MIT license asks you keep the copyright intact, just add your own copyright line right under it and copy the files into your own git repo and modify away.

If you implement any cool optimizations, if you feel like contributing them back, I would be most appreciative! But it's up to you, do what you like.

Note: Some parts are from the Go compiler, and under a different license and copyrighted by the Go Authors.

## What is it?

This compiler will take a Go package, read in all the packages it depends on in the usual way Go programs work, and compile all the code down into assembler in a style that is compatible with [customasm](https://github.com/hlorenzi/customasm).

As of this writing, customasm does not support linking, so a single large assembly file is produced. A "CPU Def" file can be included which configures the assembly language, as well as the memory layout with `#bank`s.

## Why Go?

C is great, but the language is not the easiest to parse, and while there's many great projects like [LCC](https://github.com/drh/lcc), they are not the easiest to work on and modify for a homebrew CPU.

Go is a very simple language, and is very fast to learn, yet powerful enough to be very productive.

As well, Go has a parser and type checker for Go built into the standard library. Further, there's an excellent [SSA library](https://golang.org/x/tools/go/ssa) that does most of the work of the frontend of a compiler. So all that's really required is a simple backend, which is what gorj is.

## Limitations

While all of Go is parsed, obviously there's no garbage collector, and no goroutines. In that respect it's a lot more like [tinygo](https://tinygo.org/) than Go.

It took a team many years to build the Go compiler, so as this is a hobby project written by one person, there's a lot left completely unimplemented, and error handling could be a lot better.

Defer is ignored, though it could be implemented in the future. Panic/recover will not be implemented. Runtime type reflection will not be implemented. Maps are not yet implemented.

Ints, uints, bytes, runes and pointers are 16-bits, which is fairly easy to change, see the `sizes` package. But non-standard sizes can violate some assumptions in the standard library, so anything relying on those assumptions will have bugs.

Almost none of the standard library is supported. You can try it and see if it will work, but some fundamental assumptions are violated, as well as many features relied are missing. This compiler is meant to help you write your own standard library, kernel, OS and other software for your own homebrew CPU, so you could see the lack of a standard library as a feature.

## Design

Gorj uses the [golang.org/x/tools/go/ssa](https://golang.org/x/tools/go/ssa) package to build and parse Go into [SSA form](https://en.wikipedia.org/wiki/Static_single_assignment_form).

This IR is then [converted](./parser) into a different [IR](./ir) which is [easier to work with](https://www.hpl.hp.com/techreports/Compaq-DEC/WRL-2000-2.html) and [transform](./xform).

After transformation there is a [SSA based tree register allocator](./regalloc) which uses the properties of SSA to [allocate registers in linear time](https://compilers.cs.uni-saarland.de/projects/ssara/).

After a final transformation pass, then finally [assembly code is generated](./codegen) and output.

## License

Copyright (c) 2021 rj45 ([github.com/rj45](https://github.com/rj45))

Licensed under the MIT License (MIT), see [LICENSE](./LICENSE).

Some parts Copyrighted by the Go Authors, under a BSD like license, see [Go's LICENSE](https://github.com/golang/go/blob/666fc173c02ff3004ac9ef867aa4eec7e243dde3/LICENSE).
