# Things left to do

- [x] build register allocation verification
- [x] fix load store bugs
- [x] properly handle parameters with a copy
- [.] switch allocator to use uses data structure for liveness tracking
  - [ ] in uses liveness calculator, special case phis
    - [ ] phis should act like copies in source block
    - [ ] if the value does not live
      - [ ] it should not be in live out
      - [ ] it should be in block kill list
  - [ ] alternately:
    - [ ] color each value in each block in one pass (not with regs, just colors)
    - [ ] join colors at the block level across the function
    - [ ] use affinity list to try to merge colors for copies
    - [ ] only at end choose registers for colors
- [ ] top to bottom register allocation verification?
- [ ] it seems like back-links (loops) are not handled properly in the register allocator yet
- [ ] fix value overwrites in the register allocator

- [ ] if a function has multiple returns
  - [ ] create a common exit block
  - [ ] set successors for each return block to it
  - [ ] create a phi in the return block for each returned value
  - [ ] change return to jump, remove args
- [ ] create spills for temp vars still live at a call
- [ ] build way to count register pressure
  - [ ] add spill support to lower register pressure
