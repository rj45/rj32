#include "../cpudef.asm"

move r1, 5
move r2, 3
add r1, r2
add r1, 2
if.ne r1, 10
error
halt

