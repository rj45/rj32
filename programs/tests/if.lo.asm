#include "../cpudef.asm"

move r1, 0
move r2, -1 ; 0xffff

; skip r-r
if.lo r2, r1
error

; skip r-r equal
if.lo r2, r2
error

; skip r-i
if.lo r2, 0
error

; skip r-i equal
if.lo r2, -1
error

; no-skip r-r
if.lo r1, r2
jump p1
error
p1:

; no-skip r-i
if.lo r1, 1
jump p2
error
p2:

halt