  ; -------------------------------------------------------------
  ;  Basic tests
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  lui x1, 0x00000
  li x7, 0x00000000
  bne x1, x7, fail
  
  
test_3:
  li gp, 3
  lui x1, 0xfffff
  sra x1,x1,1
  li x7, 0xfffff800
  bne x1, x7, fail
  
  
test_4:
  li gp, 4
  lui x1, 0x7ffff
  sra x1,x1,20
  li x7, 0x000007ff
  bne x1, x7, fail
  
  
test_5:
  li gp, 5
  lui x1, 0x80000
  sra x1,x1,20
  li x7, 0xfffff800
  bne x1, x7, fail
  
  

test_6:
  li gp, 6
  lui x0, 0x80000
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
