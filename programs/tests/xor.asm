#include "../cpudef.asm"

move r1,  0b110
move r2,  0b010

xor r1, r2
if.ne r1, 0b100
  error

xor r1, r2
if.ne r1, 0b110
  error

xor r1,   0b011
if.ne r1, 0b101
  error

halt

