#include "../cpudef.asm"

move r1, 2
move r2, 3
move r3, 3

; skip r-r
if.eq r1, r2
error

; skip r-i
if.eq r1, 3
error

; no-skip r-r
if.eq r2, r3
jump p1
error
p1:

; no-skip r-i
if.eq r2, 3
jump p2
error
p2:

halt