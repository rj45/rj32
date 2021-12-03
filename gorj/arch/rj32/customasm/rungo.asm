stackStartAddress = 0xFEFF

; run go's main__main function

; initialize the stack and global pointer
move sp, stackStartAddress
move gp, 0

; initialize all the global variables
call main__init

; check that the stack is not corrupted
if.ne sp, stackStartAddress
  error

; run the main program
call main__main

; check that the stack is not corrupted
if.ne sp, stackStartAddress
  error

halt

