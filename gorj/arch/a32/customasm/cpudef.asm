; Copyright (c) 2021 Mathis "Artentus" Rech; MIT Licensed, similar to LICENSE.
; https://pastebin.com/4TsiQeHU

#bits 8

#subruledef reg
{
    r{i: u5} => i
    zero => 0`5
    ra => 1`5
    bp => 2`5
    sp => 3`5
    a{i: u3} => (i`5 + 4)`5
    t{i: u5} => {
        assert(i < 10)
        (i + 12)`5
    }
    s{i: u5} => {
        assert(i < 10)
        (i + 22)`5
    }
}

#subruledef imm
{
    {v: i32} => {
        assert((v < 8192) && (v >= -8192))
        0b0 @ v`14
    }
    {v: i32} => {
        assert((v >= 8192) || (v < -8192))
        0`14 @ (v >> 14)`18 @ 0b1 @ v`14
    }
}

#subruledef rel
{
    {v: i32} => {
        offset = v - $ - 4
        assert((offset < 8192) && (offset >= -8192))
        0b0 @ offset`14
    }
    {v: i32} => {
        offset = v - $ - 4
        assert((offset >= 8192) || (offset < -8192))
        0`14 @ (offset >> 14)`18 @ 0b1 @ offset`14
    }
}

; real instructions
#ruledef
{
    NOP => le(0`10 @ 0`5 @ 0`5 @ 0`5 @ 0x0 @ 0b000)
    BRK => le(0`10 @ 0`5 @ 0`5 @ 0`5 @ 0x1 @ 0b000)
    HLT => le(0`10 @ 0`5 @ 0`5 @ 0`5 @ 0x2 @ 0b000)
    ERR => le(0`10 @ 0`5 @ 0`5 @ 0`5 @ 0x3 @ 0b000)

    ADD  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x1 @ 0b001)
    ADDC {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x2 @ 0b001)
    SUB  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x3 @ 0b001)
    SUBB {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x4 @ 0b001)
    AND  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x5 @ 0b001)
    OR   {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x6 @ 0b001)
    XOR  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x7 @ 0b001)
    SHL  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x8 @ 0b001)
    ASR  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0x9 @ 0b001)
    LSR  {d: reg}, {l: reg}, {r: reg} => le(0`10 @ {r} @ {l} @ {d} @ 0xA @ 0b001)

    ADD  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x1 @ 0b010)
    ADDC {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x2 @ 0b010)
    SUB  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x3 @ 0b010)
    SUBB {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x4 @ 0b010)
    AND  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x5 @ 0b010)
    OR   {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x6 @ 0b010)
    XOR  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x7 @ 0b010)
    SHL  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x8 @ 0b010)
    ASR  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0x9 @ 0b010)
    LSR  {d: reg}, {l: reg}, {v: imm} => le({v} @ {l} @ {d} @ 0xA @ 0b010)

    ADD  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x1 @ 0b011)
    ADDC {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x2 @ 0b011)
    SUB  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x3 @ 0b011)
    SUBB {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x4 @ 0b011)
    AND  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x5 @ 0b011)
    OR   {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x6 @ 0b011)
    XOR  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x7 @ 0b011)
    SHL  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x8 @ 0b011)
    ASR  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0x9 @ 0b011)
    LSR  {d: reg}, {v: imm}, {r: reg} => le({v} @ {r} @ {d} @ 0xA @ 0b011)

    LD   {d: reg}, [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ {d} @ 0x0 @ 0b100)
    LD   {d: reg}, [{s: reg} + {v: imm}] => le(       {v} @ {s} @ {d} @ 0x1 @ 0b100)
    ST   [{d: reg} + {o: reg}], {s: reg} => le(0`10 @ {o} @ {s} @ {d} @ 0x2 @ 0b100)
    ST   [{d: reg} + {v: imm}], {s: reg} => le(       {v} @ {s} @ {d} @ 0x3 @ 0b100)
    LD8  {d: reg}, [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ {d} @ 0x4 @ 0b100)
    LD8  {d: reg}, [{s: reg} + {v: imm}] => le(       {v} @ {s} @ {d} @ 0x5 @ 0b100)
    ST8  [{d: reg} + {o: reg}], {s: reg} => le(0`10 @ {o} @ {s} @ {d} @ 0x6 @ 0b100)
    ST8  [{d: reg} + {v: imm}], {s: reg} => le(       {v} @ {s} @ {d} @ 0x7 @ 0b100)
    LD16 {d: reg}, [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ {d} @ 0x8 @ 0b100)
    LD16 {d: reg}, [{s: reg} + {v: imm}] => le(       {v} @ {s} @ {d} @ 0x9 @ 0b100)
    ST16 [{d: reg} + {o: reg}], {s: reg} => le(0`10 @ {o} @ {s} @ {d} @ 0xA @ 0b100)
    ST16 [{d: reg} + {v: imm}], {s: reg} => le(       {v} @ {s} @ {d} @ 0xB @ 0b100)

    JP.C    {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x1 @ 0b101)
    JP.Z    {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x2 @ 0b101)
    JP.S    {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x3 @ 0b101)
    JP.O    {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x4 @ 0b101)
    JP.NC   {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x5 @ 0b101)
    JP.NZ   {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x6 @ 0b101)
    JP.NS   {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x7 @ 0b101)
    JP.NO   {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x8 @ 0b101)
    JP.U.LE {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0x9 @ 0b101)
    JP.U.G  {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0xA @ 0b101)
    JP.S.L  {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0xB @ 0b101)
    JP.S.GE {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0xC @ 0b101)
    JP.S.LE {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0xD @ 0b101)
    JP.S.G  {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0xE @ 0b101)
    JMP     {s: reg} + {o: reg} => le(0`10 @ {o} @ {s} @ 0`5 @ 0xF @ 0b101)

    JP.C    [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x1 @ 0b101)
    JP.Z    [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x2 @ 0b101)
    JP.S    [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x3 @ 0b101)
    JP.O    [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x4 @ 0b101)
    JP.NC   [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x5 @ 0b101)
    JP.NZ   [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x6 @ 0b101)
    JP.NS   [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x7 @ 0b101)
    JP.NO   [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x8 @ 0b101)
    JP.U.LE [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0x9 @ 0b101)
    JP.U.G  [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0xA @ 0b101)
    JP.S.L  [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0xB @ 0b101)
    JP.S.GE [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0xC @ 0b101)
    JP.S.LE [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0xD @ 0b101)
    JP.S.G  [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0xE @ 0b101)
    JMP     [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ 1`5 @ 0xF @ 0b101)

    JP.C    {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x1 @ 0b101)
    JP.Z    {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x2 @ 0b101)
    JP.S    {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x3 @ 0b101)
    JP.O    {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x4 @ 0b101)
    JP.NC   {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x5 @ 0b101)
    JP.NZ   {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x6 @ 0b101)
    JP.NS   {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x7 @ 0b101)
    JP.NO   {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x8 @ 0b101)
    JP.U.LE {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0x9 @ 0b101)
    JP.U.G  {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0xA @ 0b101)
    JP.S.L  {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0xB @ 0b101)
    JP.S.GE {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0xC @ 0b101)
    JP.S.LE {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0xD @ 0b101)
    JP.S.G  {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0xE @ 0b101)
    JMP     {s: reg} + {v: imm} => le({v} @ {s} @ 2`5 @ 0xF @ 0b101)

    JP.C    [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x1 @ 0b101)
    JP.Z    [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x2 @ 0b101)
    JP.S    [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x3 @ 0b101)
    JP.O    [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x4 @ 0b101)
    JP.NC   [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x5 @ 0b101)
    JP.NZ   [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x6 @ 0b101)
    JP.NS   [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x7 @ 0b101)
    JP.NO   [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x8 @ 0b101)
    JP.U.LE [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0x9 @ 0b101)
    JP.U.G  [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0xA @ 0b101)
    JP.S.L  [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0xB @ 0b101)
    JP.S.GE [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0xC @ 0b101)
    JP.S.LE [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0xD @ 0b101)
    JP.S.G  [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0xE @ 0b101)
    JMP     [{s: reg} + {v: imm}] => le({v} @ {s} @ 3`5 @ 0xF @ 0b101)

    BR.C    {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x1 @ 0b101)
    BR.Z    {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x2 @ 0b101)
    BR.S    {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x3 @ 0b101)
    BR.O    {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x4 @ 0b101)
    BR.NC   {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x5 @ 0b101)
    BR.NZ   {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x6 @ 0b101)
    BR.NS   {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x7 @ 0b101)
    BR.NO   {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x8 @ 0b101)
    BR.U.LE {v: rel} => le({v} @ 0`5 @ 4`5 @ 0x9 @ 0b101)
    BR.U.G  {v: rel} => le({v} @ 0`5 @ 4`5 @ 0xA @ 0b101)
    BR.S.L  {v: rel} => le({v} @ 0`5 @ 4`5 @ 0xB @ 0b101)
    BR.S.GE {v: rel} => le({v} @ 0`5 @ 4`5 @ 0xC @ 0b101)
    BR.S.LE {v: rel} => le({v} @ 0`5 @ 4`5 @ 0xD @ 0b101)
    BR.S.G  {v: rel} => le({v} @ 0`5 @ 4`5 @ 0xE @ 0b101)
    BRA     {v: rel} => le({v} @ 0`5 @ 4`5 @ 0xF @ 0b101)

    IN  {d: reg}, [{s: reg} + {o: reg}] => le(0`10 @ {o} @ {s} @ {d} @ 0x0 @ 0b110)
    IN  {d: reg}, [{s: reg} + {v: imm}] => le(       {v} @ {s} @ {d} @ 0x1 @ 0b110)
    OUT [{d: reg} + {o: reg}], {s: reg} => le(0`10 @ {o} @ {s} @ {d} @ 0x2 @ 0b110)
    OUT [{d: reg} + {v: imm}], {s: reg} => le(       {v} @ {s} @ {d} @ 0x3 @ 0b110)

    SYS  => le(0`10 @ 0`5 @ 0`5 @ 0`5 @ 0x0 @ 0b111)
    CLRK => le(0`10 @ 0`5 @ 0`5 @ 0`5 @ 0x1 @ 0b111)
}

; aliases
#ruledef
{
    MOV {d: reg}, {s: reg} => asm { OR {d}, {s}, zero }
    LD  {d: reg}, {v: i32} => asm { OR {d},  v , zero }

    CMP {l: reg}, {r: reg} => asm { SUB zero, {l}, {r}}
    CMP {l: reg}, {v: i32} => asm { SUB zero, {l},  v }
    CMP {v: i32}, {r: reg} => asm { SUB zero,  v , {r}}

    BIT {l: reg}, {r: reg} => asm { AND zero, {l}, {r}}
    BIT {l: reg}, {v: i32} => asm { AND zero, {l},  v }
    BIT {v: i32}, {r: reg} => asm { AND zero,  v , {r}}

    TEST {s: reg} => asm { OR zero, {s}, zero }

    INC  {d: reg} => asm { ADD  {d}, {d}, 1 }
    INCC {d: reg} => asm { ADDC {d}, {d}, zero }
    DEC  {d: reg} => asm { SUB  {d}, {d}, 1 }
    DECB {d: reg} => asm { SUBB {d}, {d}, zero }

    NEG  {d: reg}, {s: reg} => asm { SUB  {d}, zero, {s} }
    NEGB {d: reg}, {s: reg} => asm { SUBB {d}, zero, {s} }

    NOT {d: reg}, {s: reg} => asm { XOR {d}, {s}, -1 }

    LD   {d: reg}, [{s: reg}] => asm { LD   {d}, [{s} + zero] }
    LD   {d: reg}, [{v: i32}] => asm { LD   {d}, [zero + v] }
    ST   [{d: reg}], {s: reg} => asm { ST   [{d} + zero], {s} }
    ST   [{v: i32}], {s: reg} => asm { ST   [zero + v], {s} }
    LD8  {d: reg}, [{s: reg}] => asm { LD8  {d}, [{s} + zero] }
    LD8  {d: reg}, [{v: i32}] => asm { LD8  {d}, [zero + v] }
    ST8  [{d: reg}], {s: reg} => asm { ST8  [{d} + zero], {s} }
    ST8  [{v: i32}], {s: reg} => asm { ST8  [zero + v], {s} }
    LD16 {d: reg}, [{s: reg}] => asm { LD16 {d}, [{s} + zero] }
    LD16 {d: reg}, [{v: i32}] => asm { LD16 {d}, [zero + v] }
    ST16 [{d: reg}], {s: reg} => asm { ST16 [{d} + zero], {s} }
    ST16 [{v: i32}], {s: reg} => asm { ST16 [zero + v], {s} }

    JP.EQ   {s: reg} + {o: reg} => asm { JP.Z  {s} + {o} }
    JP.NEQ  {s: reg} + {o: reg} => asm { JP.NZ {s} + {o} }
    JP.U.L  {s: reg} + {o: reg} => asm { JP.NC {s} + {o} }
    JP.U.GE {s: reg} + {o: reg} => asm { JP.C  {s} + {o} }

    JP.EQ   [{s: reg} + {o: reg}] => asm { JP.Z  [{s} + {o}] }
    JP.NEQ  [{s: reg} + {o: reg}] => asm { JP.NZ [{s} + {o}] }
    JP.U.L  [{s: reg} + {o: reg}] => asm { JP.NC [{s} + {o}] }
    JP.U.GE [{s: reg} + {o: reg}] => asm { JP.C  [{s} + {o}] }

    JP.EQ   {s: reg} + {v: u32} => asm { JP.Z  {s} + v }
    JP.NEQ  {s: reg} + {v: u32} => asm { JP.NZ {s} + v }
    JP.U.L  {s: reg} + {v: u32} => asm { JP.NC {s} + v }
    JP.U.GE {s: reg} + {v: u32} => asm { JP.C  {s} + v }

    JP.EQ   [{s: reg} + {v: u32}] => asm { JP.Z  [{s} + v] }
    JP.NEQ  [{s: reg} + {v: u32}] => asm { JP.NZ [{s} + v] }
    JP.U.L  [{s: reg} + {v: u32}] => asm { JP.NC [{s} + v] }
    JP.U.GE [{s: reg} + {v: u32}] => asm { JP.C  [{s} + v] }

    BR.EQ   {v: u32} => asm { BR.Z  v }
    BR.NEQ  {v: u32} => asm { BR.NZ v }
    BR.U.L  {v: u32} => asm { BR.NC v }
    BR.U.GE {v: u32} => asm { BR.C  v }

    JP.C    {s: reg} => asm { JP.C    {s} + zero }
    JP.Z    {s: reg} => asm { JP.Z    {s} + zero }
    JP.S    {s: reg} => asm { JP.S    {s} + zero }
    JP.O    {s: reg} => asm { JP.O    {s} + zero }
    JP.NC   {s: reg} => asm { JP.NC   {s} + zero }
    JP.NZ   {s: reg} => asm { JP.NZ   {s} + zero }
    JP.NS   {s: reg} => asm { JP.NS   {s} + zero }
    JP.NO   {s: reg} => asm { JP.NO   {s} + zero }
    JP.EQ   {s: reg} => asm { JP.Z    {s} + zero }
    JP.NEQ  {s: reg} => asm { JP.NZ   {s} + zero }
    JP.U.L  {s: reg} => asm { JP.NC   {s} + zero }
    JP.U.GE {s: reg} => asm { JP.C    {s} + zero }
    JP.U.LE {s: reg} => asm { JP.U.LE {s} + zero }
    JP.U.G  {s: reg} => asm { JP.U.G  {s} + zero }
    JP.S.L  {s: reg} => asm { JP.S.L  {s} + zero }
    JP.S.GE {s: reg} => asm { JP.S.GE {s} + zero }
    JP.S.LE {s: reg} => asm { JP.S.LE {s} + zero }
    JP.S.G  {s: reg} => asm { JP.S.G  {s} + zero }
    JMP     {s: reg} => asm { JMP     {s} + zero }

    JP.C    [{s: reg}] => asm { JP.C    [{s} + zero] }
    JP.Z    [{s: reg}] => asm { JP.Z    [{s} + zero] }
    JP.S    [{s: reg}] => asm { JP.S    [{s} + zero] }
    JP.O    [{s: reg}] => asm { JP.O    [{s} + zero] }
    JP.NC   [{s: reg}] => asm { JP.NC   [{s} + zero] }
    JP.NZ   [{s: reg}] => asm { JP.NZ   [{s} + zero] }
    JP.NS   [{s: reg}] => asm { JP.NS   [{s} + zero] }
    JP.NO   [{s: reg}] => asm { JP.NO   [{s} + zero] }
    JP.EQ   [{s: reg}] => asm { JP.Z    [{s} + zero] }
    JP.NEQ  [{s: reg}] => asm { JP.NZ   [{s} + zero] }
    JP.U.L  [{s: reg}] => asm { JP.NC   [{s} + zero] }
    JP.U.GE [{s: reg}] => asm { JP.C    [{s} + zero] }
    JP.U.LE [{s: reg}] => asm { JP.U.LE [{s} + zero] }
    JP.U.G  [{s: reg}] => asm { JP.U.G  [{s} + zero] }
    JP.S.L  [{s: reg}] => asm { JP.S.L  [{s} + zero] }
    JP.S.GE [{s: reg}] => asm { JP.S.GE [{s} + zero] }
    JP.S.LE [{s: reg}] => asm { JP.S.LE [{s} + zero] }
    JP.S.G  [{s: reg}] => asm { JP.S.G  [{s} + zero] }
    JMP     [{s: reg}] => asm { JMP     [{s} + zero] }

    JP.C    {v: u32} => asm { JP.C    zero + v }
    JP.Z    {v: u32} => asm { JP.Z    zero + v }
    JP.S    {v: u32} => asm { JP.S    zero + v }
    JP.O    {v: u32} => asm { JP.O    zero + v }
    JP.NC   {v: u32} => asm { JP.NC   zero + v }
    JP.NZ   {v: u32} => asm { JP.NZ   zero + v }
    JP.NS   {v: u32} => asm { JP.NS   zero + v }
    JP.NO   {v: u32} => asm { JP.NO   zero + v }
    JP.EQ   {v: u32} => asm { JP.Z    zero + v }
    JP.NEQ  {v: u32} => asm { JP.NZ   zero + v }
    JP.U.L  {v: u32} => asm { JP.NC   zero + v }
    JP.U.GE {v: u32} => asm { JP.C    zero + v }
    JP.U.LE {v: u32} => asm { JP.U.LE zero + v }
    JP.U.G  {v: u32} => asm { JP.U.G  zero + v }
    JP.S.L  {v: u32} => asm { JP.S.L  zero + v }
    JP.S.GE {v: u32} => asm { JP.S.GE zero + v }
    JP.S.LE {v: u32} => asm { JP.S.LE zero + v }
    JP.S.G  {v: u32} => asm { JP.S.G  zero + v }
    JMP     {v: u32} => asm { JMP     zero + v }

    JP.C    [{v: u32}] => asm { JP.C    [zero + v] }
    JP.Z    [{v: u32}] => asm { JP.Z    [zero + v] }
    JP.S    [{v: u32}] => asm { JP.S    [zero + v] }
    JP.O    [{v: u32}] => asm { JP.O    [zero + v] }
    JP.NC   [{v: u32}] => asm { JP.NC   [zero + v] }
    JP.NZ   [{v: u32}] => asm { JP.NZ   [zero + v] }
    JP.NS   [{v: u32}] => asm { JP.NS   [zero + v] }
    JP.NO   [{v: u32}] => asm { JP.NO   [zero + v] }
    JP.EQ   [{v: u32}] => asm { JP.Z    [zero + v] }
    JP.NEQ  [{v: u32}] => asm { JP.NZ   [zero + v] }
    JP.U.L  [{v: u32}] => asm { JP.NC   [zero + v] }
    JP.U.GE [{v: u32}] => asm { JP.C    [zero + v] }
    JP.U.LE [{v: u32}] => asm { JP.U.LE [zero + v] }
    JP.U.G  [{v: u32}] => asm { JP.U.G  [zero + v] }
    JP.S.L  [{v: u32}] => asm { JP.S.L  [zero + v] }
    JP.S.GE [{v: u32}] => asm { JP.S.GE [zero + v] }
    JP.S.LE [{v: u32}] => asm { JP.S.LE [zero + v] }
    JP.S.G  [{v: u32}] => asm { JP.S.G  [zero + v] }
    JMP     [{v: u32}] => asm { JMP     [zero + v] }

    IN  {d: reg}, [{s: reg}] => asm { IN   {d}, [{s} + zero] }
    IN  {d: reg}, [{v: i32}] => asm { IN   {d}, [zero + v] }
    OUT [{d: reg}], {s: reg} => asm { OUT  [{d} + zero], {s} }
    OUT [{v: i32}], {s: reg} => asm { OUT  [zero + v], {s} }
}

; macros
#ruledef
{
    PUSH {s: reg} => asm {
        ST [sp], {s}
        SUB sp, sp, 4
    }

    POP {d: reg} => asm {
        ADD sp, sp, 4
        LD {d}, [sp]
    }

    PUSH8 {s: reg} => asm {
        ST8 [sp], {s}
        SUB sp, sp, 4
    }

    POP8 {d: reg} => asm {
        ADD sp, sp, 4
        LD8 {d}, [sp]
    }

    PUSH16 {s: reg} => asm {
        ST16 [sp], {s}
        SUB sp, sp, 4
    }

    POP16 {d: reg} => asm {
        ADD sp, sp, 4
        LD16 {d}, [sp]
    }

    CALL {s: reg} => {
        addr = $ + 12
        assert(addr < 8192)

        asm {
            MOV bp, sp
            LD ra, addr`14
            JMP {s}
        }
    }
    CALL {s: reg} => {
        addr = $ + 12
        assert(addr >= 8192)

        asm {
            MOV bp, sp
            LD ra, (addr + 4)
            JMP {s}
        }
    }

    CALL [{s: reg}] => {
        addr = $ + 12
        assert(addr < 8192)

        asm {
            MOV bp, sp
            LD ra, addr`14
            JMP [{s}]
        }
    }
    CALL [{s: reg}] => {
        addr = $ + 12
        assert(addr >= 8192)

        asm {
            MOV bp, sp
            LD ra, (addr + 4)
            JMP [{s}]
        }
    }

    CALL {v: u32} => {
        assert(v < 8192)
        addr = $ + 12
        assert(addr < 8192)

        asm {
            MOV bp, sp
            LD ra, addr`14
            JMP v`14
        }
    }
    CALL {v: u32} => {
        assert(v < 8192)
        addr = $ + 12
        assert(addr >= 8192)

        asm {
            MOV bp, sp
            LD ra, (addr + 4)
            JMP v`14
        }
    }
    CALL {v: u32} => {
        assert(v >= 8192)
        addr = $ + 16
        assert(addr < 8192)

        asm {
            MOV bp, sp
            LD ra, addr`14
            JMP v
        }
    }
    CALL {v: u32} => {
        assert(v >= 8192)
        addr = $ + 16
        assert(addr >= 8192)

        asm {
            MOV bp, sp
            LD ra, (addr + 4)
            JMP v
        }
    }

    CALL [{v: u32}] => {
        assert(v < 8192)
        addr = $ + 12
        assert(addr < 8192)

        asm {
            MOV bp, sp
            LD ra, addr`14
            JMP [v`14]
        }
    }
    CALL [{v: u32}] => {
        assert(v < 8192)
        addr = $ + 12
        assert(addr >= 8192)

        asm {
            MOV bp, sp
            LD ra, (addr + 4)
            JMP [v`14]
        }
    }
    CALL [{v: u32}] => {
        assert(v >= 8192)
        addr = $ + 16
        assert(addr < 8192)

        asm {
            MOV bp, sp
            LD ra, addr`14
            JMP [v]
        }
    }
    CALL [{v: u32}] => {
        assert(v >= 8192)
        addr = $ + 16
        assert(addr >= 8192)

        asm {
            MOV bp, sp
            LD ra, (addr + 4)
            JMP [v]
        }
    }

    RET {v: u12} => asm {
        ADD sp, bp, (v`14 << 2)`14
        JMP ra
    }

    CALLS => {
        addr = $ + 8
        assert(addr < 8192)

        asm {
            LD ra, addr`14
            SYS
        }
    }
    CALLS => {
        addr = $ + 8
        assert(addr >= 8192)

        asm {
            LD ra, (addr + 4)
            SYS
        }
    }

    RETS => asm {
        CLRK
        JMP ra
    }
}
