; -------------------------------------------------------------
;  Arithmetic tests
; -------------------------------------------------------------

test_2:
  li gp, 2
  li x1, 0x00000000
  addi x14, x1, ((0x000) | (-(((0x000) >> 11) & 1) << 11))

  li x7, 0x00000000
  bne x14, x7, fail


test_3:
  li gp, 3
  li x1, 0x00000001
  addi x14, x1, ((0x001) | (-(((0x001) >> 11) & 1) << 11))

  li x7, 0x00000002
  bne x14, x7, fail


test_4:
  li gp, 4
  li x1, 0x00000003
  addi x14, x1, ((0x007) | (-(((0x007) >> 11) & 1) << 11))

  li x7, 0x0000000a
  bne x14, x7, fail



test_5:
  li gp, 5
  li x1, 0x00000000
  addi x14, x1, ((0x800) | (-(((0x800) >> 11) & 1) << 11))

  li x7, 0xfffff800
  bne x14, x7, fail


test_6:
  li gp, 6
  li x1, 0x80000000
  addi x14, x1, ((0x000) | (-(((0x000) >> 11) & 1) << 11))

  li x7, 0x80000000
  bne x14, x7, fail


test_7:
  li gp, 7
  li x1, 0x80000000
  addi x14, x1, ((0x800) | (-(((0x800) >> 11) & 1) << 11))

  li x7, 0x7ffff800
  bne x14, x7, fail



test_8:
  li gp, 8
  li x1, 0x00000000
  addi x14, x1, ((0x7ff) | (-(((0x7ff) >> 11) & 1) << 11))

  li x7, 0x000007ff
  bne x14, x7, fail


test_9:
  li gp, 9
  li x1, 0x7fffffff
  addi x14, x1, ((0x000) | (-(((0x000) >> 11) & 1) << 11))

  li x7, 0x7fffffff
  bne x14, x7, fail


test_10:
  li gp, 10
  li x1, 0x7fffffff
  addi x14, x1, ((0x7ff) | (-(((0x7ff) >> 11) & 1) << 11))

  li x7, 0x800007fe
  bne x14, x7, fail



test_11:
  li gp, 11
  li x1, 0x80000000
  addi x14, x1, ((0x7ff) | (-(((0x7ff) >> 11) & 1) << 11))

  li x7, 0x800007ff
  bne x14, x7, fail


test_12:
  li gp, 12
  li x1, 0x7fffffff
  addi x14, x1, ((0x800) | (-(((0x800) >> 11) & 1) << 11))

  li x7, 0x7ffff7ff
  bne x14, x7, fail



test_13:
  li gp, 13
  li x1, 0x00000000
  addi x14, x1, ((0xfff) | (-(((0xfff) >> 11) & 1) << 11))

  li x7, 0xffffffff
  bne x14, x7, fail


test_14:
  li gp, 14
  li x1, 0xffffffff
  addi x14, x1, ((0x001) | (-(((0x001) >> 11) & 1) << 11))

  li x7, 0x00000000
  bne x14, x7, fail


test_15:
  li gp, 15
  li x1, 0xffffffff
  addi x14, x1, ((0xfff) | (-(((0xfff) >> 11) & 1) << 11))

  li x7, 0xfffffffe
  bne x14, x7, fail



test_16:
  li gp, 16
  li x1, 0x7fffffff
  addi x14, x1, ((0x001) | (-(((0x001) >> 11) & 1) << 11))

  li x7, 0x80000000
  bne x14, x7, fail



; -------------------------------------------------------------
;  Source/Destination tests
; -------------------------------------------------------------

test_17:
  li gp, 17
  li x1, 13
  addi x1, x1, ((11) | (-(((11) >> 11) & 1) << 11))

  li x7, 24
  bne x1, x7, fail



; -------------------------------------------------------------
;  Bypassing tests
; -------------------------------------------------------------

test_18:
  li gp, 18
  li x4, 0
.L1:
  li x1, 13
  addi x14, x1, ((11) | (-(((11) >> 11) & 1) << 11))
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 24
  bne x6, x7, fail


test_19:
  li gp, 19
  li x4, 0
.L1:
  li x1, 13
  addi x14, x1, ((10) | (-(((10) >> 11) & 1) << 11))
  nop
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 23
  bne x6, x7, fail


test_20:
  li gp, 20
  li x4, 0
.L1:
  li x1, 13
  addi x14, x1, ((9) | (-(((9) >> 11) & 1) << 11))
  nop
  nop
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 22
  bne x6, x7, fail



test_21:
  li gp, 21
  li x4, 0
.L1:
  li x1, 13
  addi x14, x1, ((11) | (-(((11) >> 11) & 1) << 11))
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 24
  bne x14, x7, fail


test_22:
  li gp, 22
  li x4, 0
.L1:
  li x1, 13
  nop
  addi x14, x1, ((10) | (-(((10) >> 11) & 1) << 11))
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 23
  bne x14, x7, fail


test_23:
  li gp, 23
  li x4, 0
.L1:
  li x1, 13
  nop
  nop
  addi x14, x1, ((9) | (-(((9) >> 11) & 1) << 11))
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 22
  bne x14, x7, fail



test_24:
  li gp, 24
  addi x1, x0, ((32) | (-(((32) >> 11) & 1) << 11))

  li x7, 32
  bne x1, x7, fail


test_25:
  li gp, 25
  li x1, 33
  addi x0, x1, ((50) | (-(((50) >> 11) & 1) << 11))

  li x7, 0
  bne x0, x7, fail



  bne x0, gp, pass
fail:
  fence
.L1:
  beqz gp, .L1
  sll gp, gp, 1
  or gp, gp, 1
  li a7, 93
  addi a0, gp, 0
  ecall
pass:
  fence
  li gp, 1
  li a7, 93
  li a0, 0
  ecall
