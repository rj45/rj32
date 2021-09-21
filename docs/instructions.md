# Instructions

The instruction set of the rj32 processor, so named for the 32 core instructions it implements.

## Instruction Encodings

### Instruction Encoding Overview

![Instruction Encoding Overview](isa_encodings.png)

Bits 0-1 are a 2 bit format (`fmt`) code
Bits 2-6 are the `op` code field
Bit 7 on many instructions is reserved for future expansion

- `rd`: the left operand register and destination register
- `rs`: the right operand register
- `imm4`-`imm12`: various widths of immediate values
  - `imm4` in the `ls` format is unsigned
  - All other immediates are sign extended to 16 bits
  - The `imm` instruction can be used to extend the next instruction's immediate to 16 bits
- `func`: the ALU func

Depending on the `fmt` code, the opcode is constructed from different patterns of bits denoted in the decoding logic section. There are 32 opcodes.

### Instruction Decoding Details

This diagram shows how the opcodes are decoded into each instruction class.

![Instruction Decoding Details](isa_decoding.png)

On the left is the various formats, then the `fmt` code for each instruction class, then (optionally) the zeroth bit of op if it's required, then whether the instruction has an immediate. Then you get the 5 opcode bits.

### Prefix Instructions and State

There are a few prefix instructions (`imm`, `addc` and `subc`) that carry state over into the next instruction.

If an interrupt would happen between these instructions and the ones they modify, state would be lost.

To avoid that, interrupts are disabled until the instruction they modify completes and that state is fully transferred. Then the state is reset to ensure it cannot affect anything while interrupts are enabled.

Also, since these instructions might be automatically inserted by the assembler, the same flag used to defer interrupts is used to extend a skip as well. So skip instructions will skip over any prefix instructions as well as the instruction they modify.

## Instruction Summary

A summary of the instructions supported, categorized by their instruction format.

### Memory Ops

- `ls`:
  - `load`, `store` - load/store word in data memory
  - `loadb`, `storeb` - load/store byte in data memory

### ALU Ops

- psuedoinstructions:
  - `not A` = `xor A, -1`
  - `neg A` = `xor A, -1; add A, 1`
- `ri8` & `rr`:
  - `move` - move value from register/immediate to register
- `ri6` & `rr`:
  - `add` - add
  - `sub` - subtract
  - `add` - add & save carry
  - `sub` - subtract & save carry for next instruction
  - `and`, `or`, `xor` - bitwise ops
  - `shl`, `shr`, `asr` - shifts

#### If Skip Ops

These instructions check a condition and skip the next instruction if the condition is false.

- psuedoinstructions:
  - `if.gt A, B`     = `if.lt B, A`
  - `if.gt A, imm`   = `if.ge A, imm+1`
  - `if.le A, B`     = `if.ge B, A`
  - `if.le A, imm`   = `if.lt A, imm-1`
  - `if.ugt A, B`    = `if.ult B, A`
  - `if.ugt A, imm`  = `if.uge B, imm+1`
  - `if.ule A, B`    = `if.ugt B, A`
  - `if.ule A, imm`  = `if.ult B, imm-1`
- `ri6` & `rr`:
  - `if.ne`  (Z==0) - not equal
  - `if.eq`  (Z==1) - equal
  - `if.ge`  (N==V) - signed greater than or equal `>=`
  - `if.lt`  (N!=V) - signed less than `<`
  - `if.uge` (C==1) - unsigned greater or equal `>=`
  - `if.ult` (C==0) - unsigned less than `<`

### CSRs

- `rr` only:
  - `rcsr`, `wcsr` - read and write CSRs (Computer Status Register)

### Branches

- `i11`/`r`:
  - `jump` - jump to a different program address
  - `jal` - jump and link, that is, call a function

### Sys Ops

- `i11`:
  - `imm` - prefix to extend next instruction's immediate
- `rr`:
  - `nop` - no operation
  - `rets` - return from/to system
  - `error` - halt with error
  - `halt` - halt with success

## Instruction Details

This is quite incomplete... a few instructions are done as an example. The details may be inacurate due to the instruction set design not being finished yet.

