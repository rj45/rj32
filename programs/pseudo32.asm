#ruledef {
  ; pseudoinstructions
  nop                          => asm { add zero, zero, zero }
  mv  {rd:reg}, {rs:reg}       => asm { add {rd}, zero, {rs} }
  not {rd:reg}, {rs:reg}       => asm { xor {rd}, {rs}, -1 }
  neg {rd:reg}, {rs:reg}       => asm { sub {rd}, zero, {rs} }

  beqz {rd:reg}, {imm}         => asm { beq {rd}, zero, imm }

  li {rd:reg}, {imm: s13} => asm {
    addi {rd}, zero, imm
  }

  li {rd:reg}, {imm: u32} => asm {
    lui {rd}, imm`32
    addi {rd}, {rd}, imm[11:0]
  }

}


; code bank is the main program memory bank
#bankdef code
{
  #bits 8
  #addr 0x80000000
  #size 0x10000
  #outp 0
}

; data is the bank where strings, constants and pre-initialized
; values goes.
#bankdef data
{
  #bits 8
  #addr 0x80010000
  #size 0x10000
  #outp 0x10000*16
}

; bss is the main data memory bank for uninitialized variables
#bankdef bss
{
  #bits 8
  #addr 0x80020000
  #size 0x80000
}

#bank code
