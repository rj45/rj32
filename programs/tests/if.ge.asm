#include "../cpudef.asm"

move r1, -1
move r2, 0

; skip r-r
if.ge r1, r2
error

; skip r-i
if.ge r1, 0
error

; no-skip r-r
if.ge r2, r1
jump p1
error
p1:

; no-skip r-r equal
if.ge r2, r2
jump p2
error
p2:

; no-skip r-i
if.ge r2, 0
jump p3
error
p3:

; no-skip r-i equal
if.ge r1, -1
jump p4
error
p4:

halt