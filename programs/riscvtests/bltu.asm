  ; -------------------------------------------------------------
  ;  Branch tests
  ; -------------------------------------------------------------

  ;  Each test checks both forward and backward branches

test_2:
  li gp, 2
  li x1, 0x00000000
  li x2, 0x00000001
  bltu x1, x2, .L2
  bne x0, gp, fail
.L1:
  bne x0, gp, .L3
.L2:
  bltu x1, x2, .L1
  bne x0, gp, fail
.L3:
  
test_3:
  li gp, 3
  li x1, 0xfffffffe
  li x2, 0xffffffff
  bltu x1, x2, .L2
  bne x0, gp, fail
.L1:
  bne x0, gp, .L3
.L2:
  bltu x1, x2, .L1
  bne x0, gp, fail
.L3:
  
test_4:
  li gp, 4
  li x1, 0x00000000
  li x2, 0xffffffff
  bltu x1, x2, .L2
  bne x0, gp, fail
.L1:
  bne x0, gp, .L3
.L2:
  bltu x1, x2, .L1
  bne x0, gp, fail
.L3:
  

test_5:
  li gp, 5
  li x1, 0x00000001
  li x2, 0x00000000
  bltu x1, x2, .L1
  bne x0, gp, .L2
.L1:
  bne x0, gp, fail
.L2:
  bltu x1, x2, .L1
.L3:
  
test_6:
  li gp, 6
  li x1, 0xffffffff
  li x2, 0xfffffffe
  bltu x1, x2, .L1
  bne x0, gp, .L2
.L1:
  bne x0, gp, fail
.L2:
  bltu x1, x2, .L1
.L3:
  
test_7:
  li gp, 7
  li x1, 0xffffffff
  li x2, 0x00000000
  bltu x1, x2, .L1
  bne x0, gp, .L2
.L1:
  bne x0, gp, fail
.L2:
  bltu x1, x2, .L1
.L3:
  
test_8:
  li gp, 8
  li x1, 0x80000000
  li x2, 0x7fffffff
  bltu x1, x2, .L1
  bne x0, gp, .L2
.L1:
  bne x0, gp, fail
.L2:
  bltu x1, x2, .L1
.L3:
  

  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_9:
  li gp, 9
  li x4, 0
.L1:
  li x1, 0xf0000000
  li x2, 0xefffffff
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_10:
  li gp, 10
  li x4, 0
.L1:
  li x1, 0xf0000000
  li x2, 0xefffffff
  nop
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_11:
  li gp, 11
  li x4, 0
.L1:
  li x1, 0xf0000000
  li x2, 0xefffffff
  nop
  nop
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_12:
  li gp, 12
  li x4, 0
.L1:
  li x1, 0xf0000000
  nop
  li x2, 0xefffffff
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_13:
  li gp, 13
  li x4, 0
.L1:
  li x1, 0xf0000000
  nop
  li x2, 0xefffffff
  nop
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_14:
  li gp, 14
  li x4, 0
.L1:
  li x1, 0xf0000000
  nop
  nop
  li x2, 0xefffffff
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

test_15:
  li gp, 15
  li x4, 0
.L1:
  li x1, 0xf0000000
  li x2, 0xefffffff
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_16:
  li gp, 16
  li x4, 0
.L1:
  li x1, 0xf0000000
  li x2, 0xefffffff
  nop
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_17:
  li gp, 17
  li x4, 0
.L1:
  li x1, 0xf0000000
  li x2, 0xefffffff
  nop
  nop
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_18:
  li gp, 18
  li x4, 0
.L1:
  li x1, 0xf0000000
  nop
  li x2, 0xefffffff
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_19:
  li gp, 19
  li x4, 0
.L1:
  li x1, 0xf0000000
  nop
  li x2, 0xefffffff
  nop
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_20:
  li gp, 20
  li x4, 0
.L1:
  li x1, 0xf0000000
  nop
  nop
  li x2, 0xefffffff
  bltu x1, x2, fail
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

  ; -------------------------------------------------------------
  ;  Test delay slot instructions not executed nor bypassed
  ; -------------------------------------------------------------

test_21:
  li gp, 21
  li x1, 1
  bltu x0, x1, .L1
  addi x1, x1, 1
  addi x1, x1, 1
  addi x1, x1, 1
  addi x1, x1, 1
.L1:
  addi x1, x1, 1
  addi x1, x1, 1
  
  li x7, 3
  bne x1, x7, fail
  

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
