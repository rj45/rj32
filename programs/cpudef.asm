#bits 16

#subruledef reg {
  r0 => 0
  r1 => 1
  r2 => 2
  r3 => 3
}

#ruledef {
  add   {rd:reg}, {value}           => rd`4 @ value`6          @ 0`3 @ 0b001
  add   {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 0`3 @ 0b000
  move  {rd:reg}, {value}           => rd`4 @ value`6          @ 6`3 @ 0b001
  move  {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 0b10   @ 6`3 @ 0b000
  jump  {value}                     => value`12                @ 0`1 @ 0b101
  jump  {rd:reg}                    => rd`4 @ 0`4              @ 0`1 @ 0b000
  if.eq {rd:reg}, {value}           => rd`4 @ value`5 @ 1`3    @ 0`1 @ 0b011
  if.eq {rd:reg}, {rs:reg}          => rd`4 @ rs`4    @ 1`4    @ 0`1 @ 0b010

  load  {rd:reg}, [{rs:reg}, {imm}] => rd`4 @ rs`4    @ imm`5        @ 0b110
  store [{rs:reg}, {imm}], {rd:reg} => rd`4 @ rs`4    @ imm`5        @ 0b111

  halt                              => 0`4  @ 0`4     @ 0b0000 @ 1`1 @ 0b000
  nop                               => 0`4  @ 0`4     @ 0b0000 @ 0`1 @ 0b000
}
