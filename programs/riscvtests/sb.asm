  ; -------------------------------------------------------------
  ;  Basic tests
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  la x1, tdat
  li x2, 0xffffffaa
  la x15, .L7:
  sb x2, 0(x1)
  lb x14, 0(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffffaa
  bne x14, x7, fail
  
  
test_3:
  li gp, 3
  la x1, tdat
  li x2, 0x00000000
  la x15, .L7:
  sb x2, 1(x1)
  lb x14, 1(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x00000000
  bne x14, x7, fail
  
  
test_4:
  li gp, 4
  la x1, tdat
  li x2, 0xffffefa0
  la x15, .L7:
  sb x2, 2(x1)
  lh x14, 2(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffefa0
  bne x14, x7, fail
  
  
test_5:
  li gp, 5
  la x1, tdat
  li x2, 0x0000000a
  la x15, .L7:
  sb x2, 3(x1)
  lb x14, 3(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x0000000a
  bne x14, x7, fail
  
  

  ;  Test with negative offset

test_6:
  li gp, 6
  la x1, tdat8
  li x2, 0xffffffaa
  la x15, .L7:
  sb x2, -3(x1)
  lb x14, -3(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffffaa
  bne x14, x7, fail
  
  
test_7:
  li gp, 7
  la x1, tdat8
  li x2, 0x00000000
  la x15, .L7:
  sb x2, -2(x1)
  lb x14, -2(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x00000000
  bne x14, x7, fail
  
  
test_8:
  li gp, 8
  la x1, tdat8
  li x2, 0xffffffa0
  la x15, .L7:
  sb x2, -1(x1)
  lb x14, -1(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0xffffffa0
  bne x14, x7, fail
  
  
test_9:
  li gp, 9
  la x1, tdat8
  li x2, 0x0000000a
  la x15, .L7:
  sb x2, 0(x1)
  lb x14, 0(x1)
  j .L8:
.L7:
  mv x14, x2
.L8:
  li x7, 0x0000000a
  bne x14, x7, fail
  
  

  ;  Test with a negative base

test_10:
  li gp, 10
  la x1, tdat9
  li x2, 0x12345678
  addi x4, x1, -32
  sb x2, 32(x4)
  lb x5, 0(x1)
  
  li x7, 0x78
  bne x5, x7, fail
  







  ;  Test with unaligned base

test_11:
  li gp, 11
  la x1, tdat9
  li x2, 0x00003098
  addi x1, x1, -6
  sb x2, 7(x1)
  la x4, tdat10
  lb x5, 0(x4)
  
  li x7, 0xffffff98
  bne x5, x7, fail
  

  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_12:
  li gp, 12
  li x4, 0
.L1:
  li x1, 0xffffffdd
  la x2, tdat
  sb x1, 0(x2)
  lb x14, 0(x2)
  li x7, 0xffffffdd
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_13:
  li gp, 13
  li x4, 0
.L1:
  li x1, 0xffffffcd
  la x2, tdat
  nop
  sb x1, 1(x2)
  lb x14, 1(x2)
  li x7, 0xffffffcd
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_14:
  li gp, 14
  li x4, 0
.L1:
  li x1, 0xffffffcc
  la x2, tdat
  nop
  nop
  sb x1, 2(x2)
  lb x14, 2(x2)
  li x7, 0xffffffcc
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_15:
  li gp, 15
  li x4, 0
.L1:
  li x1, 0xffffffbc
  nop
  la x2, tdat
  sb x1, 3(x2)
  lb x14, 3(x2)
  li x7, 0xffffffbc
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_16:
  li gp, 16
  li x4, 0
.L1:
  li x1, 0xffffffbb
  nop
  la x2, tdat
  nop
  sb x1, 4(x2)
  lb x14, 4(x2)
  li x7, 0xffffffbb
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_17:
  li gp, 17
  li x4, 0
.L1:
  li x1, 0xffffffab
  nop
  nop
  la x2, tdat
  sb x1, 5(x2)
  lb x14, 5(x2)
  li x7, 0xffffffab
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

test_18:
  li gp, 18
  li x4, 0
.L1:
  la x2, tdat
  li x1, 0x33
  sb x1, 0(x2)
  lb x14, 0(x2)
  li x7, 0x33
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_19:
  li gp, 19
  li x4, 0
.L1:
  la x2, tdat
  li x1, 0x23
  nop
  sb x1, 1(x2)
  lb x14, 1(x2)
  li x7, 0x23
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_20:
  li gp, 20
  li x4, 0
.L1:
  la x2, tdat
  li x1, 0x22
  nop
  nop
  sb x1, 2(x2)
  lb x14, 2(x2)
  li x7, 0x22
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
  li x1, 0x12
  sb x1, 3(x2)
  lb x14, 3(x2)
  li x7, 0x12
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
  li x1, 0x11
  nop
  sb x1, 4(x2)
  lb x14, 4(x2)
  li x7, 0x11
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
  li x1, 0x01
  sb x1, 5(x2)
  lb x14, 5(x2)
  li x7, 0x01
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

  li a0, 0xef
  la a1, tdat
  sb a0, 3(a1)

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
  #d8 0xef
tdat2:
  #d8 0xef
tdat3:
  #d8 0xef
tdat4:
  #d8 0xef
tdat5:
  #d8 0xef
tdat6:
  #d8 0xef
tdat7:
  #d8 0xef
tdat8:
  #d8 0xef
tdat9:
  #d8 0xef
tdat10:
  #d8 0xef
