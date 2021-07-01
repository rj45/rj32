#include "../cpudef.asm"

move r1,  0b011
move r2,  0
move r3,  1
move r4,  2

shl r1, 0
if.ne r1, 0b011
  error

shl r1, r2
if.ne r1, 0b011
  error

shl r1, 1
if.ne r1, 0b110
  error

shl r1, r3
if.ne r1, 0b1100
  error

move r1,  0b011

shl r1, 2
if.ne r1, 0b1100
  error

move r1,  0b011

shl r1, r4
if.ne r1, 0b1100
  error

halt
