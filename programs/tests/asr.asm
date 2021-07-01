#include "../cpudef.asm"

move r1,  0b110
move r2,  0
move r3,  1
move r4,  2

asr r1, 0
if.ne r1, 0b110
  error

asr r1, r2
if.ne r1, 0b110
  error

asr r1, 1
if.ne r1, 0b011
  error

asr r1, r3
if.ne r1, 0b001
  error

move r1,  0b1100

asr r1, 2
if.ne r1, 0b0011
  error

move r1,  0b1100

asr r1, r4
if.ne r1, 0b0011
  error

move r1, -8
asr r1, 2
if.ne r1, -2
  error

asr r1, 2
if.ne r1, -1
  error

move r1, -10
asr r1, 1
if.ne r1, -5
  error

halt
