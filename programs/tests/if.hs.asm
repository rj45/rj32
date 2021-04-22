#include "../cpudef.asm"

move r1, 0
move r2, -1 ; 0xffff

; skip r-r
if.hs r1, r2
error

; skip r-i
if.hs r1, -1
error

; no-skip r-r
if.hs r2, r1
jump p1
error
p1:

; no-skip r-r equal
if.hs r2, r2
jump p2
error
p2:

; no-skip r-i
if.hs r2, -1
jump p3
error
p3:

; no-skip r-i equal
if.hs r1, 0
jump p4
error
p4:

halt