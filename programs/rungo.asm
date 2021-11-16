; run go's main__main function
move sp, 0xFEFF
move gp, 0

call main__main

if.ne r1, 0
  error

; check if stack not properly handled
if.ne sp, 0xFEFF
  error

halt