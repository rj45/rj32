#include "cpudef.asm"

fibonacci:
  move r0, 0
  move r1, 0
  move r2, 1
  .loop:
    add r2, r1
    move r3, r2

    store [r0, 0], r3

    add r1, r2
    move r3, r1

    store [r0, 1], r3
    load r3, [r0,0]
    add r0, 2

    if.eq r0, 14
    jump fibonacci

    jump .loop
    move r2, 0

