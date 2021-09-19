#include "../cpudef.asm"

; put the upper 11 bits of 0xffff into r3
imm 0xffff
move r3, 0

; r3 should now equal -32 or 0xffe0
; lets invert that
xor r3, -1

; r3 should now be 0xf or 0b1111
if.ne r3, 0b1111
  error

; construct a 16 bit value in r1 without using the imm instruction
move r1, 0x75
shl r1, 8
move r2, 0x49
or r1, r2

; trigger imm to augment the if.ne instruction with a 16-bit value
if.ne r1, 0x7549
  error

; try imm augmentation on add in addition to if.ne
move r5, 70
add r5, 128
if.ne r5, 70+128
  error

; make sure 63 doesn't become -1 in the if.ne
move r5, 63
if.ne r5, 63
  error

; handle the special case of skipping the imm instruction
; in the first if.ne, the add instruction will have an imm prefix
; we want to skip that imm prefix and the add
; otherwise the assembler needs detect this case and perform a very
; complicated workaround that customasm doesn't support
move r5, 10
if.ne r5, 10
  add r5, 100
if.ne r5, 10
  error

move r5, 0x7e12
if.ne r5, 0x7e12
  error

; grab the pc and put it in r4
call grabpc
grabpc:
move r4, ra   ; 1

move r5, 0x7F ; 2

; grab the pc and put it in r2
call grabpc2  ; 3
grabpc2:
move r2, ra

; subtract pc values
sub r2, r4

; assert there is not a imm instruction inserted
if.ne r2, 3
 error

move r2, 32
move r1, 0
move r4, 0x1234

; check imm prefix is handled on load properly
store [r2], r4
imm 32
load r3, [r1, 0]
if.ne r3, 0x1234
  error

; check imm prefix for store
; move r4, 0x5678
; store [r1, 17], r4
; nop
; load r3, [r2]
; nop
; if.ne r3, 0x5678
;   error


halt