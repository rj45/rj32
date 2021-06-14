#bits 16

HALT  = 1 << 0
ERROR = 1 << 1
JUMP  = 1 << 2

; alu ops
NOALU   = 0 << 3
ADD     = 1 << 3
SUB     = 2 << 3
MOVE    = 4 << 3
XOR     = 5 << 3
AND     = 6 << 3
OR      = 7 << 3

; mem ops
MEM   = 1 << 6
STORE = 1 << 7

; reg write
WRITE = 1 << 8

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
  done ADD | WRITE
sub:
  done SUB | WRITE
xor:
  done XOR | WRITE
and:
  done AND | WRITE
or:
  done OR | WRITE
alu_5:
  done
move:
  done MOVE | WRITE
noalu:
  done

#addr 0b11000
ifcc:
  done SUB

#addr 0b11100
jump:
  done JUMP | ADD

#addr 0b11110
load:
  done MEM | ADD | WRITE

store:
  done MEM | ADD | STORE