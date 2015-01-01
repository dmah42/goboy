package goboy

var Opcodes = [0x100]string{
	// 0x00
	"NOP", "LDBCnn", "LDBCmA", "INCBC",
	"INCr_b", "DECr_b", "LDrn_b", "RLCA",
	"LDmmSP", "ADDHLBC", "LDABCm", "DECBC",
	"INCr_c", "DECr_c", "LDrn_c", "RRCA",
	// 0x10
	"DJNZn", "LDDEnn", "LDDEmA", "INCDE",
	"INCr_d", "DECr_d", "LDrn_d", "RLA",
	"JRn", "ADDHLDE", "LDADEm", "DECDE",
	"INCr_e", "DECr_e", "LDrn_e", "RRA",

	// 0x20
	"JRNZn", "LDHLnn", "LDHLIA", "INCHL",
	"INCr_h", "DECr_h", "LDrn_h", "DAA",
	"JRZn", "ADDHLHL", "LDAHLI", "DECHL",
	"INCr_l", "DECr_l", "LDrn_l", "CPL",

	// 0x30
	"JRNCn", "LDSPnn", "LDHLDA", "INCSP",
	"INCHLm", "DECHLm", "LDHLmn", "SCF",
	"JRCn", "ADDHLSP", "LDAHLD", "DECSP",
	"INCr_a", "DECr_a", "LDrn_a", "CCF",

	// 0x40
	"LDrr_bb", "LDrr_bc", "LDrr_bd", "LDrr_be",
	"LDrr_bh", "LDrr_bl", "LDrHLm_b", "LDrr_ba",
	"LDrr_cb", "LDrr_cc", "LDrr_cd", "LDrr_ce",
	"LDrr_ch", "LDrr_cl", "LDrHLm_c", "LDrr_ca",

	// 0x50
	"LDrr_db", "LDrr_dc", "LDrr_dd", "LDrr_de",
	"LDrr_dh", "LDrr_dl", "LDrHLm_d", "LDrr_da",
	"LDrr_eb", "LDrr_ec", "LDrr_ed", "LDrr_ee",
	"LDrr_eh", "LDrr_el", "LDrHLm_e", "LDrr_ea",

	// 0x60
	"LDrr_hb", "LDrr_hc", "LDrr_hd", "LDrr_he",
	"LDrr_hh", "LDrr_hl", "LDrHLm_h", "LDrr_ha",
	"LDrr_lb", "LDrr_lc", "LDrr_ld", "LDrr_le",
	"LDrr_lh", "LDrr_ll", "LDrHLm_l", "LDrr_la",

	// 0x70
	"LDHLmr_b", "LDHLmr_c", "LDHLmr_d", "LDHLmr_e",
	"LDHLmr_h", "LDHLmr_l", "HALT", "LDHLmr_a",
	"LDrr_ab", "LDrr_ac", "LDrr_ad", "LDrr_ae",
	"LDrr_ah", "LDrr_al", "LDrHLm_a", "LDrr_aa",

	// 0x80
	"ADDr_b", "ADDr_c", "ADDr_d", "ADDr_e",
	"ADDr_h", "ADDr_l", "ADDHL", "ADDr_a",
	"ADCr_b", "ADCr_c", "ADCr_d", "ADCr_e",
	"ADCr_h", "ADCr_l", "ADCHL", "ADCr_a",

	// 0x90
	"SUBr_b", "SUBr_c", "SUBr_d", "SUBr_e",
	"SUBr_h", "SUBr_l", "SUBHL", "SUBr_a",
	"SBCr_b", "SBCr_c", "SBCr_d", "SBCr_e",
	"SBCr_h", "SBCr_l", "SBCHL", "SBCr_a",

	// 0xA0
	"ANDr_b", "ANDr_c", "ANDr_d", "ANDr_e",
	"ANDr_h", "ANDr_l", "ANDHL", "ANDr_a",
	"XORr_b", "XORr_c", "XORr_d", "XORr_e",
	"XORr_h", "XORr_l", "XORHL", "XORr_a",

	// 0xB0
	"ORr_b", "ORr_c", "ORr_d", "ORr_e",
	"ORr_h", "ORr_l", "ORHL", "ORr_a",
	"CPr_b", "CPr_c", "CPr_d", "CPr_e",
	"CPr_h", "CPr_l", "CPHL", "CPr_a",

	// 0xC0
	"RETNZ", "POPBC", "JPNZnn", "JPnn",
	"CALLNZnn", "PUSHBC", "ADDn", "RST00",
	"RETZ", "RET", "JPZnn", "MAPcb",
	"CALLZnn", "CALLnn", "ADCn", "RST08",

	// 0xD0
	"RETNC", "POPDE", "JPNCnn", "XX",
	"CALLNCnn", "PUSHDE", "SUBn", "RST10",
	"RETC", "RETI", "JPCnn", "XX",
	"CALLCnn", "XX", "SBCn", "RST18",

	// 0xE0
	"LDIOnA", "POPHL", "LDIOCA", "XX",
	"XX", "PUSHHL", "ANDn", "RST20",
	"ADDSPn", "JPHL", "LDmmA", "XX",
	"XX", "XX", "XORn", "RST28",

	// 0xF0
	"LDAIOn", "POPAF", "LDAIOC", "DI",
	"XX", "PUSHAF", "ORn", "RST30",
	"LDHLSPn", "XX", "LDAmm", "EI",
	"XX", "XX", "CPn", "RST38",
}
