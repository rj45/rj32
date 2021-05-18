#include "../cpudef.asm"

move r1,  0b0110
move r2,  0b1010

or r1, r2
if.ne r1, 0b1110
  error

or r2,   0b0001
if.ne r2, 0b1011
  error

halt
