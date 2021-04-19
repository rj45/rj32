#include "../cpudef.asm"

move r1, 2
move r2, 3
move r3, 3

; skip r-r
if.ne r2, r2
error

; skip r-i
if.ne r1, 2
error

; no-skip r-r
if.ne r2, r1
jump p1
error
p1:

; no-skip r-i
if.ne r2, 2
jump p2
error
p2:

halt