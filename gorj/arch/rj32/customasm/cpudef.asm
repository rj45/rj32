#bits 16

; special label for nil pointers
nil = 0

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

  ; ABI names
  ra => 0
  a0 => 1
  a1 => 2
  a2 => 3
  s0 => 4
  s1 => 5
  s2 => 6
  s3 => 7
  t0 => 8
  t1 => 9
  t2 => 10
  t3 => 11
  t4 => 12
  t5 => 13
  gp => 14
  bp => 14
  sp => 15
}

#subruledef op {
  nop     => 0
  error   => 2
  halt    => 3
  move    => 6
  movei   => 0
  jump    => 8
  jumpi   => 0
  calli   => 1
  load    => 0
  store   => 1
  loadb   => 2
  storeb  => 3
  addi    => 0
  subi    => 1
  addci   => 2
  subci   => 3
  xori    => 4
  andi    => 5
  ori     => 6
  shli    => 7
  shri    => 8
  asri    => 9
  if.eqi  => 10
  if.nei  => 11
  if.lti  => 12
  if.gei  => 13
  if.ulti => 14
  if.ugei => 15
  add     => 16 | 0
  sub     => 16 | 1
  addc    => 16 | 2
  subc    => 16 | 3
  xor     => 16 | 4
  and     => 16 | 5
  or      => 16 | 6
  shl     => 16 | 7
  shr     => 16 | 8
  asr     => 16 | 9
  if.eq   => 16 | 10
  if.ne   => 16 | 11
  if.lt   => 16 | 12
  if.ge   => 16 | 13
  if.ult  => 16 | 14
  if.uge  => 16 | 15
}

