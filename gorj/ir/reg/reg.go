package reg

import (
	"math/bits"
	"strings"
)

type Arch interface {
	RegNames() []string
	SavedRegs() []Reg
	TempRegs() []Reg
	ArgRegs() []Reg
	SpecialRegs() map[string]Reg
}

func SetArch(arch Arch) {
	names = arch.RegNames()
	SavedRegs = arch.SavedRegs()
	ArgRegs = arch.ArgRegs()
	TempRegs = arch.TempRegs()

	RevSavedRegs = make([]Reg, len(SavedRegs))
	savedRegMask = 0
	for i, reg := range SavedRegs {
		RevSavedRegs[len(SavedRegs)-i-1] = reg
		savedRegMask |= reg
	}

	argRegMask = 0
	for _, reg := range ArgRegs {
		argRegMask |= reg
	}

	spec := arch.SpecialRegs()
	RA = spec["RA"]
	GP = spec["GP"]
	SP = spec["SP"]
	FP = spec["FP"]
}

var names []string

type Reg uint

const (
	None Reg = 0
)

var SavedRegs []Reg
var RevSavedRegs []Reg
var TempRegs []Reg
var ArgRegs []Reg

var savedRegMask Reg
var argRegMask Reg

var SP Reg
var FP Reg
var GP Reg
var RA Reg

func FromRegNum(num int) Reg {
	return 1 << num
}

func (reg Reg) String() string {
	if reg == None {
		return "none"
	}

	if reg.IsMany() {
		regs := reg.RegNumbers()
		strs := make([]string, len(regs))
		for i, num := range regs {
			strs[i] = names[num]
		}
		return strings.Join(strs, ",")
	}

	return names[reg.RegNumber()]
}

func (reg Reg) RegNumber() int {
	if reg == None {
		return -1
	}
	return bits.TrailingZeros(uint(reg))
}

func (reg Reg) NumRegs() int {
	return bits.OnesCount(uint(reg))
}

func (reg Reg) IsMany() bool {
	return bits.OnesCount(uint(reg)) > 1
}

func (reg Reg) RegNumbers() []int {
	count := bits.OnesCount(uint(reg))
	if count == 0 {
		return nil
	}
	regs := make([]int, count)
	for i := range regs {
		num := bits.TrailingZeros(uint(reg))
		regs[i] = num
		reg &^= 1 << num
	}
	return regs
}

func (reg Reg) IsArgReg() bool {
	return reg&argRegMask != 0
}

func (reg Reg) IsSpecialReg() bool {
	return reg == GP || reg == SP || reg == RA
}

func (reg Reg) IsSavedReg() bool {
	return reg&savedRegMask != 0
}

func (reg Reg) CanAffinity() bool {
	return !reg.IsSpecialReg()
}
