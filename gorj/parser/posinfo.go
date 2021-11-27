// Copyright (c) 2018-2021 TinyGo Authors. All rights reserved.
// Licensed under a 3 clause BSD license. See LICENSE.tinygo.
//
// Copyright (c) 2021 rj45 (github.com/rj45), MIT Licensed, see LICENSE.

package parser

import (
	"go/token"

	"golang.org/x/tools/go/ssa"
)

// posser is an interface that's implemented by both ssa.Value and
// ssa.Instruction. It is implemented by everything that has a Pos() method,
// which is all that getPos() needs.
type posser interface {
	Pos() token.Pos
}

// getPos returns position information for a ssa.Value or ssa.Instruction.
//
// Not all instructions have position information, especially when they're
// implicit (such as implicit casts or implicit returns at the end of a
// function). In these cases, it makes sense to try a bit harder to guess what
// the position really should be.
func getPos(val posser) token.Pos {
	pos := val.Pos()
	if pos != token.NoPos {
		// Easy: position is known.
		return pos
	}

	// No position information is known.
	switch val := val.(type) {
	case *ssa.MakeInterface:
		return getPos(val.X)
	case *ssa.MakeClosure:
		return val.Fn.(*ssa.Function).Pos()
	case *ssa.Return:
		syntax := val.Parent().Syntax()
		if syntax != nil {
			// non-synthetic
			return syntax.End()
		}
		return token.NoPos
	case *ssa.FieldAddr:
		return getPos(val.X)
	case *ssa.IndexAddr:
		return getPos(val.X)
	case *ssa.Slice:
		return getPos(val.X)
	case *ssa.Store:
		return getPos(val.Addr)
	case *ssa.Extract:
		return getPos(val.Tuple)
	case *ssa.If:
		return getPos(val.Cond)
	default:
		// This is reachable, for example with *ssa.Const, *ssa.If, and
		// *ssa.Jump. They might be implemented in some way in the future.
		return token.NoPos
	}
}
