; run go's main__main function
move sp, 0
sub sp, 1
move bp, 0

call main__main

if.ne r1, 0
  error
sub r1, 1

; stack not properly handled
if.ne sp, r1
  error

halt