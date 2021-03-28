#bits 16

#subruledef reg {
  r0 => 0
  r1 => 1
  r2 => 2
  r3 => 3
}

#ruledef {
  add  {rd:reg}, {value}    =>   value`5 @  0`4 @ rd`4 @ 0b011
  add  {rd:reg}, {rs:reg}   => 1`1 @ 0`4 @ rs`4 @ rd`4 @ 0b010
  move {rd:reg}, {value}    =>   value`5 @  6`4 @ rd`4 @ 0b011
  move {rd:reg}, {rs:reg}   => 1`1 @ 6`4 @ rs`4 @ rd`4 @ 0b010
  eq   {rd:reg}, {value}    =>   value`5 @  4`4 @ rd`4 @ 0b011
  eq   {rd:reg}, {rs:reg}   => 1`1 @ 4`4 @ rs`4 @ rd`4 @ 0b010
  jump {value}              =>  value`11 @         0`2 @ 0b001
  jump {rs:reg}             => 0`1 @ 4`4 @  0`4 @ rd`4 @ 0b010
  brt  {value}              =>  value`11 @         2`2 @ 0b001
  brt  {rs:reg}             => 0`1 @ 4`4 @  2`4 @ rd`4 @ 0b010
  brf  {value}              =>  value`11 @         3`2 @ 0b001
  brf  {rs:reg}             => 0`1 @ 4`4 @  3`4 @ rd`4 @ 0b010

  halt                      => 0`11             @  1`2 @ 0b000
}
