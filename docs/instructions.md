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
- `func`: the ALU func

Depending on the `fmt` code, the opcode is constructed from different patterns of bits denoted in the decoding logic section. There are 32 opcodes.


### Instruction Decoding Details

This diagram shows how the opcodes are decoded into each instruction class.

![Instruction Decoding Details](isa_decoding.png)

On the left is the various formats, then the `fmt` code for each instruction class, then (optionally) the zeroth bit of op if it's required, then whether the instruction has an immediate. Then you get the 5 opcode bits.

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
  - `if.gt A, B`    = `if.lt B, A`
  - `if.gt A, imm`  = `if.ge A, imm+1`
  - `if.hi A, B`    = `if.lo B, A`
  - `if.hi A, imm`  = `if.hs B, imm+1`
  - `if.le A, B`    = `if.ge B, A`
  - `if.le A, imm`  = `if.lt A, imm-1`
  - `if.ls A, B`    = `if.hi B, A`
  - `if.ls A, imm`  = `if.lo B, imm-1`
- `ri6` & `rr`:
  - `if.ne` (Z==0) - not equal
  - `if.eq` (Z==1) - equal
  - `if.ge` (N==V) - signed greater than or equal `>=`
  - `if.lt` (N!=V) - signed less than `<`
  - `if.hs` (C==1) - higher same / unsigned greater or equal `>=`
  - `if.lo` (C==0) - lower / unsigned less than `<`

### CSRs

- `r` only:
  - `rcsr`, `wcsr` - read and write CSRs (Computer Status Register)

### Branches

- `i11`/`r`:
  - `jump` - jump to a different program address
  - `jal` - jump and link, that is, call a function

### Sys Ops

- `i11`:
  - `imm` - prefix to extend next instruction's immediate
  - `sys` - syscall?
- `rr`:
  - `nop` - no operation
  - `rets` - return from/to system
  - `error` - halt with error
  - `halt` - halt with success

## Instruction Details

This is quite incomplete... a few instructions are done as an example. The details may be inacurate due to the instruction set design not being finished yet.

### system

    rr format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |    N/A    |    N/A    | -| 0  0  0|  op | 0  0|

    nop - no operation
      format:    rr
      assembler: nop
      example:   nop
      symbolic:
      operation:
                 Do nothing.

    break - breakpoint
      format:    rr
      assembler: break
      example:   break
      symbolic:
      operation:
                 Break to debugger, or do nothing if not
                 debugging.

    error - stop unit test with error
      format:    rr
      assembler: error
      example:   error
      symbolic:  error
      operation:
                 Asserts the error line. Does nothing
                 other than that unless in a unit test,
                 in which case it exits with a failure.

    halt - stop unit test with success
      format:    rr
      assembler: halt
      example:   halt
      symbolic:  halt
      operation:
                 Asserts the halt line. Does nothing
                 other than that unless in a unit test,
                 in which case it exits with success.
                 May halt the clock if in system mode in
                 the future.

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

    ri8 format:

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
      format:    i11; rr
      assembler: jump rd
      example:   jump r3
      symbolic:  pc <- rd
      operation:
                 Set the program counter to the absolute
                 value stored in the register. This is
                 particularly useful to return after a
                 call.

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
                 `rs` and is stored in register `rd`.

    storeb - store byte
      format:    ls
      assembler: storeb rd, [rs, imm4]
      example:   storeb r5, [r2, 15]
      symbolic:  memb[rs + imm4] <- rd & 0xff
      operation:
                 An 8-bit byte in the lower byte of register
                 `rd` is stored to memory at the absolute
                 address provided by the register `rs`.
