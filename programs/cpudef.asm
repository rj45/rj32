#bits 16

#subruledef reg {
  r0  => 0
  r1  => 1
  r2  => 2
  r3  => 3
  r4  => 4
  r5  => 5
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
  nop                                => 0`4  @ 0`4    @ 0b0000  @ 0`2 @ 0b00
  error                              => 0`4  @ 0`4    @ 0b0000  @ 2`2 @ 0b00
  halt                               => 0`4  @ 0`4    @ 0b0000  @ 3`2 @ 0b00

  jump   {rd:reg}                    => rd`4 @ rd`4   @ 0b00010 @ 0`1 @ 0b00
  return                             =>  0`4 @ 0`4    @ 0b00010 @ 0`1 @ 0b00

  move   {rd:reg}, {value}           => rd`4 @ value`8          @ 0`1 @ 0b001
  move   {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b00011 @ 0`1 @ 0b00

  jump   {value}                     => (value - pc - 1)`11     @ 1`2 @ 0b101
  call   {value}                     => (value - pc - 1)`11     @ 2`2 @ 0b101

  load   {rd:reg}, [{rs:reg}, {imm}] => rd`4 @ rs`4   @ imm`4   @ 0`2 @ 0b10
  store  [{rs:reg}, {imm}], {rd:reg} => rd`4 @ rs`4   @ imm`4   @ 1`2 @ 0b10
  loadb  {rd:reg}, [{rs:reg}, {imm}] => rd`4 @ rs`4   @ imm`4   @ 2`2 @ 0b10
  storeb [{rs:reg}, {imm}], {rd:reg} => rd`4 @ rs`4   @ imm`4   @ 3`2 @ 0b10

  add    {rd:reg}, {value}           => rd`4 @ value`6         @  0`4 @ 0b11
  add    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  0`4 @ 0b00
  sub    {rd:reg}, {value}           => rd`4 @ value`6         @  1`4 @ 0b11
  sub    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  1`4 @ 0b00
  addc   {rd:reg}, {value}           => rd`4 @ value`6         @  2`4 @ 0b11
  addc   {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  2`4 @ 0b00
  subc   {rd:reg}, {value}           => rd`4 @ value`6         @  3`4 @ 0b11
  subc   {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  3`4 @ 0b00
  xor    {rd:reg}, {value}           => rd`4 @ value`6         @  4`4 @ 0b11
  xor    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  4`4 @ 0b00
  and    {rd:reg}, {value}           => rd`4 @ value`6         @  5`4 @ 0b11
  and    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  5`4 @ 0b00
  or     {rd:reg}, {value}           => rd`4 @ value`6         @  6`4 @ 0b11
  or     {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  6`4 @ 0b00
  shl    {rd:reg}, {value}           => rd`4 @ value`6         @  7`4 @ 0b11
  shl    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  7`4 @ 0b00
  shr    {rd:reg}, {value}           => rd`4 @ value`6         @  8`4 @ 0b11
  shr    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  8`4 @ 0b00
  asr    {rd:reg}, {value}           => rd`4 @ value`6         @  9`4 @ 0b11
  asr    {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @  9`4 @ 0b00

  if.eq  {rd:reg}, {value}           => rd`4 @ value`6         @ 10`4 @ 0b11
  if.eq  {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @ 10`4 @ 0b00
  if.ne  {rd:reg}, {value}           => rd`4 @ value`6         @ 11`4 @ 0b11
  if.ne  {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @ 11`4 @ 0b00
  if.lt  {rd:reg}, {value}           => rd`4 @ value`6         @ 12`4 @ 0b11
  if.lt  {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @ 12`4 @ 0b00
  if.ge  {rd:reg}, {value}           => rd`4 @ value`6         @ 13`4 @ 0b11
  if.ge  {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @ 13`4 @ 0b00
  if.lo  {rd:reg}, {value}           => rd`4 @ value`6         @ 14`4 @ 0b11
  if.lo  {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @ 14`4 @ 0b00
  if.hs  {rd:reg}, {value}           => rd`4 @ value`6         @ 15`4 @ 0b11
  if.hs  {rd:reg}, {rs:reg}          => rd`4 @ rs`4   @ 0b01   @ 15`4 @ 0b00
}
