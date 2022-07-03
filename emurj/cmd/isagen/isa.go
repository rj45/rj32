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
	"s":      {26, 26},
	"i1":     {25, 16},
	"i2":     {15, 14},
	"i3":     {13, 6},
	"i4":     {31, 30},
	"x1":     {13, 11},
	"x2":     {29, 27},
	"x3":     {15, 15},
	"opcode": {5, 0},
}

var fmtOrder = []string{"N", "R", "I", "B", "J", "U"}

var fmtFields = map[string][]string{
	"N": {"rd", "s", "i1", "rs2", "rs1", "opcode"},
	"R": {"rd", "s", "i1", "rs2", "rs1", "opcode"},
	"I": {"rd", "s", "i1", "i2", "x1", "rs1", "opcode"},
	"B": {"i4", "x2", "s", "i1", "rs2", "rs1", "opcode"},
	"J": {"rd", "s", "i1", "i2", "i3", "opcode"},
	"U": {"rd", "s", "i1", "x3", "i2", "i3", "opcode"},
}

var cpudefOperands = map[string][][2]string{
	"N": {},
	"R": {{"rd", "reg"}, {"rs1", "reg"}, {"rs2", "reg"}},
	"I": {{"rd", "reg"}, {"rs1", "reg"}, {"imm", "s13"}},
	"B": {{"rs1", "reg"}, {"rs2", "reg"}, {"imm", "s13"}},
	"J": {{"rd", "reg"}, {"imm", "s21"}},
	"U": {{"rd", "reg"}, {"imm", "s32"}},
}

var goOperandTypes = map[string]string{
	"opcode": "Opcode",
	"rd":     "Reg",
	"rs1":    "Reg",
	"rs2":    "Reg",
	"imm":    "int32",
}

var operandNames = [...]string{
	"opcode", "rd", "rs1", "rs2", "imm",
}

var fieldOperands = map[string]string{
	"rd":     "rd",
	"rs1":    "rs1",
	"rs2":    "rs2",
	"s":      "imm",
	"i1":     "imm",
	"i2":     "imm",
	"i3":     "imm",
	"i4":     "imm",
	"opcode": "opcode",
}

var immBits = map[string]map[string][2]int{
	"I": {"s": {12, 12}, "i1": {11, 2}, "i2": {1, 0}},
	"B": {"s": {12, 12}, "i1": {11, 2}, "i4": {1, 0}},
	"J": {"s": {20, 20}, "i3": {19, 12}, "i1": {11, 2}, "i2": {1, 0}},
	"U": {"s": {31, 31}, "i1": {30, 21}, "i2": {20, 20}, "i3": {19, 12}},
}

type opcodeMatcher struct {
	mask, match uint
	fmt         string
}

var opcodeFmtMatchers = []opcodeMatcher{
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
