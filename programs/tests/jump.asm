#include "../cpudef.asm"

  jump one
  halt
three:
  add r0, r1
  add r0, r2
  jump end
  halt
two:
  move r2, 2
  jump three
  halt
one:
  move r1, 3
  jump two
  halt
end:
  eq r0, 5
  halt