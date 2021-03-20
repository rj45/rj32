#include "../cpudef.asm"

  jump one
three:
  add r0, r1
  add r0, r2
  jump end
two:
  move r2, 2
  jump three
one:
  move r1, 3
  jump two
end:
  musteq r0, 5