### system

    i12 format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |              imm12                | 1  1  0  1|

    rr format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |    N/A    |    N/A    | -| 0  0  0|  op | 0  0|

    imm - extended immediate prefix
      format:    i12
      assembler: imm imm12
      example:   imm 0x1234
      symbolic:  imm <- imm12 & 0xfff0
      operation:
                 This instruction extends the immediate of
                 the next instruction with the upper 12 bits
                 of the 16 bit value. The lower 4 bits come
                 from the next instruction.

                 This instruction is automatically inserted
                 by the assembler when necessary. If a skip
                 is active, both the imm instruction and the
                 following instruction will be skipped.

                 Interrupts are disabled during this
                 instruction since the state of the imm
                 register cannot be saved.

    nop - no operation
      format:    rr
      assembler: nop
      example:   nop
      symbolic:
      operation:
                 Do nothing.

    error - stop unit test with error
      format:    rr
      assembler: error
      example:   error
      symbolic:  error
      operation:
                 Assert the error line and spin until reset.
                 Useful in unit tests to indicate a failure.

    halt - stop unit test with success
      format:    rr
      assembler: halt
      example:   halt
      symbolic:  halt
      operation:
                 Assert the halt line and spin until reset.
                 Useful in unit tests to indicate success.
                 Exits the emulator.

### move

    ri8 format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     rd    |     rs    |    imm4   |op| 0  1  0|

    rr format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     rd    |     rs    | -| 0  0  1  1|op| 0  0|

    move - move word
      format:    ri8
      assembler: move rd, imm8
      example:   move r5, 121
      symbolic:  rd <- imm8
      operation:
                 Copy a value from an 8 bit signed
                 immediate into the destination
                 register rd.

    move - move word
      format:    rr
      assembler: move rd, rs
      example:   move r3, r9
      symbolic:  rd <- rs
      operation:
                 Copy a value from another register rs
                 into the destination register rd.

### branches

    imm11 format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |               imm11            |  op | 1  0  1|

    rr format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     rd    |    N/A    | -| 0  1  0|  op | 0  0|

    jump - set program counter
      format:    i11
      assembler: jump imm11
      example:   jump label
      symbolic:  pc <- pc + imm11
      operation:
                 Increment the program counter by the
                 given offset, effectively jumping to
                 the relative address given.

    jump - set program counter
      format:    rr
      assembler: jump rd
      example:   jump r3
      symbolic:  pc <- rd
      operation:
                 Set the program counter to the absolute
                 value stored in the register. This is
                 particularly useful to return after a
                 call.

    call - save program counter and set it
      format:    i11
      assembler: call imm11
      example:   call label
      symbolic:  r0 <- pc; pc <- pc + imm11
      operation:
                 Increment the program counter by the
                 given offset, effectively calling a
                 function at the relative address given
                 while simultaneously saving the old
                 program counter in r0 (ra). The called
                 function can optionally save r0 to the
                 stack.

    call - save program counter and set it
      format:    rr
      assembler: call rd
      example:   call r3
      symbolic:  pc <- rd
      operation:
                 Set the program counter to the absolute
                 value stored in the register while also
                 saving the program counter in r0 (ra).
                 This can be useful for virtual function
                 calls and calling function pointers.

### load / store

    ls format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     rd    |     rs    |    imm4   |  op | 1  0|

    load - load word
      format:    ls
      assembler: load rd, [rs, imm4]
      example:   load r5, [r2, 15]
      symbolic:  rd <- memw[(rs & ~1) + imm4 * 2]
      operation:
                 A 16-bit word is loaded from memory at the
                 offset provided by the doubled, zero
                 extended immediate plus a base register
                 `rs` and is stored in register `rd`. The
                 least significant bit of `rs` is ignored.

    store - store word
      format:    ls
      assembler: store rd, [rs, imm4]
      example:   store r5, [r2, 15]
      symbolic:  memw[(rs & ~1) + imm4 * 2] <- rd
      operation:
                 A 16-bit word is stored in memory at the
                 offset provided by the doubled, zero
                 extended immediate plus a base register `rs`
                 from the register `rd`.The least significant
                 bit of `rs` is ignored.

    loadb - load byte
      format:    ls
      assembler: loadb rd, [rs, imm4]
      example:   loadb r5, [r2, 15]
      symbolic:  rd <- sext(memb[rs + imm4])
      operation:
                 An 8-bit byte is loaded from memory at the
                 absolute address provided by the register
                 `rs` and is stored in register `rd`. The
                 loaded byte is sign extended to 16 bits.
                 To undo the sign extension, and with
                 0xff00.

    storeb - store byte
      format:    ls
      assembler: storeb rd, [rs, imm4]
      example:   storeb r5, [r2, 15]
      symbolic:  memb[rs + imm4] <- rd & 0xff
      operation:
                 An 8-bit byte in the lower byte of register
                 `rd` is stored to memory at the absolute
                 address provided by the register `rs`.

