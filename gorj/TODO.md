# Things left to do

- [x] build register allocation verification
- [x] fix load store bugs
- [x] properly handle parameters with a copy
- [ ] top to bottom register allocation verification?
- [ ] fix value overwrites in the register allocator
- [ ] switch allocator to use uses data structure for liveness tracking
- [ ] if a function has multiple returns
  - [ ] create a common exit block
  - [ ] set successors for each return block to it
  - [ ] create a phi in the return block for each returned value
  - [ ] change return to jump, remove args
- [ ] create spills for temp vars still live at a call
- [ ] build way to count register pressure
  - [ ] add spill support to lower register pressure
