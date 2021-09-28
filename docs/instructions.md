# Instructions

The instruction set of the rj32 processor, so named for the 32 core instructions it implements.

## Instruction Summary

Some definitions:

- `rd`: destination register and often also left source
- `rs`: right source register
- `immediate`: a number provided in the instruction itself
- `imm4`: unsigned 4 bit immediate
- `imm6`, `imm8`, `imm11`, `imm12`: signed 6, 8, 11 and 12 bit immediates
- `csr`: computer status register - special registers used to control processor state

|  op | asm                    | description              |
| --: | ---------------------- | ------------------------ |
|   0 | `nop`                  | no operation             |
|   1 | `rets`                 | return to/from system    |
|   2 | `error`                | halt with error          |
|   3 | `halt`                 | halt without error       |
|   4 | `rcsr rd, csr`         | read csr                 |
|   5 | `wcsr csr, rd`         | write csr                |
|   6 | `move rd, rs/imm8`     | move into register       |
|   7 | `loadc rd, rs/imm8`    | load constant            |
|   8 | `jump rd/imm11`        | set program counter      |
|   9 | `imm imm12`            | extend immediate         |
|  10 | `call imm12`           | jump and save `PC`       |
|  11 | _reserved_             | also `imm` instruction   |
|  12 | `load rd, [rs, imm4]`  | load word                |
|  13 | `store [rs,imm4], rd`  | store word               |
|  14 | `loadb rd, [rs, imm4]` | load byte                |
|  15 | `storeb [rs,imm4], rd` | store byte               |
|  16 | `add rd, rs/imm6`      | add                      |
|  17 | `sub rd, rs/imm6`      | subtract                 |
|  18 | `addc rd, rs/imm6`     | add with carry           |
|  19 | `subc rd, rs/imm6`     | subtract with carry      |
|  20 | `xor rd, rs/imm6`      | exclusive or             |
|  21 | `and rd, rs/imm6`      | logical and              |
|  22 | `or rd, rs/imm6`       | logical or               |
|  23 | `shl rd, rs/imm6`      | logical shift left       |
|  24 | `shr rd, rs/imm6`      | logical shift right      |
|  25 | `asr rd, rs/imm6`      | arithmetic shift right   |
|  26 | `if.eq rd, rs/imm6`    | if equal                 |
|  27 | `if.ne rd, rs/imm6`    | if not equal             |
|  28 | `if.lt rd, rs/imm6`    | if less than             |
|  29 | `if.ge rd, rs/imm6`    | if greater or equal      |
|  30 | `if.ult rd, rs/imm6`   | if unsigned less than    |
|  31 | `if.uge rd, rs/imm6`   | if unsigned greater than |

There are also the following pseudoinstructions:

| asm              | implementation          |
| ---------------- | ----------------------- |
| `return`         | `jump r0`               |
| `sxt rd`         | `shl rd, 8; asr rd, 8`  |
| `not rd`         | `xor rd, -1`            |
| `neg rd`         | `xor rd, -1; add rd, 1` |
| `if.gt rd, rs`   | `if.lt rs, rd`          |
| `if.gt rd, imm`  | `if.ge rd, imm+1`       |
| `if.le rd, rs`   | `if.ge rs, rd`          |
| `if.le rd, imm`  | `if.lt rd, imm-1`       |
| `if.ugt rd, rs`  | `if.ult rs, rd`         |
| `if.ugt rd, imm` | `if.uge rs, imm+1`      |
| `if.ule rd, rs`  | `if.ugt rs, rd`         |
| `if.ule rd, imm` | `if.ult rs, imm-1`      |

### Registers

There are 16 registers, only r0 is hard coded with a special function, the other registers can be used in any way. Their calling convention purpose is denoted.

| reg   | alias | conventional usage    |
| ----- | ----- | ------------------------ |
| `r0`  | `ra`  | return address           |
| `r1`  | `a0`  | return value / 1st arg   |
| `r2`  | `a1`  | second function argument |
| `r3`  | `s0`  | callee saved reg  |
| `r4`  | `s1`  | callee saved reg         |
| `r5`  | `s2`  | callee saved reg         |
| `r6`  | `s3`  | callee saved reg         |
| `r7`  | `s4`  | callee saved reg         |
| `r8`  | `t0`  | caller saved temp reg    |
| `r9`  | `t1`  | caller saved temp reg    |
| `r10` | `t2`  | caller saved temp reg    |
| `r11` | `t3`  | caller saved temp reg    |
| `r12` | `t4`  | caller saved temp reg    |
| `r13` | `t5`  | caller saved temp reg    |
| `r14` | `bp`  | data base pointer        |
| `r15` | `sp`  | stack pointer            |

Callee saved registers should be saved to the stack if used in a function.

Caller saved registers are expected to be clobbered in a function call, so if they are used, they are saved in the stack before a function is called.

