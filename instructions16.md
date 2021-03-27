# Instructions

The instruction set of the rj32 processor, so named for the 32 core instructions it implements.

## Summary

### Memory Ops

  ls:
    - lw, lb, sw, sb

### ALU Ops

  ri & rr:
    - add, sub, ltu, gtu, and, or, xor, move
    - addc, subc, lt, gt, shl, shr, asr, eq
  r only:
    - nop, neg, swap, sxt

### CSRs

  r only:
    - rcsr, wcsr, scsr, ccsr

### Branches

  i/r:
    - jump, jal, brt, brf

### Sys Ops

  i:
    - imm, break, halt, rets

## Instruction Details

### load / store

    ls format:

    |15|14|13|12|11|10| 9| 8| 7| 6| 5| 4| 3| 2| 1| 0|
    |     imm5     |     rs    |     rd    | 1| op0 |

    load - load word
      format:    ls
      assembler: load rd, [rs, imm5]
      example:   load r5, [r2, 15]
      symbolic:  rd <- mem[rs + imm5 * 2]
      operation: A 16-bit word is loaded from memory at the
                 offset provided by the doubled, zero extended
                 immediate plus a base register `rs` and is
                 stored in register `rd`.

    store - store word
      format:    ls
      assembler: store rd, [rs, imm5]
      example:   store r5, [r2, 15]
      symbolic:  mem[rs + imm5 * 2] <- rd
      operation: A 16-bit word is stored in memory at the
                 offset provided by the doubled, zero extended
                 immediate plus a base register `rs` from the register `rd`.

    loadb - load byte
      format:    ls
      assembler: load rd, [rs, imm5]
      example:   load r5, [r2, 15]
      symbolic:  rd <- mem[rs + imm5]
      operation: An 8-bit byte is loaded from memory at the
                 offset provided by the zero extended
                 immediate plus a base register `rs` and is
                 stored in register `rd`.

    storeb - store byte
      format:    ls
      assembler: store rd, [rs, imm5]
      example:   store r5, [r2, 15]
      symbolic:  mem[rs + imm5] <- rd
      operation: An 8-bit byte is stored in memory at the
                 offset provided by the zero extended immediate plus a base register `rs` from the register `rd`.


