#include "cpudef.asm"

; r1 -- loop counter
; r2 -- next prime candidate
; r3 -- temp
; r4 -- max n
; r7 -- output register

move r4, 8
add r4, r4
add r4, r4
add r4, r4
add r4, r4
add r4, r4
add r4, r4
add r4, r4
add r4, r4

move r1, 1
move r2, 2
move r3, 3
move r5, 5
move r6, 1

initloop:
  add r1, 1
  if.ne r1, r2
    jump .next3
  store [r1,0], r6
  add r2, 2
  .next3:
    if.ne r1, r3
      jump .next5
    store [r1,0], r6
    add r3, 3
  .next5:
    if.ne r1, r5
      jump .finloop
    store [r1,0], r6
    add r5, 5
  .finloop:
  if.lo r1, r4
    jump initloop


sieve:
  move r2, 6
  move r1, 0
  move r3, 1

  ; main loop checking for next prime
  .mainloop:
    add r2, 1

    load r3, [r2, 0]
    if.ne r3, 0
      jump .mainloop

  ; found a prime
    move r7, r2     ; output prime

    if.hs r2, r4
      jump sieve    ; all done, start again

    move r1, r2
    move r3, 1

  .setloop:
    add r1, r2
    if.hs r1, r4
      jump .mainloop
    store [r1, 0], r3
    jump .setloop
