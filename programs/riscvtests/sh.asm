  ; -------------------------------------------------------------
  ;  Basic tests
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  la x1, tdat
  li x2, 0x000000aa
  la x15, .L7:
  sh x2, 0(x1)
  lh x14, 0(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x000000aa
  bne x14, x7, fail
  
  
test_3:
  li gp, 3
  la x1, tdat
  li x2, 0xffffaa00
  la x15, .L7:
  sh x2, 2(x1)
  lh x14, 2(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffaa00
  bne x14, x7, fail
  
  
test_4:
  li gp, 4
  la x1, tdat
  li x2, 0xbeef0aa0
  la x15, .L7:
  sh x2, 4(x1)
  lw x14, 4(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xbeef0aa0
  bne x14, x7, fail
  
  
test_5:
  li gp, 5
  la x1, tdat
  li x2, 0xffffa00a
  la x15, .L7:
  sh x2, 6(x1)
  lh x14, 6(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffa00a
  bne x14, x7, fail
  
  

  ;  Test with negative offset

test_6:
  li gp, 6
  la x1, tdat8
  li x2, 0x000000aa
  la x15, .L7:
  sh x2, -6(x1)
  lh x14, -6(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x000000aa
  bne x14, x7, fail
  
  
test_7:
  li gp, 7
  la x1, tdat8
  li x2, 0xffffaa00
  la x15, .L7:
  sh x2, -4(x1)
  lh x14, -4(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffaa00
  bne x14, x7, fail
  
  
test_8:
  li gp, 8
  la x1, tdat8
  li x2, 0x00000aa0
  la x15, .L7:
  sh x2, -2(x1)
  lh x14, -2(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x00000aa0
  bne x14, x7, fail
  
  
test_9:
  li gp, 9
  la x1, tdat8
  li x2, 0xffffa00a
  la x15, .L7:
  sh x2, 0(x1)
  lh x14, 0(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffa00a
  bne x14, x7, fail
  
  

  ;  Test with a negative base

test_10:
  li gp, 10
  la x1, tdat9
  li x2, 0x12345678
  addi x4, x1, -32
  sh x2, 32(x4)
  lh x5, 0(x1)
  
  li x7, 0x5678
  bne x5, x7, fail
  







  ;  Test with unaligned base

test_11:
  li gp, 11
  la x1, tdat9
  li x2, 0x00003098
  addi x1, x1, -5
  sh x2, 7(x1)
  la x4, tdat10
  lh x5, 0(x4)
  
  li x7, 0x3098
  bne x5, x7, fail
  

  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_12:
  li gp, 12
  li x4, 0
.L1:
  li x1, 0xffffccdd
  la x2, tdat
  sh x1, 0(x2)
  lh x14, 0(x2)
  li x7, 0xffffccdd
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_13:
  li gp, 13
  li x4, 0
.L1:
  li x1, 0xffffbccd
  la x2, tdat
  nop
  sh x1, 2(x2)
  lh x14, 2(x2)
  li x7, 0xffffbccd
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_14:
  li gp, 14
  li x4, 0
.L1:
  li x1, 0xffffbbcc
  la x2, tdat
  nop
  nop
  sh x1, 4(x2)
  lh x14, 4(x2)
  li x7, 0xffffbbcc
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_15:
  li gp, 15
  li x4, 0
.L1:
  li x1, 0xffffabbc
  nop
  la x2, tdat
  sh x1, 6(x2)
  lh x14, 6(x2)
  li x7, 0xffffabbc
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_16:
  li gp, 16
  li x4, 0
.L1:
  li x1, 0xffffaabb
  nop
  la x2, tdat
  nop
  sh x1, 8(x2)
  lh x14, 8(x2)
  li x7, 0xffffaabb
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_17:
  li gp, 17
  li x4, 0
.L1:
  li x1, 0xffffdaab
  nop
  nop
  la x2, tdat
  sh x1, 10(x2)
  lh x14, 10(x2)
  li x7, 0xffffdaab
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

test_18:
  li gp, 18
  li x4, 0
.L1:
  la x2, tdat
  li x1, 0x2233
  sh x1, 0(x2)
  lh x14, 0(x2)
  li x7, 0x2233
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_19:
  li gp, 19
  li x4, 0
.L1:
  la x2, tdat
  li x1, 0x1223
  nop
  sh x1, 2(x2)
  lh x14, 2(x2)
  li x7, 0x1223
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_20:
  li gp, 20
  li x4, 0
.L1:
  la x2, tdat
  li x1, 0x1122
  nop
  nop
  sh x1, 4(x2)
  lh x14, 4(x2)
  li x7, 0x1122
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_21:
  li gp, 21
  li x4, 0
.L1:
  la x2, tdat
  nop
  li x1, 0x0112
  sh x1, 6(x2)
  lh x14, 6(x2)
  li x7, 0x0112
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_22:
  li gp, 22
  li x4, 0
.L1:
  la x2, tdat
  nop
  li x1, 0x0011
  nop
  sh x1, 8(x2)
  lh x14, 8(x2)
  li x7, 0x0011
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_23:
  li gp, 23
  li x4, 0
.L1:
  la x2, tdat
  nop
  nop
  li x1, 0x3001
  sh x1, 10(x2)
  lh x14, 10(x2)
  li x7, 0x3001
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

  li a0, 0xbeef
  la a1, tdat
  sh a0, 6(a1)

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

unimp

  #bank data
 

tdat:
tdat1:
  #d16 le(0xbeef)
tdat2:
  #d16 le(0xbeef)
tdat3:
  #d16 le(0xbeef)
tdat4:
  #d16 le(0xbeef)
tdat5:
  #d16 le(0xbeef)
tdat6:
  #d16 le(0xbeef)
tdat7:
  #d16 le(0xbeef)
tdat8:
  #d16 le(0xbeef)
tdat9:
  #d16 le(0xbeef)
tdat10:
  #d16 le(0xbeef)
