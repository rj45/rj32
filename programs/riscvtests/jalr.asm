  ; -------------------------------------------------------------
  ;  Test 2:
  ; Basic test
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  li t0, 0
  la t1, target_2

  jalr t0, t1, 0
linkaddr_2:
  j fail

target_2:
  la t1, linkaddr_2
  bne t0, t1, fail

  ; -------------------------------------------------------------
  ;  Test 3:
  ; Basic test2, rs = rd
  ; -------------------------------------------------------------

test_3:
  li gp, 3
  la t0, target_3

  jalr t0, t0, 0
linkaddr_3:
  j fail

target_3:
  la t1, linkaddr_3
  bne t0, t1, fail

  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_4:
  li gp, 4
  li x4, 0
.L1:
  la x6, .L2
  jalr x13, x6, 0
  bne x0, gp, fail
.L2:
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_5:
  li gp, 5
  li x4, 0
.L1:
  la x6, .L2
  nop
  jalr x13, x6, 0
  bne x0, gp, fail
.L2:
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  
test_6:
  li gp, 6
  li x4, 0
.L1:
  la x6, .L2
  nop
  nop
  jalr x13, x6, 0
  bne x0, gp, fail
.L2:
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  

  ; -------------------------------------------------------------
  ;  Test delay slot instructions not executed nor bypassed
  ; -------------------------------------------------------------

  ; .option push
  ; .align 2
  ; .option norvc
test_7:
  li gp, 7
  li t0, 1
  la t1, .L1
  jr t1, -4
  addi t0, t0, 1
  addi t0, t0, 1
  addi t0, t0, 1
  addi t0, t0, 1
.L1:
  addi t0, t0, 1
  addi t0, t0, 1
  
  li x7, 4
  bne t0, x7, fail
  

  ; .option pop

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

