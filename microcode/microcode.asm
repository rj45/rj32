#bits 32

HALT  = 1 << 0
ERROR = 1 << 1
JUMP  = 1 << 2

; alu ops
ADD     = 0 << 3
SUB     = 1 << 3
SHR     = 2 << 3
SHL     = 3 << 3
ASR     = 4 << 3
XOR     = 5 << 3
AND     = 6 << 3
OR      = 7 << 3

; mem ops
MEM   = 1 << 6
STORE = 1 << 7

; reg write
WRITE = 1 << 8

; wrmux -- write mux
WM_RESULT  = 0 << 9
WM_L       = 1 << 9
WM_MEMDATA = 2 << 9
WM_R       = 3 << 9

; cond
EQ = 1 << 11
NE = 2 << 11
LT = 3 << 11
GE = 4 << 11
LO = 5 << 11
HS = 6 << 11

; imm bit
IMM = 1 << 14

; flags
FL_ACK  = 1 << 0
FL_NOP1 = 1 << 1
FL_NOP2 = 1 << 2

#ruledef {
  done {value}                            => le(    0`5 @     0`5 @     0`3 @ value`19)
  done                                    => le(    0`5 @     0`5 @     0`3 @     0`19)
  next {value}, {flags}, {nextt}, {nextf} => le(nextf`5 @ nextt`5 @ flags`3 @ value`19)
  loop {value}, {nextt}                   => le(nextt`5 @ nextt`5 @     0`3 @ value`19)
}

nop:
  done
rets:
  done
error:
  loop ERROR, infinierror
halt:
  loop HALT, infinihalt

#addr 0b00110
move:
  done WM_R | WRITE

#addr 0b01000
jump:
  done JUMP | ADD
imm:
  done IMM
call:
  done JUMP | ADD | WM_L | WRITE
imm2:
  done IMM

#addr 0b01100
load:
  next MEM | ADD, FL_ACK, loaddone, loadwait
store:
  next MEM | ADD | STORE, FL_ACK, storedone, storewait
loadb:
  done WM_MEMDATA | MEM | ADD | WRITE
storeb:
  done ADD | MEM | STORE

#addr 0b10000
add:
  done ADD | WRITE
sub:
  done SUB | WRITE
addc:
  done ADD | WRITE
subc:
  done SUB | WRITE
xor:
  done XOR | WRITE
and:
  done AND | WRITE
or:
  done OR | WRITE
shl:
  done SHL | WRITE
shr:
  done SHR | WRITE
asr:
  done ASR | WRITE
ifeq:
  done SUB | EQ
ifne:
  done SUB | NE
iflt:
  done SUB | LT
ifge:
  done SUB | GE
iflo:
  done SUB | LO
ifhs:
  done SUB | HS

#addr 0b100001
infinierror:
  loop ERROR, infinierror
infinihalt:
  loop HALT, infinihalt
loaddone:
  done WRITE | WM_MEMDATA
loadwait:
  next MEM | ADD, FL_ACK, loaddone, loadwait
storedone:
  done
storewait:
  next MEM | ADD | STORE, FL_ACK, storedone, storewait
