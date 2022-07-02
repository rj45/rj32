package main

const registers = 32

const abiNames = "zero ra sp gp tp t0 t1 t2 s0 s1 " +
	"a0 a1 a2 a3 a4 a5 a6 a7 s2 s3 s4 s5 s6 s7 s8 " +
	"s9 s10 s11 t3 t4 t5 t6"

var opmap = [8][8]string{
	{"add", "slt", "sltu", "and", "or", "xor", "srl", "sll"},
	{"sub", "mul", "mulh", "mulhu", "", "", "sra", ""},
	{"ebreak", "ecall", "sret", "mret", "", "", "", ""},
	{"csrrw", "csrrs", "csrrc", "", "sfence", "fence", "fence.i", "wfi"},
	{"addi", "slti", "sltiu", "andi", "ori", "xori", "srli", "slli"},
	{"lw", "lr", "lh", "lb", "lhu", "lbu", "srai", "jalr"},
	{"beq", "bne", "blt", "bge", "bltu", "bgeu", "", ""},
	{"sw", "sc", "sh", "sb", "lui", "auipc", "", "jal"},
}

var fieldBits = map[string][2]int{
	"rd":     {31, 27},
	"rs1":    {10, 6},
	"rs2":    {15, 11},
	"i1":     {26, 16},
	"i2":     {15, 14},
	"i3":     {13, 6},
	"i4":     {31, 27},
	"x1":     {13, 11},
	"x2":     {29, 27},
	"x3":     {15, 15},
	"opcode": {5, 0},
}

var fmtOrder = []string{"N", "R", "I", "B", "J", "U"}

var fmtFields = map[string][]string{
	"N": {"rd", "i1", "rs2", "rs1", "opcode"},
	"R": {"rd", "i1", "rs2", "rs1", "opcode"},
	"I": {"rd", "i1", "i2", "x1", "rs1", "opcode"},
	"B": {"i4", "x2", "i1", "rs2", "rs1", "opcode"},
	"J": {"rd", "i1", "i2", "i3", "opcode"},
	"U": {"rd", "i1", "x3", "i2", "i3", "opcode"},
}

var operands = map[string][][2]string{
	"N": {},
	"R": {{"rd", "reg"}, {"rs1", "reg"}, {"rs2", "reg"}},
	"I": {{"rd", "reg"}, {"rs1", "reg"}, {"val", "s13"}},
	"B": {{"rs1", "reg"}, {"rs2", "reg"}, {"val", "s13"}},
	"J": {{"rd", "reg"}, {"val", "s21"}},
	"U": {{"rd", "reg"}, {"val", "s32"}},
}

var fieldOperands = map[string]string{
	"rd":     "rd",
	"rs1":    "rs1",
	"rs2":    "rs2",
	"i1":     "val",
	"i2":     "val",
	"i3":     "val",
	"i4":     "val",
	"opcode": "opcode",
}

var immBits = map[string]map[string][2]int{
	"I": {"i1": {12, 2}, "i2": {1, 0}},
	"B": {"i1": {12, 2}, "i4": {1, 0}},
	"J": {"i3": {20, 12}, "i1": {11, 2}, "i2": {1, 0}},
	"U": {"i1": {31, 21}, "i2": {20, 20}, "i3": {19, 12}},
}

var opcodeFmtMatchers = []struct {
	mask, match uint
	fmt         string
}{
	{0b111_111, 0b111_111, "J"},
	{0b111_100, 0b111_100, "U"},
	{0b110_000, 0b110_000, "B"},
	{0b100_000, 0b100_000, "I"},
	{0b111_100, 0b010_000, "N"},
	{0b111_100, 0b011_100, "N"},
	{0b000_000, 0b000_000, "R"},
}

func opcodeIsPcRelative(quad, i int) bool {
	if opmap[quad][i][0] == 'b' ||
		opmap[quad][i][0] == 'j' ||
		opmap[quad][i] == "auipc" {
		return true
	}
	return false
}
