  ; -------------------------------------------------------------
  ;  Logical tests
  ; -------------------------------------------------------------

test_2:
  li gp, 2
  li x1, 0xff00ff00
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  
  li x7, 0x0f000f00
  bne x14, x7, fail
  
  
test_3:
  li gp, 3
  li x1, 0x0ff00ff0
  li x2, 0xf0f0f0f0
  and x14, x1, x2
  
  li x7, 0x00f000f0
  bne x14, x7, fail
  
  
test_4:
  li gp, 4
  li x1, 0x00ff00ff
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  
  li x7, 0x000f000f
  bne x14, x7, fail
  
  
test_5:
  li gp, 5
  li x1, 0xf00ff00f
  li x2, 0xf0f0f0f0
  and x14, x1, x2
  
  li x7, 0xf000f000
  bne x14, x7, fail
  
  

  ; -------------------------------------------------------------
  ;  Source/Destination tests
  ; -------------------------------------------------------------

test_6:
  li gp, 6
  li x1, 0xff00ff00
  li x2, 0x0f0f0f0f
  and x1, x1, x2
  
  li x7, 0x0f000f00
  bne x1, x7, fail
  
  
test_7:
  li gp, 7
  li x1, 0x0ff00ff0
  li x2, 0xf0f0f0f0
  and x2, x1, x2
  
  li x7, 0x00f000f0
  bne x2, x7, fail
  
  
test_8:
  li gp, 8
  li x1, 0xff00ff00
  and x1, x1, x1
  
  li x7, 0xff00ff00
  bne x1, x7, fail
  
  

  ; -------------------------------------------------------------
  ;  Bypassing tests
  ; -------------------------------------------------------------

test_9:
  li gp, 9
  li x4, 0
.L1:
  li x1, 0xff00ff00
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x0f000f00
  bne x6, x7, fail
  
  
test_10:
  li gp, 10
  li x4, 0
.L1:
  li x1, 0x0ff00ff0
  li x2, 0xf0f0f0f0
  and x14, x1, x2
  nop
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x00f000f0
  bne x6, x7, fail
  
  
test_11:
  li gp, 11
  li x4, 0
.L1:
  li x1, 0x00ff00ff
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  nop
  nop
  addi x6, x14, 0
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x000f000f
  bne x6, x7, fail
  
  

test_12:
  li gp, 12
  li x4, 0
.L1:
  li x1, 0xff00ff00
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x0f000f00
  bne x14, x7, fail
  
  
test_13:
  li gp, 13
  li x4, 0
.L1:
  li x1, 0x0ff00ff0
  li x2, 0xf0f0f0f0
  nop
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x00f000f0
  bne x14, x7, fail
  
  
test_14:
  li gp, 14
  li x4, 0
.L1:
  li x1, 0x00ff00ff
  li x2, 0x0f0f0f0f
  nop
  nop
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x000f000f
  bne x14, x7, fail
  
  
test_15:
  li gp, 15
  li x4, 0
.L1:
  li x1, 0xff00ff00
  nop
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x0f000f00
  bne x14, x7, fail
  
  
test_16:
  li gp, 16
  li x4, 0
.L1:
  li x1, 0x0ff00ff0
  nop
  li x2, 0xf0f0f0f0
  nop
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x00f000f0
  bne x14, x7, fail
  
  
test_17:
  li gp, 17
  li x4, 0
.L1:
  li x1, 0x00ff00ff
  nop
  nop
  li x2, 0x0f0f0f0f
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x000f000f
  bne x14, x7, fail
  
  

test_18:
  li gp, 18
  li x4, 0
.L1:
  li x2, 0x0f0f0f0f
  li x1, 0xff00ff00
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x0f000f00
  bne x14, x7, fail
  
  
test_19:
  li gp, 19
  li x4, 0
.L1:
  li x2, 0xf0f0f0f0
  li x1, 0x0ff00ff0
  nop
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x00f000f0
  bne x14, x7, fail
  
  
test_20:
  li gp, 20
  li x4, 0
.L1:
  li x2, 0x0f0f0f0f
  li x1, 0x00ff00ff
  nop
  nop
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x000f000f
  bne x14, x7, fail
  
  
test_21:
  li gp, 21
  li x4, 0
.L1:
  li x2, 0x0f0f0f0f
  nop
  li x1, 0xff00ff00
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x0f000f00
  bne x14, x7, fail
  
  
test_22:
  li gp, 22
  li x4, 0
.L1:
  li x2, 0xf0f0f0f0
  nop
  li x1, 0x0ff00ff0
  nop
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x00f000f0
  bne x14, x7, fail
  
  
test_23:
  li gp, 23
  li x4, 0
.L1:
  li x2, 0x0f0f0f0f
  nop
  nop
  li x1, 0x00ff00ff
  and x14, x1, x2
  addi x4, x4, 1
  li x5, 2
  bne x4, x5, .L1
  li x7, 0x000f000f
  bne x14, x7, fail
  
  

test_24:
  li gp, 24
  li x1, 0xff00ff00
  and x2, x0, x1
  
  li x7, 0
  bne x2, x7, fail
  
  
test_25:
  li gp, 25
  li x1, 0x00ff00ff
  and x2, x1, x0
  
  li x7, 0
  bne x2, x7, fail
  
  
test_26:
  li gp, 26
  and x1, x0, x0
  
  li x7, 0
  bne x1, x7, fail
  
  
test_27:
  li gp, 27
  li x1, 0x11111111
  li x2, 0x22222222
  and x0, x1, x2
  
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

