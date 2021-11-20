/*
Copyright (c) 2018-2021 TinyGo Authors. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of the copyright holder nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

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
