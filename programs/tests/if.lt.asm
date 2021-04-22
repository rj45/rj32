#include "../cpudef.asm"

move r1, -1
move r2, 0

; skip r-r
if.lt r2, r1
error

; skip r-r equal
if.lt r2, r2
error

; skip r-i
if.lt r2, -1
error

; skip r-i equal
if.lt r2, 0
error

; no-skip r-r
if.lt r1, r2
jump p1
error
p1:

; no-skip r-i
if.lt r1, 1
jump p2
error
p2:

halt