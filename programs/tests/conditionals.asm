#include "../cpudef.asm"

  brt end
  brf one
three:
  add r0, r1
  add r0, r2
  eq r2, 5
  brt one
  jump end
two:
  move r2, 2
  eq r1, 5
  brf three
one:
  move r1, 3
  eq r1, 3
  brt two
end:
  eq r0, 5
  halt