// Code generated by "enumer -transform lower -trimprefix Fmt -type Op,Fmt ."; DO NOT EDIT.

package rj32

import (
	"fmt"
	"strings"
)

const _OpName = "nopretserrorhaltrcsrwcsrmoveloadcjumpimmcallimm2loadstoreloadbstorebaddsubaddcsubcxorandorshlshrasrifeqifneifltifgeifultifuge"

var _OpIndex = [...]uint8{0, 3, 7, 12, 16, 20, 24, 28, 33, 37, 40, 44, 48, 52, 57, 62, 68, 71, 74, 78, 82, 85, 88, 90, 93, 96, 99, 103, 107, 111, 115, 120, 125}

const _OpLowerName = "nopretserrorhaltrcsrwcsrmoveloadcjumpimmcallimm2loadstoreloadbstorebaddsubaddcsubcxorandorshlshrasrifeqifneifltifgeifultifuge"

func (i Op) String() string {
	if i < 0 || i >= Op(len(_OpIndex)-1) {
		return fmt.Sprintf("Op(%d)", i)
	}
	return _OpName[_OpIndex[i]:_OpIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OpNoOp() {
	var x [1]struct{}
	_ = x[Nop-(0)]
	_ = x[Rets-(1)]
	_ = x[Error-(2)]
	_ = x[Halt-(3)]
	_ = x[Rcsr-(4)]
	_ = x[Wcsr-(5)]
	_ = x[Move-(6)]
	_ = x[Loadc-(7)]
	_ = x[Jump-(8)]
	_ = x[Imm-(9)]
	_ = x[Call-(10)]
	_ = x[Imm2-(11)]
	_ = x[Load-(12)]
	_ = x[Store-(13)]
	_ = x[Loadb-(14)]
	_ = x[Storeb-(15)]
	_ = x[Add-(16)]
	_ = x[Sub-(17)]
	_ = x[Addc-(18)]
	_ = x[Subc-(19)]
	_ = x[Xor-(20)]
	_ = x[And-(21)]
	_ = x[Or-(22)]
	_ = x[Shl-(23)]
	_ = x[Shr-(24)]
	_ = x[Asr-(25)]
	_ = x[IfEq-(26)]
	_ = x[IfNe-(27)]
	_ = x[IfLt-(28)]
	_ = x[IfGe-(29)]
	_ = x[IfUlt-(30)]
	_ = x[IfUge-(31)]
}

var _OpValues = []Op{Nop, Rets, Error, Halt, Rcsr, Wcsr, Move, Loadc, Jump, Imm, Call, Imm2, Load, Store, Loadb, Storeb, Add, Sub, Addc, Subc, Xor, And, Or, Shl, Shr, Asr, IfEq, IfNe, IfLt, IfGe, IfUlt, IfUge}

var _OpNameToValueMap = map[string]Op{
	_OpName[0:3]:          Nop,
	_OpLowerName[0:3]:     Nop,
	_OpName[3:7]:          Rets,
	_OpLowerName[3:7]:     Rets,
	_OpName[7:12]:         Error,
	_OpLowerName[7:12]:    Error,
	_OpName[12:16]:        Halt,
	_OpLowerName[12:16]:   Halt,
	_OpName[16:20]:        Rcsr,
	_OpLowerName[16:20]:   Rcsr,
	_OpName[20:24]:        Wcsr,
	_OpLowerName[20:24]:   Wcsr,
	_OpName[24:28]:        Move,
	_OpLowerName[24:28]:   Move,
	_OpName[28:33]:        Loadc,
	_OpLowerName[28:33]:   Loadc,
	_OpName[33:37]:        Jump,
	_OpLowerName[33:37]:   Jump,
	_OpName[37:40]:        Imm,
	_OpLowerName[37:40]:   Imm,
	_OpName[40:44]:        Call,
	_OpLowerName[40:44]:   Call,
	_OpName[44:48]:        Imm2,
	_OpLowerName[44:48]:   Imm2,
	_OpName[48:52]:        Load,
	_OpLowerName[48:52]:   Load,
	_OpName[52:57]:        Store,
	_OpLowerName[52:57]:   Store,
	_OpName[57:62]:        Loadb,
	_OpLowerName[57:62]:   Loadb,
	_OpName[62:68]:        Storeb,
	_OpLowerName[62:68]:   Storeb,
	_OpName[68:71]:        Add,
	_OpLowerName[68:71]:   Add,
	_OpName[71:74]:        Sub,
	_OpLowerName[71:74]:   Sub,
	_OpName[74:78]:        Addc,
	_OpLowerName[74:78]:   Addc,
	_OpName[78:82]:        Subc,
	_OpLowerName[78:82]:   Subc,
	_OpName[82:85]:        Xor,
	_OpLowerName[82:85]:   Xor,
	_OpName[85:88]:        And,
	_OpLowerName[85:88]:   And,
	_OpName[88:90]:        Or,
	_OpLowerName[88:90]:   Or,
	_OpName[90:93]:        Shl,
	_OpLowerName[90:93]:   Shl,
	_OpName[93:96]:        Shr,
	_OpLowerName[93:96]:   Shr,
	_OpName[96:99]:        Asr,
	_OpLowerName[96:99]:   Asr,
	_OpName[99:103]:       IfEq,
	_OpLowerName[99:103]:  IfEq,
	_OpName[103:107]:      IfNe,
	_OpLowerName[103:107]: IfNe,
	_OpName[107:111]:      IfLt,
	_OpLowerName[107:111]: IfLt,
	_OpName[111:115]:      IfGe,
	_OpLowerName[111:115]: IfGe,
	_OpName[115:120]:      IfUlt,
	_OpLowerName[115:120]: IfUlt,
	_OpName[120:125]:      IfUge,
	_OpLowerName[120:125]: IfUge,
}

var _OpNames = []string{
	_OpName[0:3],
	_OpName[3:7],
	_OpName[7:12],
	_OpName[12:16],
	_OpName[16:20],
	_OpName[20:24],
	_OpName[24:28],
	_OpName[28:33],
	_OpName[33:37],
	_OpName[37:40],
	_OpName[40:44],
	_OpName[44:48],
	_OpName[48:52],
	_OpName[52:57],
	_OpName[57:62],
	_OpName[62:68],
	_OpName[68:71],
	_OpName[71:74],
	_OpName[74:78],
	_OpName[78:82],
	_OpName[82:85],
	_OpName[85:88],
	_OpName[88:90],
	_OpName[90:93],
	_OpName[93:96],
	_OpName[96:99],
	_OpName[99:103],
	_OpName[103:107],
	_OpName[107:111],
	_OpName[111:115],
	_OpName[115:120],
	_OpName[120:125],
}

// OpString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OpString(s string) (Op, error) {
	if val, ok := _OpNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OpNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Op values", s)
}

// OpValues returns all values of the enum
func OpValues() []Op {
	return _OpValues
}

// OpStrings returns a slice of all String values of the enum
func OpStrings() []string {
	strs := make([]string, len(_OpNames))
	copy(strs, _OpNames)
	return strs
}

// IsAOp returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Op) IsAOp() bool {
	for _, v := range _OpValues {
		if i == v {
			return true
		}
	}
	return false
}

