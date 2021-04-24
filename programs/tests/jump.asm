#include "../cpudef.asm"

  jump one
  error
three:
  add r0, r1
  add r0, r2
  jump end
  error
two:
  move r2, 2
  jump three
  error
one:
  move r1, 3
  jump two
  error
end:
  if.ne r0, 5
  error
  halt