The first two function arguments are provided in a0 and a1, and the return value is provided in a1. If there are more arguments than two or an argument doesn't fit in 16 bits, they are provided on the stack.

## Instruction Encodings

### Instruction Encoding Overview

![Instruction Encoding Overview](isa_encodings.png)

Bits 0-1 are a 2 bit format (`fmt`) code
Bits 2-6 are the `op` code field
Bit 7 on many instructions is reserved for future expansion
The `func` code is the ALU function, which is part of the opcode.

Depending on the `fmt` code, the opcode is constructed from different patterns of bits denoted in the decoding logic section. There are 32 opcodes.

### Instruction Decoding Details

This diagram shows how the opcodes are decoded into each instruction class.

![Instruction Decoding Details](isa_decoding.png)

On the left is the various formats, then the `fmt` code for each instruction class, then (optionally) the zeroth bit of op if it's required, then whether the instruction has an immediate. Then you get the 5 opcode bits. Grey bits are added during decoding, blue or purple bits come from the instruction itself.

### Prefix Instructions and State

There are a few prefix instructions (`imm`, `addc` and `subc`) that carry state over into the next instruction.

If an interrupt would happen between these instructions and the ones they modify, state would be lost. To avoid that, interrupts are disabled until the instruction they modify completes and that state is fully transferred.

The state is reset after the instruction modified uses that state to ensure it cannot cause bugs by putting the processor into an unexpected state.

#### Imm Prefix

The `imm` prefix instruction extends the immediate of the next instruction. The 12 bits supplied become the 12 most significant bits of the next instruction, and the modified instruction supplies the remaining 4 least significant bits.

#### Addc and Subc Prefixes

It's useful to be able to be able to add or subtract numbers larger than 16 bits. In order to do that, the carry out of the previous `add` or `sub` needs to be fed into the carry in of the next `add` or `sub`. The `addc` and `subc` instructions are provided to set a carry flag as well as a "use carry" flag that will tell the next instruction to use the provided carry. Otherwise these instructions act the same as `add` and `sub`.

For example:

```asm
  ; add two 64 bit numbers
  addc r1, r5
  addc r2, r6
  addc r3, r7
  add  r4, r8  ; this instruction will use the carry

  add  r9, r10 ; this instruction won't use the carry
```

### Skipping

Instead of conditional branches, there are `if.cc` instructions which will skip the next instruction if the condition is not true. Think of these acting like `if` statements in higher level languages.

- `if.eq` - equal
- `if.ne` - not equal
- `if.lt` - less than
- `if.ge` - greater or equal
- `if.ult` - unsigned less than
- `if.uge` - unsigned greater or equal

There are also psuedoinstructions for the following conditions:

- `if.gt` - greater than
- `if.le` - less or equal
- `if.ugt` - unsigned greater than
- `if.ule` - unsigned less or equal

Any prefix instructions (`imm`, `addc` or `subc`) will be skipped in addition to one regular instruction. During a skip, interrupts are disabled. Each skipped instruction uses one clock cycle (they act like `nop`).

## Instruction Details

This is quite incomplete... a few instructions are done as an example. The details may be inaccurate due to the instruction set design not being finished yet.

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

### Arithmetic

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
                 or subc instruction preceded this
                 one, the carry flag is also added.

    add - add register
      format:    rr
      assembler: add rd, rs
      example:   add r13, r2
      symbolic:  rd <- rd + rs + C
      operation:
                 Add registers `rd` and `rs` and store
                 back into register `rd`. If an addc
                 or subc instruction preceded this
                 one, the carry flag is also added.

    addc - add immediate with carry
      format:    ri6
      assembler: add rd, imm6
      example:   add r3, 15
      symbolic:  rd <- rd + imm6 + C; C <- carry
      operation:
                 Add a 6 bit signed immediate to the
                 destination register rd. If an addc
                 or subc instruction preceded this
                 one, the carry flag is also added.
                 The carry flag is set by this
                 instruction. This instruction is not
                 interruptible, and if a skip is in
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
                 or subc instruction preceded this
                 one, the carry flag is also added.
                 The carry flag is set by this
                 instruction and will be cleared after
                 the next instruction. This instruction
                 is not interruptible, and if a skip is in
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
                 preceded this one, the carry flag
                 is also used as a borrow flag.
                 The carry flag is set by this
                 instruction and will be cleared after
                 the next instruction. This instruction
                 is not interruptible, and if a skip is in
                 progress, it will also skip the next
                 instruction.

    subc - subtract register with carry
      format:    rr
      assembler: add rd, rs
      example:   add r13, r2
      symbolic:  rd <- rd - rs - C; C <- carry
      operation:
                 Subtract registers `rd` and `rs` and store
                 back into register `rd`. If an addc
                 or subc instruction preceded this
                 one, the carry flag is also subtracted.
                 The carry flag is set by this
                 instruction and will be reset after the
                 next instruction. This instruction is
                 not interruptible, and if a skip is in
                 progress, it will also skip the next
                 instruction.