const (
	_FmtName_0      = "rrextlsri6ri8i11"
	_FmtLowerName_0 = "rrextlsri6ri8i11"
	_FmtName_1      = "i12"
	_FmtLowerName_1 = "i12"
)

var (
	_FmtIndex_0 = [...]uint8{0, 2, 5, 7, 10, 13, 16}
	_FmtIndex_1 = [...]uint8{0, 3}
)

func (i Fmt) String() string {
	switch {
	case 0 <= i && i <= 5:
		return _FmtName_0[_FmtIndex_0[i]:_FmtIndex_0[i+1]]
	case i == 13:
		return _FmtName_1
	default:
		return fmt.Sprintf("Fmt(%d)", i)
	}
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _FmtNoOp() {
	var x [1]struct{}
	_ = x[FmtRR-(0)]
	_ = x[FmtExt-(1)]
	_ = x[FmtLS-(2)]
	_ = x[FmtRI6-(3)]
	_ = x[FmtRI8-(4)]
	_ = x[FmtI11-(5)]
	_ = x[FmtI12-(13)]
}

var _FmtValues = []Fmt{FmtRR, FmtExt, FmtLS, FmtRI6, FmtRI8, FmtI11, FmtI12}

var _FmtNameToValueMap = map[string]Fmt{
	_FmtName_0[0:2]:        FmtRR,
	_FmtLowerName_0[0:2]:   FmtRR,
	_FmtName_0[2:5]:        FmtExt,
	_FmtLowerName_0[2:5]:   FmtExt,
	_FmtName_0[5:7]:        FmtLS,
	_FmtLowerName_0[5:7]:   FmtLS,
	_FmtName_0[7:10]:       FmtRI6,
	_FmtLowerName_0[7:10]:  FmtRI6,
	_FmtName_0[10:13]:      FmtRI8,
	_FmtLowerName_0[10:13]: FmtRI8,
	_FmtName_0[13:16]:      FmtI11,
	_FmtLowerName_0[13:16]: FmtI11,
	_FmtName_1[0:3]:        FmtI12,
	_FmtLowerName_1[0:3]:   FmtI12,
}

var _FmtNames = []string{
	_FmtName_0[0:2],
	_FmtName_0[2:5],
	_FmtName_0[5:7],
	_FmtName_0[7:10],
	_FmtName_0[10:13],
	_FmtName_0[13:16],
	_FmtName_1[0:3],
}

// FmtString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func FmtString(s string) (Fmt, error) {
	if val, ok := _FmtNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _FmtNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Fmt values", s)
}

// FmtValues returns all values of the enum
func FmtValues() []Fmt {
	return _FmtValues
}

// FmtStrings returns a slice of all String values of the enum
func FmtStrings() []string {
	strs := make([]string, len(_FmtNames))
	copy(strs, _FmtNames)
	return strs
}

// IsAFmt returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Fmt) IsAFmt() bool {
	for _, v := range _FmtValues {
		if i == v {
			return true
		}
	}
	return false
}
