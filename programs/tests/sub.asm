#include "../cpudef.asm"

move r1, 5
move r2, 3
sub r1, r2
if.ne r1, 2
  error

sub r1, 2
if.ne r1, 0
  error
halt

