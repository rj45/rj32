#bits 16

#subruledef reg {
  r0  => 0
  r1  => 1
  r2  => 2
  r3  => 3
  r4  => 4
  r5  => 3
  r6  => 6
  r7  => 7
  r8  => 8
  r9  => 9
  r10 => 10
  r11 => 11
  r12 => 12
  r13 => 13
  r14 => 14
  r15 => 15
}

#ruledef {
  add   {rd:reg}, {value}           => rd`4 @ value`6          @ 0`3 @ 0b001
  add   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 0`3 @ 0b000
  sub   {rd:reg}, {value}           => rd`4 @ value`6          @ 1`3 @ 0b001
  sub   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 1`3 @ 0b000
  xor   {rd:reg}, {value}           => rd`4 @ value`6          @ 2`3 @ 0b001
  xor   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 2`3 @ 0b000
  and   {rd:reg}, {value}           => rd`4 @ value`6          @ 3`3 @ 0b001
  and   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 3`3 @ 0b000
  or    {rd:reg}, {value}           => rd`4 @ value`6          @ 4`3 @ 0b001
  or    {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 4`3 @ 0b000
  shl   {rd:reg}, {value}           => rd`4 @ value`6          @ 5`3 @ 0b001
  shl   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 5`3 @ 0b000
  shr   {rd:reg}, {value}           => rd`4 @ value`6          @ 6`3 @ 0b001
  shr   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 6`3 @ 0b000
  asr   {rd:reg}, {value}           => rd`4 @ value`6          @ 7`3 @ 0b001
  asr   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 7`3 @ 0b000

  move  {rd:reg}, {value}           => rd`4 @ value`8          @ 0`1 @ 0b100
  move  {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b1101 @ 0`1 @ 0b000
  jump  {value}                     => (value - pc - 1)`12     @ 0`1 @ 0b101
  call  {value}                     => (value - pc - 1)`12     @ 1`1 @ 0b101

  jump  {rd:reg}                    => rd`4 @ rd`4    @ 0b0010 @ 0`1 @ 0b000
  return                            =>  0`4 @ 0`4     @ 0b0010 @ 0`1 @ 0b000

  if.eq {rd:reg}, {value}           => rd`4 @ value`5 @ 1`3    @ 0`1 @ 0b011
  if.eq {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 1`4    @ 0`1 @ 0b010
  if.ne {rd:reg}, {value}           => rd`4 @ value`5 @ 2`3    @ 0`1 @ 0b011
  if.ne {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 2`4    @ 0`1 @ 0b010
  if.lt {rd:reg}, {value}           => rd`4 @ value`5 @ 3`3    @ 0`1 @ 0b011
  if.lt {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 3`4    @ 0`1 @ 0b010
  if.ge {rd:reg}, {value}           => rd`4 @ value`5 @ 4`3    @ 0`1 @ 0b011
  if.ge {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 4`4    @ 0`1 @ 0b010
  if.lo {rd:reg}, {value}           => rd`4 @ value`5 @ 5`3    @ 0`1 @ 0b011
  if.lo {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 5`4    @ 0`1 @ 0b010
  if.hs {rd:reg}, {value}           => rd`4 @ value`5 @ 6`3    @ 0`1 @ 0b011
  if.hs {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 6`4    @ 0`1 @ 0b010

  load  {rd:reg}, [{rs:reg}, {imm}] => rd`4 @ rs`4    @ imm`5        @ 0b110
  store [{rs:reg}, {imm}], {rd:reg} => rd`4 @ rs`4    @ imm`5        @ 0b111

  nop                               => 0`4  @ 0`4     @ 0b000  @ 0`2 @ 0b000
  error                             => 0`4  @ 0`4     @ 0b000  @ 2`2 @ 0b000
  halt                              => 0`4  @ 0`4     @ 0b000  @ 3`2 @ 0b000
}
