#include "../cpudef.asm"

move r0, 2
move r1, 5

; memory doesn't have 5
load r2, [r0,0]
if.eq r2, 5
error

; store and check it does have 5
store [r0,0], r1
load r2, [r0,0]
if.ne r2, 5
error

; check max offset of 31
move r1, 7
store [r0,31], r1
load r2, [r0,31]
if.ne r2, 7
error

; check base equivalent to offset
move r2, 0
add r0, 15
add r0, 15
add r0, 1
load r2, [r0, 0]
if.ne r2, 7
error

halt
