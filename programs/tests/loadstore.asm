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

; check max offset of 15
move r1, 7
store [r0,15], r1
load r2, [r0,15]
if.ne r2, 7
error

; check base equivalent to offset
move r2, 0
add r0, 15
load r2, [r0, 0]
if.ne r2, 7
error

; check multiple writes and reads in a row
move r1, 1
move r2, 2
move r3, 3
move r4, 4
move r5, 5
move r6, 6
move r0, 0
store [r0, 3], r3
store [r0, 1], r1
store [r0, 6], r6
store [r0, 2], r2
store [r0, 4], r4
store [r0, 5], r5
move r1, -1
move r2, -1
move r3, -1
move r4, -1
move r5, -1
move r6, -1
load r5, [r0, 5]
load r1, [r0, 1]
load r6, [r0, 6]
load r2, [r0, 2]
load r4, [r0, 4]
load r3, [r0, 3]
if.ne r1, 1
  error
if.ne r2, 2
  error
if.ne r3, 3
  error
if.ne r4, 4
  error
if.ne r5, 5
  error
if.ne r6, 6
  error

load r3, [r0, 3]
jump skiperr1
error
skiperr1:

store [r0, 3], r3
jump skiperr2
error
skiperr2:

halt
