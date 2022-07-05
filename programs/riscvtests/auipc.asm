test_2:
  li gp, 2
  ; .align 3
  lla a0, .L1 + 10000
  jal a1, .L1
.L1:
  sub a0, a0, a1
  
  li x7, 10000
  bne a0, x7, fail
  






test_3:
  li gp, 3
  ; .align 3
  lla a0, .L1 - 10000
  jal a1, .L1
.L1:
  sub a0, a0, a1
  
  li x7, ((-10000) & ((1 << (32 - 1) << 1) - 1))
  bne a0, x7, fail
  






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
