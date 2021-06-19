#include "../cpudef.asm"

call test
if.ne r1, 5
  error
halt

error
test:
  move r1, 5
  jump r0

error
halt
