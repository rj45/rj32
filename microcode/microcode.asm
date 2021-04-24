#bits 16

HALT  = 1 << 0
ERROR = 1 << 1
JUMP  = 1 << 2

; alu ops
NOALU = 0 << 3
ADD   = 1 << 3
MOVE  = 6 << 3

; mem ops
MEM   = 1 << 6
STORE = 1 << 7

#ruledef {
  done {value}       => value`16
  done               => 0`16
}

nop:
  done
rets:
  done
error:
  done ERROR
halt:
  done HALT

#addr 0b10000
add:
  done ADD
alu_1:
  done
alu_2:
  done
alu_3:
  done
alu_4:
  done
alu_5:
  done
alu_6:
  done MOVE
noalu:
  done

#addr 0b11100
jump:
  done JUMP

#addr 0b11110
load:
  done MEM | ADD

store:
  done MEM | ADD | STORE