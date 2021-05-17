#include "cpudef.asm"

boot:
  move r0, 8
  add r0, r0 ; 16
  add r0, r0 ; 32
  add r0, r0 ; 64
  add r0, r0 ; 128
  add r0, r0 ; 256
  add r0, r0 ; 512
  add r0, r0 ; 1024
  add r0, r0 ; 2048
  add r0, r0 ; 4096
  add r0, r0 ; 8192
  move r1, r0 ; save
  add r0, r0 ; 16384
  add r0, r1 ; 24576

fibonacci:
  move r1, 0
  move r2, 1
  .loop:
    add r2, r1
    move r3, r2

    if.hs r3, r0
    jump fibonacci

    add r1, r2
    move r3, r1

    jump .loop
    error

