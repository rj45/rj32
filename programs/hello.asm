#include "cpudef.asm"

; calculate telnet base address
move r5, 1
shl r5, 14

move r1, "H"
store [r5,0], r1
move r1, "e"
store [r5,0], r1
move r1, "l"
store [r5,0], r1
move r1, "l"
store [r5,0], r1
move r1, "o"
store [r5,0], r1
move r1, ","
store [r5,0], r1
move r1, " "
store [r5,0], r1
move r1, "w"
store [r5,0], r1
move r1, "o"
store [r5,0], r1
move r1, "r"
store [r5,0], r1
move r1, "l"
store [r5,0], r1
move r1, "d"
store [r5,0], r1
move r1, "!"
store [r5,0], r1
move r1, "\r"
store [r5,0], r1
move r1, "\n"
store [r5,0], r1

echo:
  nop
  load r1, [r5, 0]
  store [r5,0], r1
  nop
  nop
  jump echo

