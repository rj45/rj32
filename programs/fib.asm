#include "cpudef.asm"

fibonacci:
  add r2, 1
  .loop:
    add r2, r1
    move r3, r2
    add r1, r2
    move r3, r1
    jump .loop

