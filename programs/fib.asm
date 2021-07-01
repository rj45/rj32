#include "cpudef.asm"

boot:
  move r1, 1
  shl r1, 13
  move r2, r1 ; save 8192
  shl r1, 1
  add r1, r2 ; 24576

  shr r2, 5
  shl r2, 4
  move r2, -42
  asr r2, 1
  asr r2, 1

.repeat:
  call fibonacci
  move r0, 0
  jump .repeat

fibonacci:
  move r2, 0
  move r3, 1
  .loop:
    add r3, r2
    move r7, r3

    if.hs r7, r1
      return

    add r2, r3
    move r7, r2

    jump .loop
    error
