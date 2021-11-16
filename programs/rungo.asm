; run go's main__main function
move sp, 0xFEFF
move gp, 0

; initialize all the global variables
call main__init

; check that the stack is not corrupted
if.ne sp, 0xFEFF
  error

call main__main

if.ne r1, 0
  error

; check that the stack is not corrupted
if.ne sp, 0xFEFF
  error

halt