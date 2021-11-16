# Things left to do

- [x] build register allocation verification
- [x] fix load store bugs
- [x] properly handle parameters with a copy
- [x] switch allocator to use uses data structure for liveness tracking
  - [x] in uses liveness calculator, special case phis
    - [x] phis should act like copies in source block
    - [x] if the value does not live
      - [x] it should not be in live out
      - [x] it should be in block kill list
- [x] top to bottom register allocation verification?
- [x] it seems like back-links (loops) are not handled properly in the register allocator yet
- [x] fix value overwrites in the register allocator
- [x] get ssa dump html to have final assembly output
- [.] get a test suite working
  - [x] ability to write .asm files
  - [x] can invoke customasm to produce .hex files
  - [x] can invoke emulator to run the .hex files
  - [ ] set up a test runner
  - [ ] start writing some simple test programs inspired by c-test-suite
- [x] storing a const needs a register `store  [gp, main__init_guard], 1`
- [.] implement function parameters
  - [x] add parameters to function entry
  - [ ] handle stack and register ABI for parameters
- [ ] create spills for temp vars still live at a call
  - [ ] reload after
- [ ] if a function has multiple returns
  - [ ] create a common exit block
  - [ ] set successors for each return block to it
  - [ ] create a phi in the return block for each returned value
  - [ ] change return to jump, remove args

## milestone get: min useful compiler

- now has ability to compile simple programs that use only word data
  - global structs and arrays supported
  - allocator does not support spills
    - crashes if it needs more registers than available


- [ ] emit ssa dump html on command line arg
- [ ] implement graph coloring in register allocator
  - [ ] implement adjacency lists for the graph
    - [ ] nodes can be marked as move related
    - [ ] nodes can be merged so multiple values are in one node
    - [ ] nodes can be removed one-by-one
    - [ ] fast finding of the nodes with least degree
  - [ ] use iterated register coalescing algo
- [ ] implement multiple function return values
- [ ] expand various ops into calls to functions implementing them
- [ ] add runtime library support for builtin ops
  - [ ] implement mul in go
  - [ ] implement div in go
  - [ ] implement rem in go
  - [ ] implement double and quad word ops
    - [ ] add/sub should use addc/subc
    - [ ] shifts should do a function call
    - [ ] either multi-def support or register pair support
- [ ] stack allocation
  - [ ] slices
- [ ] heap allocation
  - [ ] free somehow?
- [ ] slice support
- [ ] build way to count register pressure
  - [ ] add spill support to lower register pressure

## Optimizations

- [ ] fix allocator choosing wrong variable and doing extra copies
- [ ] copy propagation
  - given use of X:
    - are all reaching definitions of X
      - copies from the same variable, ie X = copy Y
    - where Y is not redefined since that variable?
    - if so, substitute use of X with use of Y instead
- [ ] move instruction defs closer to first uses to minimize reg pressure
- [ ] constant folding
- [ ] common subexpression elimination
- [ ] dead code elimination
- [ ] find all loops
  - [ ] loop invariant code motion
    - [ ] if for def X, no args refer to a phi node or def inside the loop
      - [ ] move X out of the loop into the pre-header
  - [ ] find loop induction variables?
    - [ ] do strength reduction on uses of induction variable?
      - like array indexing for example
  - [ ] support for inline assembly / extern assembly
    - [ ] get mul, div and rem converted to assembly
