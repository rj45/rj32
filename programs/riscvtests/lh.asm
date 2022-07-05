  ; -------------------------------------------------------------
  ;  Basic tests
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  li x15, 0x000000ff
  la x1, tdat
  lh x14, 0(x1)
  
  li x7, 0x000000ff
  bne x14, x7, fail
  
  
test_3:
  li gp, 3
  li x15, 0xffffff00
  la x1, tdat
  lh x14, 2(x1)
  
  li x7, 0xffffff00
  bne x14, x7, fail
  
  
test_4:
  li gp, 4
  li x15, 0x00000ff0
  la x1, tdat
  lh x14, 4(x1)
  
  li x7, 0x00000ff0
  bne x14, x7, fail
  
  
test_5:
  li gp, 5
  li x15, 0xfffff00f
  la x1, tdat
  lh x14, 6(x1)
  
  li x7, 0xfffff00f
  bne x14, x7, fail
  
  

  ;  Test with negative offset

test_6:
  li gp, 6
  li x15, 0x000000ff
  la x1, tdat4
  lh x14, -6(x1)
  
  li x7, 0x000000ff
  bne x14, x7, fail
  
  
test_7:
  li gp, 7
  li x15, 0xffffff00
  la x1, tdat4
  lh x14, -4(x1)
  
  li x7, 0xffffff00
  bne x14, x7, fail
  
  
test_8:
  li gp, 8
  li x15, 0x00000ff0
  la x1, tdat4
  lh x14, -2(x1)
  
  li x7, 0x00000ff0
  bne x14, x7, fail
  
  
test_9:
  li gp, 9
  li x15, 0xfffff00f
  la x1, tdat4
  lh x14, 0(x1)
  
  li x7, 0xfffff00f
  bne x14, x7, fail
  
  

  ;  Test with a negative base

test_10:
  li gp, 10
  la x1, tdat
  addi x1, x1, -32
  lh x5, 32(x1)
  
  li x7, 0x000000ff
  bne x5, x7, fail
  





  ;  Test with unaligned base

test_11:
  li gp, 11
  la x1, tdat
  addi x1, x1, -5
  lh x5, 7(x1)
  
  li x7, 0xffffff00
  bne x5, x7, fail
  





  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_12:
  li gp, 12
  li x4, 0
.L1:
  la x1, tdat2
  lh x14, 2(x1)
  addi x6, x14, 0
  li x7, 0x00000ff0
  bne x6, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
  
test_13:
  li gp, 13
  li x4, 0
.L1:
  la x1, tdat3
  lh x14, 2(x1)
  nop
  addi x6, x14, 0
  li x7, 0xfffff00f
  bne x6, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
  
test_14:
  li gp, 14
  li x4, 0
.L1:
  la x1, tdat1
  lh x14, 2(x1)
  nop
  nop
  addi x6, x14, 0
  li x7, 0xffffff00
  bne x6, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
  

test_15:
  li gp, 15
  li x4, 0
.L1:
  la x1, tdat2
  lh x14, 2(x1)
  li x7, 0x00000ff0
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_16:
  li gp, 16
  li x4, 0
.L1:
  la x1, tdat3
  nop
  lh x14, 2(x1)
  li x7, 0xfffff00f
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_17:
  li gp, 17
  li x4, 0
.L1:
  la x1, tdat1
  nop
  nop
  lh x14, 2(x1)
  li x7, 0xffffff00
  bne x14, x7, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

  ; -------------------------------------------------------------
  ;  Test write-after-write hazard
  ; -------------------------------------------------------------

test_18:
  li gp, 18
  la x5, tdat
  lh x2, 0(x5)
  li x2, 2
  
  li x7, 2
  bne x2, x7, fail
  





test_19:
  li gp, 19
  la x5, tdat
  lh x2, 0(x5)
  nop
  li x2, 2
  
  li x7, 2
  bne x2, x7, fail
  






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

  #bank data
 
tdat:
tdat1:
  #d16 le(0x00ff)
tdat2:
  #d16 le(0xff00)
tdat3:
  #d16 le(0x0ff0)
tdat4:
  #d16 le(0xf00f)
