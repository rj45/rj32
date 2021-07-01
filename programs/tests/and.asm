#include "../cpudef.asm"

move r1,  0b0110
move r2,  0b1010

and r1, r2
if.ne r1, 0b0010
  error

and r2,   0b1101
if.ne r2, 0b1000
  error

halt

