// Code generated by "enumer -type=Opcode"; DO NOT EDIT.

package a32

import (
	"fmt"
	"strings"
)

const _OpcodeName = "NOPBRKHLTERRADDSUBADDCSUBBANDORXORSHLASRLSRLDSTLD8ST8LD16ST16BR_EQBR_NEQBR_U_LBR_U_LEBR_U_GEBR_U_GBR_S_LBR_S_LEBR_S_GEBR_S_GBRACALLRETJMPCMPNEGNEGBNOTMOVLDISWPNumOps"

var _OpcodeIndex = [...]uint8{0, 3, 6, 9, 12, 15, 18, 22, 26, 29, 31, 34, 37, 40, 43, 45, 47, 50, 53, 57, 61, 66, 72, 78, 85, 92, 98, 104, 111, 118, 124, 127, 131, 134, 137, 140, 143, 147, 150, 153, 156, 159, 165}

const _OpcodeLowerName = "nopbrkhlterraddsubaddcsubbandorxorshlasrlsrldstld8st8ld16st16br_eqbr_neqbr_u_lbr_u_lebr_u_gebr_u_gbr_s_lbr_s_lebr_s_gebr_s_gbracallretjmpcmpnegnegbnotmovldiswpnumops"

func (i Opcode) String() string {
	if i < 0 || i >= Opcode(len(_OpcodeIndex)-1) {
		return fmt.Sprintf("Opcode(%d)", i)
	}
	return _OpcodeName[_OpcodeIndex[i]:_OpcodeIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _OpcodeNoOp() {
	var x [1]struct{}
	_ = x[NOP-(0)]
	_ = x[BRK-(1)]
	_ = x[HLT-(2)]
	_ = x[ERR-(3)]
	_ = x[ADD-(4)]
	_ = x[SUB-(5)]
	_ = x[ADDC-(6)]
	_ = x[SUBB-(7)]
	_ = x[AND-(8)]
	_ = x[OR-(9)]
	_ = x[XOR-(10)]
	_ = x[SHL-(11)]
	_ = x[ASR-(12)]
	_ = x[LSR-(13)]
	_ = x[LD-(14)]
	_ = x[ST-(15)]
	_ = x[LD8-(16)]
	_ = x[ST8-(17)]
	_ = x[LD16-(18)]
	_ = x[ST16-(19)]
	_ = x[BR_EQ-(20)]
	_ = x[BR_NEQ-(21)]
	_ = x[BR_U_L-(22)]
	_ = x[BR_U_LE-(23)]
	_ = x[BR_U_GE-(24)]
	_ = x[BR_U_G-(25)]
	_ = x[BR_S_L-(26)]
	_ = x[BR_S_LE-(27)]
	_ = x[BR_S_GE-(28)]
	_ = x[BR_S_G-(29)]
	_ = x[BRA-(30)]
	_ = x[CALL-(31)]
	_ = x[RET-(32)]
	_ = x[JMP-(33)]
	_ = x[CMP-(34)]
	_ = x[NEG-(35)]
	_ = x[NEGB-(36)]
	_ = x[NOT-(37)]
	_ = x[MOV-(38)]
	_ = x[LDI-(39)]
	_ = x[SWP-(40)]
	_ = x[NumOps-(41)]
}

var _OpcodeValues = []Opcode{NOP, BRK, HLT, ERR, ADD, SUB, ADDC, SUBB, AND, OR, XOR, SHL, ASR, LSR, LD, ST, LD8, ST8, LD16, ST16, BR_EQ, BR_NEQ, BR_U_L, BR_U_LE, BR_U_GE, BR_U_G, BR_S_L, BR_S_LE, BR_S_GE, BR_S_G, BRA, CALL, RET, JMP, CMP, NEG, NEGB, NOT, MOV, LDI, SWP, NumOps}

var _OpcodeNameToValueMap = map[string]Opcode{
	_OpcodeName[0:3]:          NOP,
	_OpcodeLowerName[0:3]:     NOP,
	_OpcodeName[3:6]:          BRK,
	_OpcodeLowerName[3:6]:     BRK,
	_OpcodeName[6:9]:          HLT,
	_OpcodeLowerName[6:9]:     HLT,
	_OpcodeName[9:12]:         ERR,
	_OpcodeLowerName[9:12]:    ERR,
	_OpcodeName[12:15]:        ADD,
	_OpcodeLowerName[12:15]:   ADD,
	_OpcodeName[15:18]:        SUB,
	_OpcodeLowerName[15:18]:   SUB,
	_OpcodeName[18:22]:        ADDC,
	_OpcodeLowerName[18:22]:   ADDC,
	_OpcodeName[22:26]:        SUBB,
	_OpcodeLowerName[22:26]:   SUBB,
	_OpcodeName[26:29]:        AND,
	_OpcodeLowerName[26:29]:   AND,
	_OpcodeName[29:31]:        OR,
	_OpcodeLowerName[29:31]:   OR,
	_OpcodeName[31:34]:        XOR,
	_OpcodeLowerName[31:34]:   XOR,
	_OpcodeName[34:37]:        SHL,
	_OpcodeLowerName[34:37]:   SHL,
	_OpcodeName[37:40]:        ASR,
	_OpcodeLowerName[37:40]:   ASR,
	_OpcodeName[40:43]:        LSR,
	_OpcodeLowerName[40:43]:   LSR,
	_OpcodeName[43:45]:        LD,
	_OpcodeLowerName[43:45]:   LD,
	_OpcodeName[45:47]:        ST,
	_OpcodeLowerName[45:47]:   ST,
	_OpcodeName[47:50]:        LD8,
	_OpcodeLowerName[47:50]:   LD8,
	_OpcodeName[50:53]:        ST8,
	_OpcodeLowerName[50:53]:   ST8,
	_OpcodeName[53:57]:        LD16,
	_OpcodeLowerName[53:57]:   LD16,
	_OpcodeName[57:61]:        ST16,
	_OpcodeLowerName[57:61]:   ST16,
	_OpcodeName[61:66]:        BR_EQ,
	_OpcodeLowerName[61:66]:   BR_EQ,
	_OpcodeName[66:72]:        BR_NEQ,
	_OpcodeLowerName[66:72]:   BR_NEQ,
	_OpcodeName[72:78]:        BR_U_L,
	_OpcodeLowerName[72:78]:   BR_U_L,
	_OpcodeName[78:85]:        BR_U_LE,
	_OpcodeLowerName[78:85]:   BR_U_LE,
	_OpcodeName[85:92]:        BR_U_GE,
	_OpcodeLowerName[85:92]:   BR_U_GE,
	_OpcodeName[92:98]:        BR_U_G,
	_OpcodeLowerName[92:98]:   BR_U_G,
	_OpcodeName[98:104]:       BR_S_L,
	_OpcodeLowerName[98:104]:  BR_S_L,
	_OpcodeName[104:111]:      BR_S_LE,
	_OpcodeLowerName[104:111]: BR_S_LE,
	_OpcodeName[111:118]:      BR_S_GE,
	_OpcodeLowerName[111:118]: BR_S_GE,
	_OpcodeName[118:124]:      BR_S_G,
	_OpcodeLowerName[118:124]: BR_S_G,
	_OpcodeName[124:127]:      BRA,
	_OpcodeLowerName[124:127]: BRA,
	_OpcodeName[127:131]:      CALL,
	_OpcodeLowerName[127:131]: CALL,
	_OpcodeName[131:134]:      RET,
	_OpcodeLowerName[131:134]: RET,
	_OpcodeName[134:137]:      JMP,
	_OpcodeLowerName[134:137]: JMP,
	_OpcodeName[137:140]:      CMP,
	_OpcodeLowerName[137:140]: CMP,
	_OpcodeName[140:143]:      NEG,
	_OpcodeLowerName[140:143]: NEG,
	_OpcodeName[143:147]:      NEGB,
	_OpcodeLowerName[143:147]: NEGB,
	_OpcodeName[147:150]:      NOT,
	_OpcodeLowerName[147:150]: NOT,
	_OpcodeName[150:153]:      MOV,
	_OpcodeLowerName[150:153]: MOV,
	_OpcodeName[153:156]:      LDI,
	_OpcodeLowerName[153:156]: LDI,
	_OpcodeName[156:159]:      SWP,
	_OpcodeLowerName[156:159]: SWP,
	_OpcodeName[159:165]:      NumOps,
	_OpcodeLowerName[159:165]: NumOps,
}

var _OpcodeNames = []string{
	_OpcodeName[0:3],
	_OpcodeName[3:6],
	_OpcodeName[6:9],
	_OpcodeName[9:12],
	_OpcodeName[12:15],
	_OpcodeName[15:18],
	_OpcodeName[18:22],
	_OpcodeName[22:26],
	_OpcodeName[26:29],
	_OpcodeName[29:31],
	_OpcodeName[31:34],
	_OpcodeName[34:37],
	_OpcodeName[37:40],
	_OpcodeName[40:43],
	_OpcodeName[43:45],
	_OpcodeName[45:47],
	_OpcodeName[47:50],
	_OpcodeName[50:53],
	_OpcodeName[53:57],
	_OpcodeName[57:61],
	_OpcodeName[61:66],
	_OpcodeName[66:72],
	_OpcodeName[72:78],
	_OpcodeName[78:85],
	_OpcodeName[85:92],
	_OpcodeName[92:98],
	_OpcodeName[98:104],
	_OpcodeName[104:111],
	_OpcodeName[111:118],
	_OpcodeName[118:124],
	_OpcodeName[124:127],
	_OpcodeName[127:131],
	_OpcodeName[131:134],
	_OpcodeName[134:137],
	_OpcodeName[137:140],
	_OpcodeName[140:143],
	_OpcodeName[143:147],
	_OpcodeName[147:150],
	_OpcodeName[150:153],
	_OpcodeName[153:156],
	_OpcodeName[156:159],
	_OpcodeName[159:165],
}

// OpcodeString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func OpcodeString(s string) (Opcode, error) {
	if val, ok := _OpcodeNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _OpcodeNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Opcode values", s)
}

// OpcodeValues returns all values of the enum
func OpcodeValues() []Opcode {
	return _OpcodeValues
}

// OpcodeStrings returns a slice of all String values of the enum
func OpcodeStrings() []string {
	strs := make([]string, len(_OpcodeNames))
	copy(strs, _OpcodeNames)
	return strs
}

// IsAOpcode returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Opcode) IsAOpcode() bool {
	for _, v := range _OpcodeValues {
		if i == v {
			return true
		}
	}
	return false
}