### arithmatic

    ri6 format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     rd    |       imm6      |     op    | 1  1|

    rr format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     rd    |     rs    | -| 1|     op    | 0  0|


    add - add immediate
      format:    ri6
      assembler: add rd, imm6
      example:   add r3, 15
      symbolic:  rd <- rd + imm6 + C
      operation:
                 Add a 6 bit signed immediate to the
                 destination register rd. If an addc
                 or subc instruction preceeded this
                 one, the carry flag is also added.

    add - add register
      format:    rr
      assembler: add rd, rs
      example:   add r13, r2
      symbolic:  rd <- rd + rs + C
      operation:
                 Add registers `rd` and `rs` and store
                 back into register `rd`. If an addc
                 or subc instruction preceeded this
                 one, the carry flag is also added.

    addc - add immediate with carry
      format:    ri6
      assembler: add rd, imm6
      example:   add r3, 15
      symbolic:  rd <- rd + imm6 + C; C <- carry
      operation:
                 Add a 6 bit signed immediate to the
                 destination register rd. If an addc
                 or subc instruction preceeded this
                 one, the carry flag is also added.
                 The carry flag is set by this
                 instruction. This instruction is not
                 interruptable, and if a skip is in
                 progress, it will also skip the next
                 instruction.

    addc - add register with carry
      format:    rr
      assembler: add rd, rs
      example:   add r13, r2
      symbolic:  rd <- rd + rs + C; C <- carry
      operation:
                 Add registers `rd` and `rs` and store
                 back into register `rd`. If an addc
                 or subc instruction preceeded this
                 one, the carry flag is also added.
                 The carry flag is set by this
                 instruction and will be cleared after
                 the next instruction. This instruction
                 is not interruptable, and if a skip is in
                 progress, it will also skip the next
                 instruction.

    sub - subtract immediate
      format:    ri6
      assembler: sub rd, imm6
      example:   sub r3, 43
      symbolic:  rd <- rd - imm6 - C
      operation:
                 Subtract a 6 bit signed immediate from the
                 destination register rd. If the previous
                 instruction was an `addc` or `subc` the
                 carry flag will also be subtracted.

    sub - subtract register
      format:    rr
      assembler: sub rd, rs
      example:   sub r1, r5
      symbolic:  rd <- rd - rs - C
      operation:
                 Subtract registers `rd` and `rs` and store
                 back into register `rd`. If the previous
                 instruction was an `addc` or `subc` the
                 carry flag will also be subtracted.

    subc - Subtract immediate with carry
      format:    ri6
      assembler: sub rd, imm6
      example:   sub r3, 15
      symbolic:  rd <- rd - imm6 - C; C <- carry
      operation:
                 Subtract a 6 bit signed immediate
                 from the destination register rd.
                 If an addc or subc instruction
                 preceeded this one, the carry flag
                 is also used as a borrow flag.
                 The carry flag is set by this
                 instruction and will be cleared after
                 the next instruction. This instruction
                 is not interruptable, and if a skip is in
                 progress, it will also skip the next
                 instruction.

    subc - subtract register with carry
      format:    rr
      assembler: add rd, rs
      example:   add r13, r2
      symbolic:  rd <- rd - rs - C; C <- carry
      operation:
                 Subtrat registers `rd` and `rs` and store
                 back into register `rd`. If an addc
                 or subc instruction preceeded this
                 one, the carry flag is also subtracted.
                 The carry flag is set by this
                 instruction and will be reset after the
                 next instruction. This instruction is
                 not interruptable, and if a skip is in
                 progress, it will also skip the next
                 instruction.
