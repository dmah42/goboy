package z80

type opcode int

const (
	// 0x00
	NOP opcode = iota
	LDBCnn
	LDBCmA
	INCBC
	INCr_b
	DECr_b
	LDrn_b
	RLCA
	LDmmSP
	ADDHLBC
	LDABCm
	DECBC
	INCr_c
	DECr_c
	LDrn_c
	RRCA

	// 0x10
)


