
; kernel bank is the code for the bootloader/kernel
#bankdef kernel
{
  #bits 8
  #addr 0x0000
  #size 8192
  #outp 0
}

; code bank is the main program memory bank
#bankdef code
{
  #bits 8
  #addr 8192
  #size 0x10000-8192
  #outp 8192*8
}

; data is the bank where strings, constants and pre-initialized
; values goes.
#bankdef data
{
  #bits 8
  #addr 0x10000
  #size 0x10000
  #outp 0x10000*8
}

; bss is the main data memory bank for uninitialized variables
; this bank is not stored in the output file, and in hardware
; it should be zeroed out on initialization
#bankdef bss
{
  #bits 8
  #addr 0x20000
  #size 0x10000
}

#bank kernel


stackStartAddress = 0x01FFFFFC

; run go's main__main function
init:
  ; initialize the stack
  LD   sp, stackStartAddress

  ; initialize all the global variables
  CALL  main__init

  ; check that the stack is not corrupted
  CMP   sp, stackStartAddress
  BR.EQ .stackok
  ERR
  BRA   .looperr
.stackok:

  ; run the main program
  CALL   main__main

  ; check that the stack is not corrupted
  CMP    sp, stackStartAddress
  BR.EQ  .stackok2
  ERR
  BRA    .looperr
.stackok2:

  ; halt or loop forever
  HLT

.loophalt:
  BRA    .loophalt

.looperr:
  ; if the ERR instruction does nothing, then the fact
  ; it's looping here can be used to see if stack
  ; corruption happened
  BRA    .looperr

#bank code