#ruledef {
  fmt_ri6 {op:op}, {rd:reg}, {value}           => {
    assert(value < (1<<5) && value >= -(1<<5))
    rd`4 @ value`6         @ op`4 @ 0b11
  }
  fmt_ri6 {op:op}, {rd:reg}, {value}           => {
    assert(value >= (1<<5) || value < -(1<<5))
    asm { imm value } @
    rd`4 @ value`6         @ op`4 @ 0b11
  }

  fmt_ri8 {op:op}, {rd:reg}, {value}           => {
    assert(value < (1<<7) && value >= -(1<<7))
    rd`4 @ value`8         @ op`1 @ 0b001
  }
  fmt_ri8 {op:op}, {rd:reg}, {value}           => {
    assert(value >= (1<<7) || value < -(1<<7))
    asm { imm value } @
    rd`4 @ value`8         @ op`1 @ 0b001
  }

  fmt_i11 {op:op}, {value}                     => {
    assert(value < (1<<10) && value >= -(1<<10))
    value`11               @ op`1 @ 0b0101
  }
  fmt_i11 {op:op}, {value}                     => {
    assert(value >= (1<<10) || value < -(1<<10))
    asm { imm value } @
    value`11               @ op`1 @ 0b0101
  }

  fmt_ls  {op:op}, {rd:reg}, {rs:reg}, {value} => {
    assert(value < (1<<4) && value >= 0)
    rd`4 @ rs`4 @ value`4  @ op`2 @ 0b10
  }
  fmt_ls  {op:op}, {rd:reg}, {rs:reg}, {value} => {
    assert(value >= (1<<4) || value < 0)
    asm { imm value } @
    rd`4 @ rs`4 @ value`4  @ op`2 @ 0b10
  }

  fmt_i12 {value}                     => {
    value`12               @  0b1101
  }
  fmt_rr  {op:op}, {rd:reg}, {rs:reg}          => {
    rd`4 @ rs`4            @ op`6 @ 0b00
  }

  nop                                => asm { fmt_rr nop, r0, r0 }
  error                              => asm { fmt_rr error, r0, r0 }
  halt                               => asm { fmt_rr halt, r0, r0 }

  return                             => asm { jump ra }

  move   {rd:reg}, {value}           => asm { fmt_ri8 movei, {rd}, value }
  move   {rd:reg}, {rs:reg}          => asm { fmt_rr move, {rd}, {rs} }

  imm    {value}                     => asm { fmt_i12 value[15:4] }
  jump   {value}                     => asm { fmt_i11 jumpi, value - pc - 1 }
  jump   {rd:reg}                    => asm { fmt_rr jump, {rd}, {rd} }
  call   {value}                     => asm { fmt_i11 calli, value - pc - 1 }

  load   {rd:reg}, [{rs:reg}, {imm}] => asm { fmt_ls load, {rd}, {rs}, imm }
  store  [{rs:reg}, {imm}], {rd:reg} => asm { fmt_ls store, {rd}, {rs}, imm }
  loadb  {rd:reg}, [{rs:reg}, {imm}] => asm { fmt_ls loadb, {rd}, {rs}, imm }
  storeb [{rs:reg}, {imm}], {rd:reg} => asm { fmt_ls storeb, {rd}, {rs}, imm }

  add    {rd:reg}, {value}           => asm { fmt_ri6 addi   , {rd}, value }
  add    {rd:reg}, {rs:reg}          => asm { fmt_rr  add    , {rd}, {rs} }
  sub    {rd:reg}, {value}           => asm { fmt_ri6 subi   , {rd}, value }
  sub    {rd:reg}, {rs:reg}          => asm { fmt_rr  sub    , {rd}, {rs} }
  addc   {rd:reg}, {value}           => asm { fmt_ri6 addci  , {rd}, value }
  addc   {rd:reg}, {rs:reg}          => asm { fmt_rr  addc   , {rd}, {rs} }
  subc   {rd:reg}, {value}           => asm { fmt_ri6 subci  , {rd}, value }
  subc   {rd:reg}, {rs:reg}          => asm { fmt_rr  subc   , {rd}, {rs} }
  xor    {rd:reg}, {value}           => asm { fmt_ri6 xori   , {rd}, value }
  xor    {rd:reg}, {rs:reg}          => asm { fmt_rr  xor    , {rd}, {rs} }
  and    {rd:reg}, {value}           => asm { fmt_ri6 andi   , {rd}, value }
  and    {rd:reg}, {rs:reg}          => asm { fmt_rr  and    , {rd}, {rs} }
  or     {rd:reg}, {value}           => asm { fmt_ri6 ori    , {rd}, value }
  or     {rd:reg}, {rs:reg}          => asm { fmt_rr  or     , {rd}, {rs} }
  shl    {rd:reg}, {value}           => asm { fmt_ri6 shli   , {rd}, value }
  shl    {rd:reg}, {rs:reg}          => asm { fmt_rr  shl    , {rd}, {rs} }
  shr    {rd:reg}, {value}           => asm { fmt_ri6 shri   , {rd}, value }
  shr    {rd:reg}, {rs:reg}          => asm { fmt_rr  shr    , {rd}, {rs} }
  asr    {rd:reg}, {value}           => asm { fmt_ri6 asri   , {rd}, value }
  asr    {rd:reg}, {rs:reg}          => asm { fmt_rr  asr    , {rd}, {rs} }
  if.eq  {rd:reg}, {value}           => asm { fmt_ri6 if.eqi , {rd}, value }
  if.eq  {rd:reg}, {rs:reg}          => asm { fmt_rr  if.eq  , {rd}, {rs} }
  if.ne  {rd:reg}, {value}           => asm { fmt_ri6 if.nei , {rd}, value }
  if.ne  {rd:reg}, {rs:reg}          => asm { fmt_rr  if.ne  , {rd}, {rs} }
  if.lt  {rd:reg}, {value}           => asm { fmt_ri6 if.lti , {rd}, value }
  if.lt  {rd:reg}, {rs:reg}          => asm { fmt_rr  if.lt  , {rd}, {rs} }
  if.ge  {rd:reg}, {value}           => asm { fmt_ri6 if.gei , {rd}, value }
  if.ge  {rd:reg}, {rs:reg}          => asm { fmt_rr  if.ge  , {rd}, {rs} }
  if.ult {rd:reg}, {value}           => asm { fmt_ri6 if.ulti, {rd}, value }
  if.ult {rd:reg}, {rs:reg}          => asm { fmt_rr  if.ult , {rd}, {rs} }
  if.uge {rd:reg}, {value}           => asm { fmt_ri6 if.ugei, {rd}, value }
  if.uge {rd:reg}, {rs:reg}          => asm { fmt_rr  if.uge , {rd}, {rs} }

  ; pseudoinstructions
  load   {rd:reg}, [{rs:reg}]  => asm { load {rd}, [{rs}, 0] }
  store  [{rs:reg}], {rd:reg}  => asm { store [{rs}, 0], {rd} }
  loadb   {rd:reg}, [{rs:reg}] => asm { loadb {rd}, [{rs}, 0] }
  storeb  [{rs:reg}], {rd:reg} => asm { storeb [{rs}, 0], {rd} }
  not {rd:reg}                 => asm { xor {rd}, -1 }
  neg {rd:reg}                 => asm {
    xor {rd}, -1
    add {rd}, 1
  }
  swap {rd:reg}, {rs:reg}      => asm {
    xor {rd}, {rs}
    xor {rs}, {rd}
    xor {rd}, {rs}
  }

  ; other names for these instructions
  if.lo  {rd:reg}, {value}     => asm { if.ult {rd}, value }
  if.lo  {rd:reg}, {rs:reg}    => asm { if.ult {rd}, {rs} }
  if.hs  {rd:reg}, {value}     => asm { if.uge {rd}, value }
  if.hs  {rd:reg}, {rs:reg}    => asm { if.uge {rd}, {rs} }

  ; variants
  if.gt  {rd:reg}, {value}     => asm { if.ge  {rd}, value+1 }
  if.gt  {rd:reg}, {rs:reg}    => asm { if.lt  {rs}, {rd} }
  if.le  {rd:reg}, {value}     => asm { if.lt  {rd}, value-1 }
  if.le  {rd:reg}, {rs:reg}    => asm { if.ge  {rs}, {rd} }
  if.ugt  {rd:reg}, {value}    => asm { if.uge {rd}, value+1 }
  if.ugt  {rd:reg}, {rs:reg}   => asm { if.ult {rs}, {rd} }
  if.ule  {rd:reg}, {value}    => {
    assert(value == 0)
    asm { if.eq {rd}, 0 }
  }
  if.ule  {rd:reg}, {value}    => {
    assert(value > 0)
    asm { if.ult {rd}, value-1 }
  }
  if.ule  {rd:reg}, {rs:reg}   => asm { if.uge {rs}, {rd} }


  sxt {rd:reg} => asm {
    shl {rd}, 8
    asr {rd}, 8
  }
}

; code bank is the main program memory bank
#bankdef code
{
  #bits 16
  #addr 0x0000
  #size 0x10000
  #outp 0
}

; data is the bank where strings, constants and pre-initialized
; values goes.
#bankdef data
{
  #bits 16
  #addr 0x8000
  #size 0x4000
  #outp 0x10000*16
}

; bss is the main data memory bank for uninitialized variables
#bankdef bss
{
  #bits 16
  #addr 0x0000
  #size 0x8000
}

#bank code
