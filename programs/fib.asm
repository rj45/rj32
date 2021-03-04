#bits 16

#subruledef reg {
  r0 => 0
  r1 => 1
  r2 => 2
  r3 => 3
}

#ruledef {
  add {rd:reg}, {value}   => value`8 @ 0`2  @ rd`2 @ 0`3 @ 0`1
  add {rd:reg}, {rs:reg}  => 0`8     @ rs`2 @ rd`2 @ 0`3 @ 1`1
  move {rd:reg}, {value}  => value`8 @ 0`2  @ rd`2 @ 6`3 @ 0`1
  move {rd:reg}, {rs:reg} => 0`8     @ rs`2 @ rd`2 @ 6`3 @ 1`1
  jump {value}            => value`8 @ 0`2  @ 0`2  @ 7`3 @ 0`1
  jump {rs:reg}           => value`8 @ rs`2 @ 0`2  @ 7`3 @ 0`1
}

fibonacci:
  add r2, 1
  .loop:
    add r2, r1
    move r3, r2
    add r1, r2
    move r3, r1
    jump .loop

