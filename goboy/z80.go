package goboy

import (
	"log"
)

type saveRegisters struct {
	a, b, c, d, e, h, l, f uint8
}

type registers struct {
	a, b, c, d, e, h, l, f uint8
	Pc, sp, i, r           uint16
	M                      int
	Ime                    bool
}

type instruction struct {
	name string
	f interface{}
}

type z80 struct {
	R	registers
	M	int
	Halt	bool
}

func makeZ80() z80 {
	var z z80
	z.R.Ime = true
	return z
}

var (
	savedRegisters saveRegisters

	instructions = [0x100]instruction {
		// 0x00
		{"NOP", func() { Z80.R.M = 1 }},
		{"LDBCnn", func() {
			Z80.R.c = MMU.ReadByte(Z80.R.Pc)
			Z80.R.b = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		}},
		{"LDBCmA", func() {
			MMU.WriteByte((uint16(Z80.R.b)<<8)+uint16(Z80.R.c), Z80.R.a)
			Z80.R.M = 2
		}},
		{"INCBC", func() {
			Z80.R.c = (Z80.R.c + 1) & 0xFF
			if Z80.R.c == 0 {
				Z80.R.b = (Z80.R.b + 1) & 0xFF
			}
			Z80.R.M = 1
		}},

		{"INCr_b", func() {
			Z80.R.b = (Z80.R.b + 1) & 0xFF
			Z80.R.f = 0
			if Z80.R.b == 0 {
				Z80.R.f = 0x80
			}
			Z80.R.M = 1
		}},
		{"DEC_r_b", nil},
		{"LDrn_b", func() {
			Z80.R.b = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"RLCA", nil},

		{"LDmmSP", nil},
		{"ADDHLBC", func() {
			hl := uint(Z80.R.h << 8) + uint(Z80.R.l)
			hl += uint(Z80.R.b << 8) + uint(Z80.R.c)
			if hl > 0xFFFF {
				Z80.R.f |= 0x10
			} else {
				Z80.R.f &= 0xEF
			}
			Z80.R.h = uint8((hl >> 8) & 0xFF)
			Z80.R.l = uint8(hl & 0xFF)
			Z80.R.M = 3
		}},
		{"LDABCm", nil},
		{"DECBC", nil},

		{"INCr_c", nil},
		{"DECr_c", nil},
		{"LDrn_c", func() {
			Z80.R.c = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"RRCA", nil},

		// 0x10
		{"DJNZn", nil},
		{"LDDEnn", func() {
			Z80.R.d = MMU.ReadByte(Z80.R.Pc)
			Z80.R.e = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		}},
		{"LDDEmA", func() {
			MMU.WriteByte((uint16(Z80.R.d)<<8)+uint16(Z80.R.e), Z80.R.a)
			Z80.R.M = 2
		}},
		{"INCDE", func() {
			Z80.R.e = (Z80.R.e + 1) & 0xFF
			if Z80.R.e == 0 {
				Z80.R.d = (Z80.R.d + 1) & 0xFF
			}
			Z80.R.M = 1
		}},

		{"INCr_d", nil},
		{"DECr_d", nil},
		{"LDrn_d", func() {
			Z80.R.d = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"RLA", nil},

		{"JRn", nil},
		{"ADDHLDE", func() {
			hl := uint(Z80.R.h << 8) + uint(Z80.R.l)
			hl += uint(Z80.R.d << 8) + uint(Z80.R.e)
			if hl > 0xFFFF {
				Z80.R.f |= 0x10
			} else {
				Z80.R.f &= 0xEF
			}
			Z80.R.h = uint8((hl >> 8) & 0xFF)
			Z80.R.l = uint8(hl & 0xFF)
			Z80.R.M = 3
		}},
		{"LDADEm", nil},
		{"DECDE", nil},

		{"INCr_e", nil},
		{"DECr_e", nil},
		{"LDrn_e", func() {
			Z80.R.e = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"RRA", nil},

		// 0x20
		{"JRNZn", nil},
		{"LDHLnn", func() {
			Z80.R.l = MMU.ReadByte(Z80.R.Pc)
			Z80.R.h = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		}},
		{"LDHLIA", nil},
		{"INCHL", func() {
			Z80.R.l = (Z80.R.l + 1) & 0xFF
			if Z80.R.l == 0 {
				Z80.R.h = (Z80.R.h + 1) & 0xFF
			}
			Z80.R.M = 1
		}},

		{"INCr_h", nil},
		{"DECr_h", nil},
		{"LDrn_h", func() {
			Z80.R.h = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"DAA", nil},

		{"JRZn", nil},
		{"ADDHLHL", func() {
			hl := uint(Z80.R.h << 8) + uint(Z80.R.l)
			hl += uint(Z80.R.h << 8) + uint(Z80.R.l)
			if hl > 0xFFFF {
				Z80.R.f |= 0x10
			} else {
				Z80.R.f &= 0xEF
			}
			Z80.R.h = uint8((hl >> 8) & 0xFF)
			Z80.R.l = uint8(hl & 0xFF)
			Z80.R.M = 3
		}},
		{"LDAHLI", nil},
		{"DECHL", nil},

		{"INCr_l", nil},
		{"DECr_l", nil},
		{"LDrn_l", func() {
			Z80.R.l = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"CPL", nil},

		// 0x30
		{"JRNCn", nil},
		{"LDSPnn", func() {
			Z80.R.sp = MMU.ReadWord(Z80.R.Pc)
			Z80.R.Pc += 2
			Z80.R.M = 3
		}},
		{"LDHLDA", nil},
		{"INCSP", func() {
			Z80.R.sp = (Z80.R.sp + 1) & 0xFFFF
			Z80.R.M = 1
		}},

		{"INCHLm", nil},
		{"DECHLm", nil},
		{"LDHLmn", nil},
		{"SCF", nil},

		{"JRCn", nil},
		{"ADDHLSP", func() {
			hl := uint(Z80.R.h << 8) + uint(Z80.R.l)
			hl += uint(Z80.R.sp)
			if hl > 0xFFFF {
				Z80.R.f |= 0x10
			} else {
				Z80.R.f &= 0xEF
			}
			Z80.R.h = uint8((hl >> 8) & 0xFF)
			Z80.R.l = uint8(hl & 0xFF)
			Z80.R.M = 3
		}},
		{"LDAHLD", nil},
		{"DECSP", nil},

		{"INCr_a", nil},
		{"DECr_a", nil},
		{"LDrn_a", func() {
			Z80.R.a = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"CCF", nil},

		// 0x40
		{"LDrr_bb", func() {
			Z80.R.b = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_bc", func() {
			Z80.R.b = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_bd", func() {
			Z80.R.b = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_be", func() {
			Z80.R.b = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_bh", func() {
			Z80.R.b = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_bl", func() {
			Z80.R.b = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_b", func() {
			Z80.R.b = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_ba", func() {
			Z80.R.b = Z80.R.a
			Z80.R.M = 1
		}},

		{"LDrr_cb", func() {
			Z80.R.c = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_cc", func() {
			Z80.R.c = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_cd", func() {
			Z80.R.c = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_ce", func() {
			Z80.R.c = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_ch", func() {
			Z80.R.c = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_cl", func() {
			Z80.R.c = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_c", func() {
			Z80.R.c = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_ca", func() {
			Z80.R.c = Z80.R.a
			Z80.R.M = 1
		}},

		// 0x50
		{"LDrr_db", func() {
			Z80.R.d = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_dc", func() {
			Z80.R.d = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_dd", func() {
			Z80.R.d = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_de", func() {
			Z80.R.d = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_dh", func() {
			Z80.R.d = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_dl", func() {
			Z80.R.d = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_d", func() {
			Z80.R.d = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_da", func() {
			Z80.R.d = Z80.R.a
			Z80.R.M = 1
		}},

		{"LDrr_eb", func() {
			Z80.R.e = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_ec", func() {
			Z80.R.e = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_ed", func() {
			Z80.R.e = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_ee", func() {
			Z80.R.e = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_eh", func() {
			Z80.R.e = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_el", func() {
			Z80.R.e = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_e", func() {
			Z80.R.e = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_ea", func() {
			Z80.R.e = Z80.R.a
			Z80.R.M = 1
		}},

		// 0x60
		{"LDrr_hb", func() {
			Z80.R.h = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_hc", func() {
			Z80.R.h = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_hd", func() {
			Z80.R.h = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_he", func() {
			Z80.R.h = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_hh", func() {
			Z80.R.h = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_hl", func() {
			Z80.R.h = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_h", func() {
			Z80.R.h = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_ha", func() {
			Z80.R.h = Z80.R.a
			Z80.R.M = 1
		}},

		{"LDrr_lb", func() {
			Z80.R.l = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_lc", func() {
			Z80.R.l = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_ld", func() {
			Z80.R.l = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_le", func() {
			Z80.R.l = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_lh", func() {
			Z80.R.l = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_ll", func() {
			Z80.R.l = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_l", func() {
			Z80.R.l = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_la", func() {
			Z80.R.l = Z80.R.a
			Z80.R.M = 1
		}},

		// 0x70
		{"LDHLmr_b", nil},
		{"LDHLmr_c", nil},
		{"LDHLmr_d", nil},
		{"LDHLmr_e", nil},

		{"LDHLmr_h", nil},
		{"LDHLmr_l", nil},
		{"HALT", func() {
			Z80.Halt = true
			Z80.R.M = 1
		}},
		{"LDHLmr_a", nil},

		{"LDrr_ab", func() {
			Z80.R.a = Z80.R.b
			Z80.R.M = 1
		}},
		{"LDrr_ac", func() {
			Z80.R.a = Z80.R.c
			Z80.R.M = 1
		}},
		{"LDrr_ad", func() {
			Z80.R.a = Z80.R.d
			Z80.R.M = 1
		}},
		{"LDrr_ae", func() {
			Z80.R.a = Z80.R.e
			Z80.R.M = 1
		}},

		{"LDrr_ah", func() {
			Z80.R.a = Z80.R.h
			Z80.R.M = 1
		}},
		{"LDrr_al", func() {
			Z80.R.a = Z80.R.l
			Z80.R.M = 1
		}},
		{"LDrHLm_a", func() {
			Z80.R.a = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		}},
		{"LDrr_aa", func() {
			Z80.R.a = Z80.R.a
			Z80.R.M = 1
		}},

		// 0x80
		{"ADDr_b", nil},
		{"ADDr_c", nil},
		{"ADDr_d", nil},
		{"ADDr_e", nil},

		{"ADDr_h", nil},
		{"ADDr_l", nil},
		{"ADDHL", nil},
		{"ADDr_a", func() {
			a := Z80.R.a
			newA := uint16(Z80.R.a) + uint16(Z80.R.a)
			Z80.R.a = uint8(newA & 0xFF)
			Z80.R.f = 0x0
			if newA > 0xFF {
				Z80.R.f = 0x10
			}
			if Z80.R.a == 0 {
				Z80.R.f |= 0x80
			}
			if (Z80.R.a ^ Z80.R.a ^ a) & 0x10 != 0 {
				Z80.R.f |= 0x20
			}
			Z80.R.M = 1
		}},

		{"ADCr_b", nil},
		{"ADCr_c", nil},
		{"ADCr_d", nil},
		{"ADCr_e", nil},

		{"ADCr_h", nil},
		{"ADCr_l", nil},
		{"ADCHL", nil},
		{"ADCr_a", nil},

		// 0x90
		{"SUBr_b", nil},
		{"SUBr_c", nil},
		{"SUBr_d", nil},
		{"SUBr_e", nil},

		{"SUBr_h", nil},
		{"SUBr_l", nil},
		{"SUBHL", nil},
		{"SUBr_a", func() {
			a := Z80.R.a
			newA := int8(Z80.R.a) - int8(Z80.R.a)
			Z80.R.a = uint8(newA)
			Z80.R.f = 0x40
			if newA < 0 {
				Z80.R.f = 0x50
			}
			if Z80.R.a == 0 {
				Z80.R.f |= 0x80
			}
			if (Z80.R.a ^ Z80.R.a ^ a) & 0x10 != 0 {
				Z80.R.f |= 0x20
			}
			Z80.R.M = 1
		}},

		{"SBCr_b", nil},
		{"SBCr_c", nil},
		{"SBCr_d", nil},
		{"SBCr_e", nil},

		{"SBCr_h", nil},
		{"SBCr_l", nil},
		{"SBCHL", nil},
		{"SBCr_a", nil},

		// 0xA0
		{"ANDr_b", nil},
		{"ANDr_c", nil},
		{"ANDr_d", nil},
		{"ANDr_e", nil},

		{"ANDr_h", nil},
		{"ANDr_l", nil},
		{"ANDHL", nil},
		{"ANDr_a", nil},

		{"XORr_b", nil},
		{"XORr_c", nil},
		{"XORr_d", nil},
		{"XORr_e", nil},

		{"XORr_h", nil},
		{"XORr_l", nil},
		{"XORHL", nil},
		{"XORr_a", nil},

		// 0xB0
		{"ORr_b", nil},
		{"ORr_c", nil},
		{"ORr_d", nil},
		{"ORr_e", nil},

		{"ORr_h", nil},
		{"ORr_l", nil},
		{"ORHL", nil},
		{"ORr_a", nil},

		{"CPr_b", nil},
		{"CPr_c", nil},
		{"CPr_d", nil},
		{"CPr_e", nil},

		{"CPr_h", nil},
		{"CPr_l", nil},
		{"CPHL", nil},
		{"CPr_a", nil},

		// 0xC0
		{"RETNZ", func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x80) == 0 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		}},
		{"POPBC", func() {
			Z80.R.c = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.b = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		}},
		{"JPNZnn", nil},
		{"JPnn", func() {
			Z80.R.Pc = MMU.ReadWord(Z80.R.Pc)
			Z80.R.M = 3
		}},

		{"CALLNZnn", nil},
		{"PUSHBC", func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.b)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.c)
			Z80.R.M = 3
		}},
		{"ADDn", nil},
		{"RST00", func() { Z80.reset(0x00) }},

		{"RETZ", func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x80) == 0x80 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		}},
		{"RET", func() {
			Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
			Z80.R.sp += 2
			Z80.R.M = 3
		}},
		{"JPZnn", nil},
		{"MAPcb", nil},

		{"CALLZnn", nil},
		{"CALLnn", func() {
			Z80.R.sp -= 2
			MMU.WriteWord(Z80.R.sp, Z80.R.Pc + 2)
			Z80.R.Pc = MMU.ReadWord(Z80.R.Pc)
			Z80.R.M = 5
		}},
		{"ADCn", nil},
		{"RST08", func() { Z80.reset(0x08) }},

		// 0xD0
		{"RETNC", func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x10) == 0 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		}},
		{"POPDE", func() {
			Z80.R.e = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.d = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		}},
		{"JPNCnn", nil},
		{"XX", nil},

		{"CALLNCnn", nil},
		{"PUSHDE", func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.d)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.e)
			Z80.R.M = 3
		}},
		{"SUBn", nil},
		{"RST10", func() { Z80.reset(0x10) }},

		{"RETC", func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x10) == 0x10 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		}},
		{"RETI", func() {
			Z80.R.Ime = true
			Z80.loadRegisters()
			Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
			Z80.R.sp += 2
			Z80.R.M = 3
		}},
		{"JPCnn", nil},
		{"XX", nil},

		{"CALLCnn", nil},
		{"XX", nil},
		{"SBCn", nil},
		{"RST18", func() { Z80.reset(0x18) }},

		// 0xE0
		{"LDIOnA", func() {
			MMU.WriteByte(uint16(0xFF00) + uint16(MMU.ReadByte(Z80.R.Pc)), Z80.R.a)
			Z80.R.Pc += 1
			Z80.R.M = 3
		}},
		{"POPHL", func() {
			Z80.R.l = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.h = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		}},
		{"LDIOCA", func() {
			MMU.WriteByte(uint16(0xFF00) + uint16(Z80.R.c), Z80.R.a)
			Z80.R.M = 2
		}},
		{"XX", nil},

		{"XX", nil},
		{"PUSHHL", func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.h)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.l)
			Z80.R.M = 3
		}},
		{"ANDn", nil},
		{"RST20", func() { Z80.reset(0x20) }},

		{"ADDSPn", func() {
			i := int16(MMU.ReadByte(Z80.R.Pc))
			if i > 127 {
				i = -((^i+1)&0xFF)
			}
			Z80.R.Pc += 1
			sp := int16(Z80.R.sp) + i
			Z80.R.sp = uint16(sp)
			Z80.R.M = 4
		}},
		{"JPHL", nil},
		{"LDmmA", func() {
			MMU.WriteByte(MMU.ReadWord(Z80.R.Pc), Z80.R.a)
			Z80.R.Pc += 2
			Z80.R.M = 4
		}},
		{"XX", nil},

		{"XX", nil},
		{"XX", nil},
		{"XORn", nil},
		{"RST28", func() { Z80.reset(0x28) }},

		// 0xF0
		{"LDAIOn", func() {
			Z80.R.a = MMU.ReadByte(uint16(0xFF00) + uint16(MMU.ReadByte(Z80.R.Pc)))
			Z80.R.Pc += 1
			Z80.R.M = 3
		}},
		{"POPAF", func() {
			Z80.R.f = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.a = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		}},
		{"LDAIOC", func() {
			Z80.R.a = MMU.ReadByte(uint16(0xFF00) + uint16(Z80.R.c))
			Z80.R.M = 2
		}},
		{"DI", func() {
			Z80.R.Ime = false
			Z80.R.M = 1
		}},

		{"XX", nil},
		{"PUSHAF", func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.a)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.f)
			Z80.R.M = 3
		}},
		{"ORn", nil},
		{"RST30", func() { Z80.reset(0x30) }},

		{"LDHLSPn", nil},
		{"XX", nil},
		{"LDAmm", nil},
		{"EI", func() {
			Z80.R.Ime = true
			Z80.R.M = 1
		}},

		{"XX", nil},
		{"XX", nil},
		{"CPn", func() {
			i := int(Z80.R.a)
			m := int(MMU.ReadByte(Z80.R.Pc))
			i -= m
			Z80.R.Pc += 1
			Z80.R.f = 0x40
			if i < 0 {
				Z80.R.f = 0x50
			}
			i &= 0xFF
			if i == 0 {
				Z80.R.f |= 0x80
			}
			if (Z80.R.a ^ uint8(i) ^ uint8(m)) & 0x10 != 0 {
				Z80.R.f |= 0x20
			}
			Z80.R.M = 2
		}},
		{"RST38", func() { Z80.reset(0x38) }},
	}
)

func (z z80) Call(op uint8) {
	instr := instructions[op]
	if instr.f == nil {
		log.Panicf("z80: nil instruction for op [0x%x] %q\n", op, instr.name)
	}
	log.Printf("z80: op [0x%x] %q\n", op, instr.name)
	instr.f.(func())()
}

func (z *z80) storeRegisters() {
	savedRegisters.a = z.R.a
	savedRegisters.b = z.R.b
	savedRegisters.c = z.R.c
	savedRegisters.d = z.R.d
	savedRegisters.e = z.R.e
	savedRegisters.f = z.R.f
	savedRegisters.h = z.R.h
	savedRegisters.l = z.R.l
}

func (z *z80) loadRegisters() {
	z.R.a = savedRegisters.a
	z.R.b = savedRegisters.b
	z.R.c = savedRegisters.c
	z.R.d = savedRegisters.d
	z.R.e = savedRegisters.e
	z.R.f = savedRegisters.f
	z.R.h = savedRegisters.h
	z.R.l = savedRegisters.l
}

func (z *z80) reset(addr uint16) {
	z.storeRegisters()
	z.R.sp -= 2
	MMU.WriteWord(z.R.sp, z.R.Pc)
	z.R.Pc = addr
	Z80.R.M = 3
}
