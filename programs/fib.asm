#include "cpudef.asm"

fibonacci:
  move r1, 0
  move r2, 1
  .loop:
    add r2, r1
    move r3, r2
    add r1, r2
    move r3, r1
    eq r3, 8
    brt fibonacci
    jump .loop

