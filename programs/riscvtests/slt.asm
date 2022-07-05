  ; -------------------------------------------------------------
  ;  Arithmetic tests
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  li x1, 0x00000000
  li x2, 0x00000000
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  
test_3:
  li gp, 3
  li x1, 0x00000001
  li x2, 0x00000001
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  
test_4:
  li gp, 4
  li x1, 0x00000003
  li x2, 0x00000007
  slt x14, x1, x2
  
  li x7, 1
  bne x14, x7, fail
  
  
test_5:
  li gp, 5
  li x1, 0x00000007
  li x2, 0x00000003
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  

test_6:
  li gp, 6
  li x1, 0x00000000
  li x2, 0xffff8000
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  
test_7:
  li gp, 7
  li x1, 0x80000000
  li x2, 0x00000000
  slt x14, x1, x2
  
  li x7, 1
  bne x14, x7, fail
  
  
test_8:
  li gp, 8
  li x1, 0x80000000
  li x2, 0xffff8000
  slt x14, x1, x2
  
  li x7, 1
  bne x14, x7, fail
  
  

test_9:
  li gp, 9
  li x1, 0x00000000
  li x2, 0x00007fff
  slt x14, x1, x2
  
  li x7, 1
  bne x14, x7, fail
  
  
test_10:
  li gp, 10
  li x1, 0x7fffffff
  li x2, 0x00000000
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  
test_11:
  li gp, 11
  li x1, 0x7fffffff
  li x2, 0x00007fff
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  

test_12:
  li gp, 12
  li x1, 0x80000000
  li x2, 0x00007fff
  slt x14, x1, x2
  
  li x7, 1
  bne x14, x7, fail
  
  
test_13:
  li gp, 13
  li x1, 0x7fffffff
  li x2, 0xffff8000
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  

test_14:
  li gp, 14
  li x1, 0x00000000
  li x2, 0xffffffff
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  
test_15:
  li gp, 15
  li x1, 0xffffffff
  li x2, 0x00000001
  slt x14, x1, x2
  
  li x7, 1
  bne x14, x7, fail
  
  
test_16:
  li gp, 16
  li x1, 0xffffffff
  li x2, 0xffffffff
  slt x14, x1, x2
  
  li x7, 0
  bne x14, x7, fail
  
  

  ; -------------------------------------------------------------
  ;  Source/Destination tests
  ; -------------------------------------------------------------

test_17:
  li gp, 17
  li x1, 14
  li x2, 13
  slt x1, x1, x2
  
  li x7, 0
  bne x1, x7, fail
  
  
test_18:
  li gp, 18
  li x1, 11
  li x2, 13
  slt x2, x1, x2
  
  li x7, 1
  bne x2, x7, fail
  
  
test_19:
  li gp, 19
  li x1, 13
  slt x1, x1, x1
  
  li x7, 0
  bne x1, x7, fail
  
  

  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_20:
  li gp, 20
  li x4, 0
.L1:
  li x1, 11
  li x2, 13
  slt x14, x1, x2
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x6, x7, fail
  
  
test_21:
  li gp, 21
  li x4, 0
.L1:
  li x1, 14
  li x2, 13
  slt x14, x1, x2
  nop
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x6, x7, fail
  
  
test_22:
  li gp, 22
  li x4, 0
.L1:
  li x1, 12
  li x2, 13
  slt x14, x1, x2
  nop
  nop
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x6, x7, fail
  
  

test_23:
  li gp, 23
  li x4, 0
.L1:
  li x1, 14
  li x2, 13
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x14, x7, fail
  
  
test_24:
  li gp, 24
  li x4, 0
.L1:
  li x1, 11
  li x2, 13
  nop
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x14, x7, fail
  
  
test_25:
  li gp, 25
  li x4, 0
.L1:
  li x1, 15
  li x2, 13
  nop
  nop
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x14, x7, fail
  
  
test_26:
  li gp, 26
  li x4, 0
.L1:
  li x1, 10
  nop
  li x2, 13
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x14, x7, fail
  
  
test_27:
  li gp, 27
  li x4, 0
.L1:
  li x1, 16
  nop
  li x2, 13
  nop
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x14, x7, fail
  
  
test_28:
  li gp, 28
  li x4, 0
.L1:
  li x1, 9
  nop
  nop
  li x2, 13
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x14, x7, fail
  
  

test_29:
  li gp, 29
  li x4, 0
.L1:
  li x2, 13
  li x1, 17
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x14, x7, fail
  
  
test_30:
  li gp, 30
  li x4, 0
.L1:
  li x2, 13
  li x1, 8
  nop
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x14, x7, fail
  
  
test_31:
  li gp, 31
  li x4, 0
.L1:
  li x2, 13
  li x1, 18
  nop
  nop
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x14, x7, fail
  
  
test_32:
  li gp, 32
  li x4, 0
.L1:
  li x2, 13
  nop
  li x1, 7
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x14, x7, fail
  
  
test_33:
  li gp, 33
  li x4, 0
.L1:
  li x2, 13
  nop
  li x1, 19
  nop
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0
  bne x14, x7, fail
  
  
test_34:
  li gp, 34
  li x4, 0
.L1:
  li x2, 13
  nop
  nop
  li x1, 6
  slt x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 1
  bne x14, x7, fail
  
  

test_35:
  li gp, 35
  li x1, ((-1) & ((1 << (32 - 1) << 1) - 1))
  slt x2, x0, x1
  
  li x7, 0
  bne x2, x7, fail
  
  
test_36:
  li gp, 36
  li x1, ((-1) & ((1 << (32 - 1) << 1) - 1))
  slt x2, x1, x0
  
  li x7, 1
  bne x2, x7, fail
  
  
test_37:
  li gp, 37
  slt x1, x0, x0
  
  li x7, 0
  bne x1, x7, fail
  
  
test_38:
  li gp, 38
  li x1, 16
  li x2, 30
  slt x0, x1, x2
  
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
