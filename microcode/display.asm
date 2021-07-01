#bits 16

NOP = 0
PLOT = 1
CALL = 2
TYPE = 3
POS = 4
DBA = 6
COLOUR = 7

#ruledef {
  nop            =>            0`13 @ NOP`3
  plot   {value} => 0`8 @   value`5 @ PLOT`3
  call   {value} => 0`5 @   value`8 @ CALL`3
  ret            => 1`5 @       0`8 @ CALL`3
  type   {value} => 0`8 @   value`5 @ TYPE`3
  pos    {x},{y} => 0`4 @ y`4 @ x`5 @ POS`3
  value  {value} => 0`8 @   value`5 @ DBA`3
  colour {value} => 0`8 @   value`5 @ COLOUR`3
}

teal = 1
blue = 2
orange = 3
red = 4
purple = 5
green = 6

column1 = 1
column2 = column1 + 8
column3 = column2 + 8
column4 = column3 + 8

number = 0
opname = 1
regname = 2
aluop = 3
text = 7

; jump past subroutines
call start

num:
  nop ; these nops give the bin -> bcd time to calculate
  nop
  nop
  nop
  plot 0
  plot 1
  plot 2
  plot 3
  plot 4
  plot 5
  ret

start:

type number

; top row registers
colour red

; r0
value 0
pos column1, 1
call num

; r4
value 4
pos column2, 1
call num

; r8
value 8
pos column3, 1
call num

; r12
value 12
pos column4, 1
call num

; second row

; r2
value 1
pos column1, 2
call num

; r5
value 5
pos column2, 2
call num

; r9
value 9
pos column3, 2
call num

; r13
value 13
pos column4, 2
call num

; third row

; r3
value 2
pos column1, 3
call num

; r6
value 6
pos column2, 3
call num

; r10
value 10
pos column3, 3
call num

; r14
value 14
pos column4, 3
call num

; forth row

; r4
value 3
pos column1, 4
call num

; r7
value 7
pos column2, 4
call num

; r11
value 11
pos column3, 4
call num

; r15
value 15
pos column4, 4
call num

; imm
value 0x16
pos column4, 6
call num

; result
value 0x14
pos column2, 8
call num

; L
value 0x12
pos column3, 8
call num

; R
value 0x13
pos column4, 8
call num

; PC
value 0x17
pos column1, 10
call num

; op
value 0x10
colour green
type opname
pos 1, 6
plot 0
plot 1
plot 2
plot 3
plot 4
plot 5
plot 6
plot 7

; register selections
colour blue
type regname

; RS
value 0x19
pos 20, 6
plot 0
plot 1
plot 2

; RD
value 0x18
pos 12, 6
plot 0
plot 1
plot 2
pos 1, 8
plot 0
plot 1
plot 2

; aluop in the ctrl reg
value 0x11
type aluop
colour orange
pos 22, 8
plot 1
plot 2

; a couple blue `=` signs
pos 6, 8
type text
colour blue
plot 0
pos 16, 8
plot 0

; loop back to the beginning
call start