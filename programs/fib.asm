#include "cpudef.asm"

boot:
  move r1, 8
  add r1, r1 ; 16
  add r1, r1 ; 32
  add r1, r1 ; 64
  add r1, r1 ; 128
  add r1, r1 ; 256
  add r1, r1 ; 512
  add r1, r1 ; 1024
  add r1, r1 ; 2048
  add r1, r1 ; 4096
  add r1, r1 ; 8192
  move r2, r1 ; save
  add r1, r1 ; 16384
  add r1, r2 ; 24576

  and r2, 12
  or r2, 6
  xor r2, -1
  sub r2, 19

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
