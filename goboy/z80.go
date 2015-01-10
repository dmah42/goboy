package goboy

import (
	"log"
)

type saveRegisters struct {
	a, b, c, d, e, h, l, f uint8
}

type registers struct {
	a, b, c, d, e, h, l, f uint8
	Pc, sp, R	       uint16
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
		{"LDABCm", func() {
			Z80.R.a = MMU.ReadByte((uint16(Z80.R.b) << 8) + uint16(Z80.R.c))
			Z80.R.M = 2
		}},
		{"DECBC", func() {
			Z80.R.c -= 1
			if Z80.R.c == 0xff {
				Z80.R.b -= 1
			}
			Z80.R.M = 1
		}},

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

		{"INCr_d", func() {
			d := int(Z80.R.d) + 1
			Z80.R.d = uint8(d & 0xFF)
			Z80.R.f = 0
			if Z80.R.d == 0 {
				Z80.R.f = 0x80
			}
			Z80.R.M = 1
		}},
		{"DECr_d", func() {
			d := int(Z80.R.d) - 1
			Z80.R.d = uint8(d & 0xFF)
			Z80.R.f = 0
			if Z80.R.d == 0 {
				Z80.R.f = 0x80
			}
			Z80.R.M = 1
		}},
		{"LDrn_d", func() {
			Z80.R.d = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"RLA", nil},

		{"JRn", func() {
			i := int(MMU.ReadByte(Z80.R.Pc))
			if i > 127 {
				i = -((^i+1)&0xFF)
			}
			Z80.R.Pc += 1
			Z80.R.M = 2
			pc := int(Z80.R.Pc) + i
			Z80.R.Pc = uint16(pc)
			Z80.R.M += 1
		}},
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
		{"DECDE", func() {
			Z80.R.e -= 1
			if Z80.R.e == 0xff {
				Z80.R.d -= 1
			}
			Z80.R.M = 1
		}},

		{"INCr_e", func() {
			e := int(Z80.R.e) + 1
			Z80.R.e = uint8(e & 0xFF)
			Z80.R.f = 0
			if Z80.R.e == 0 {
				Z80.R.f = 0x80
			}
			Z80.R.M = 1
		}},
		{"DECr_e", func() {
			e := int(Z80.R.e) - 1
			Z80.R.e = uint8(e & 0xFF)
			Z80.R.f = 0
			if Z80.R.e == 0 {
				Z80.R.f = 0x80
			}
			Z80.R.M = 1
		}},
		{"LDrn_e", func() {
			Z80.R.e = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		}},
		{"RRA", nil},

		// 0x20
		{"JRNZn", func() {
			i := int(MMU.ReadByte(Z80.R.Pc))
			if i > 127 {
				i = -((^i+1)&0xFF)
			}
			Z80.R.Pc += 1
			Z80.R.M = 2
			if (Z80.R.f & 0x80) == 0x00 {
				pc := int(Z80.R.Pc) + i
				Z80.R.Pc = uint16(pc)
				Z80.R.M += 1
			}
		}},
		{"LDHLnn", func() {
			Z80.R.l = MMU.ReadByte(Z80.R.Pc)
			Z80.R.h = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		}},
		{"LDHLIA", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.a)
			Z80.R.l = Z80.R.l + 1  // TODO: test this overflows
			if Z80.R.l == 0 {
				Z80.R.h = Z80.R.h + 1  // TODO: test this overflows
			}
			Z80.R.M = 2
		}},
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
		{"CPL", func() {
			Z80.R.a ^= 0xFF
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},

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
		{"LDHLmr_b", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.b)
			Z80.R.M = 2
		}},
		{"LDHLmr_c", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.c)
			Z80.R.M = 2
		}},
		{"LDHLmr_d", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.d)
			Z80.R.M = 2
		}},
		{"LDHLmr_e", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.e)
			Z80.R.M = 2
		}},

		{"LDHLmr_h", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.h)
			Z80.R.M = 2
		}},
		{"LDHLmr_l", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.l)
			Z80.R.M = 2
		}},
		{"HALT", func() {
			Z80.Halt = true
			Z80.R.M = 1
		}},
		{"LDHLmr_a", func() {
			MMU.WriteByte((uint16(Z80.R.h) << 8) + uint16(Z80.R.l), Z80.R.a)
			Z80.R.M = 2
		}},

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
		{"ORr_b", func() {
			Z80.R.a |= Z80.R.b
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},
		{"ORr_c", func() {
			Z80.R.a |= Z80.R.c
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},
		{"ORr_d", func() {
			Z80.R.a |= Z80.R.d
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},
		{"ORr_e", func() {
			Z80.R.a |= Z80.R.e
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},

		{"ORr_h", func() {
			Z80.R.a |= Z80.R.h
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},
		{"ORr_l", func() {
			Z80.R.a |= Z80.R.l
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},
		{"ORHL", nil},
		{"ORr_a", func() {
			Z80.R.a |= Z80.R.a
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 1
		}},

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
		{"MAPcb", func() {
			i := MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.cb(i)
		}},

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
		{"ANDn", func() {
			Z80.R.a &= MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.M = 2
		}},
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
		{"LDAmm", func() {
			Z80.R.a = MMU.ReadByte(MMU.ReadWord(Z80.R.Pc))
			Z80.R.Pc += 2
			Z80.R.M = 4
		}},
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

	cbinstructions = [0x100]instruction {
		// CB00 
		{"RLCr_b", nil},
		{"RLCr_c", nil},
		{"RLCr_d", nil},
		{"RLCr_e", nil},

		{"RLCr_l", nil},
		{"RLCHL", nil},
		{"RLCr_a", nil},

		{"RRCr_b", nil},
		{"RRCr_c", nil},
		{"RRCr_d", nil},
		{"RRCr_e", nil},

		{"RRCr_h", nil},
		{"RRCr_l", nil},
		{"RRCHL", nil},
		{"RRCr_a", nil},

		// CB10
		{"RLr_b", nil},
		{"RLr_c", nil},
		{"RLr_d", nil},
		{"RLr_e", nil},

		{"RLr_h", nil},
		{"RLr_l", nil},
		{"RLHL", nil},
		{"RLr_a", nil},

		{"RRr_b", nil},
		{"RRr_c", nil},
		{"RRr_d", nil},
		{"RRr_e", nil},

		{"RRr_h", nil},
		{"RRr_l", nil},
		{"RRHL", nil},
		{"RRr_a", nil},

		// CB20
		{"SLAr_b", nil},
		{"SLAr_c", nil},
		{"SLAr_d", nil},
		{"SLAr_e", nil},

		{"SLAr_h", nil},
		{"SLAr_l", nil},
		{"XX", nil},
		{"SLAr_a", nil},

		{"SRAr_b", nil},
		{"SRAr_c", nil},
		{"SRAr_d", nil},
		{"SRAr_e", nil},

		{"SRAr_h", nil},
		{"SRAr_l", nil},
		{"XX", nil},
		{"SRAr_a", nil},

		// CB30
		{"SWAPr_b", nil},
		{"SWAPr_c", nil},
		{"SWAPr_d", nil},
		{"SWAPr_e", nil},

		{"SWAPr_h", nil},
		{"SWAPr_l", nil},
		{"XX", nil},
		{"SWAPr_a", nil},

		{"SRLr_b", func() {
			co := uint8(0)
			if Z80.R.b & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.b >>= 1
			Z80.R.f = 0x80
			if Z80.R.b != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},
		{"SRLr_c", func() {
			co := uint8(0)
			if Z80.R.c & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.c >>= 1
			Z80.R.f = 0x80
			if Z80.R.c != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},
		{"SRLr_d", func() {
			co := uint8(0)
			if Z80.R.d & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.d >>= 1
			Z80.R.f = 0x80
			if Z80.R.d != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},
		{"SRLr_e", func() {
			co := uint8(0)
			if Z80.R.e & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.e >>= 1
			Z80.R.f = 0x80
			if Z80.R.e != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},

		{"SRLr_h", func() {
			co := uint8(0)
			if Z80.R.h & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.h >>= 1
			Z80.R.f = 0x80
			if Z80.R.h != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},
		{"SRLr_l", func() {
			co := uint8(0)
			if Z80.R.l & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.l >>= 1
			Z80.R.f = 0x80
			if Z80.R.l != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},
		{"XX", nil},
		{"SRLr_a", func() {
			co := uint8(0)
			if Z80.R.a & 0x1 != 0 {
				co = 0x10
			}
			Z80.R.a >>= 1
			Z80.R.f = 0x80
			if Z80.R.a != 0 {
				Z80.R.f = 0
			}
			Z80.R.f = (Z80.R.f & 0xEF) + co
			Z80.R.M = 2
		}},

		// CB40
		{"BIT0b", nil},
		{"BIT0c", nil},
		{"BIT0d", nil},
		{"BIT0e", nil},

		{"BIT0h", nil},
		{"BIT0l", nil},
		{"BIT0m", nil},
		{"BIT0a", nil},

		{"BIT1b", nil},
		{"BIT1c", nil},
		{"BIT1d", nil},
		{"BIT1e", nil},

		{"BIT1h", nil},
		{"BIT1l", nil},
		{"BIT1m", nil},
		{"BIT1a", nil},

		// CB50
		{"BIT2b", nil},
		{"BIT2c", nil},
		{"BIT2d", nil},
		{"BIT2e", nil},

		{"BIT2h", nil},
		{"BIT2l", nil},
		{"BIT2m", nil},
		{"BIT2a", nil},

		{"BIT3b", nil},
		{"BIT3c", nil},
		{"BIT3d", nil},
		{"BIT3e", nil},

		{"BIT3h", nil},
		{"BIT3l", nil},
		{"BIT3m", nil},
		{"BIT3a", nil},

		// CB60
		{"BIT4b", nil},
		{"BIT4c", nil},
		{"BIT4d", nil},
		{"BIT4e,", nil},
		{"BIT4h", nil},
		{"BIT4l", nil},
		{"BIT4m", nil},
		{"BIT4a,", nil},
		{"BIT5b", nil},
		{"BIT5c", nil},
		{"BIT5d", nil},
		{"BIT5e,", nil},
		{"BIT5h", nil},
		{"BIT5l", nil},
		{"BIT5m", nil},
		{"BIT5a", nil},

		// CB70
		{"BIT6b", nil},
		{"BIT6c", nil},
		{"BIT6d", nil},
		{"BIT6e", nil},
		{"BIT6h", nil},
		{"BIT6l", nil},
		{"BIT6m", nil},
		{"BIT6a", nil},
		{"BIT7b", nil},
		{"BIT7c", nil},
		{"BIT7d", nil},
		{"BIT7e", nil},
		{"BIT7h", nil},
		{"BIT7l", nil},
		{"BIT7m", nil},
		{"BIT7a", nil},

		// CB80
		{"RES0b", nil},
		{"RES0c", nil},
		{"RES0d", nil},
		{"RES0e", nil},
		{"RES0h", nil},
		{"RES0l", nil},
		{"RES0m", nil},
		{"RES0a", nil},
		{"RES1b", nil},
		{"RES1c", nil},
		{"RES1d", nil},
		{"RES1e", nil},
		{"RES1h", nil},
		{"RES1l", nil},
		{"RES1m", nil},
		{"RES1a", nil},

		// CB90
		{"RES2b", nil},
		{"RES2c", nil},
		{"RES2d", nil},
		{"RES2e", nil},
		{"RES2h", nil},
		{"RES2l", nil},
		{"RES2m", nil},
		{"RES2a", nil},
		{"RES3b", nil},
		{"RES3c", nil},
		{"RES3d", nil},
		{"RES3e", nil},
		{"RES3h", nil},
		{"RES3l", nil},
		{"RES3m", nil},
		{"RES3a", nil},

		// CBA0
		{"RES4b", nil},
		{"RES4c", nil},
		{"RES4d", nil},
		{"RES4e", nil},
		{"RES4h", nil},
		{"RES4l", nil},
		{"RES4m", nil},
		{"RES4a", nil},
		{"RES5b", nil},
		{"RES5c", nil},
		{"RES5d", nil},
		{"RES5e", nil},
		{"RES5h", nil},
		{"RES5l", nil},
		{"RES5m", nil},
		{"RES5a", nil},

		// CBB0
		{"RES6b", nil},
		{"RES6c", nil},
		{"RES6d", nil},
		{"RES6e", nil},
		{"RES6h", nil},
		{"RES6l", nil},
		{"RES6m", nil},
		{"RES6a", nil},
		{"RES7b", nil},
		{"RES7c", nil},
		{"RES7d", nil},
		{"RES7e", nil},
		{"RES7h", nil},
		{"RES7l", nil},
		{"RES7m", nil},
		{"RES7a", nil},

		// CBC0
		{"SET0b", nil},
		{"SET0c", nil},
		{"SET0d", nil},
		{"SET0e", nil},
		{"SET0h", nil},
		{"SET0l", nil},
		{"SET0m", nil},
		{"SET0a", nil},
		{"SET1b", nil},
		{"SET1c", nil},
		{"SET1d", nil},
		{"SET1e", nil},
		{"SET1h", nil},
		{"SET1l", nil},
		{"SET1m", nil},
		{"SET1a", nil},

		// CBD0
		{"SET2b", nil},
		{"SET2c", nil},
		{"SET2d", nil},
		{"SET2e", nil},
		{"SET2h", nil},
		{"SET2l", nil},
		{"SET2m", nil},
		{"SET2a", nil},
		{"SET3b", nil},
		{"SET3c", nil},
		{"SET3d", nil},
		{"SET3e", nil},
		{"SET3h", nil},
		{"SET3l", nil},
		{"SET3m", nil},
		{"SET3a", nil},

		// CBE0
		{"SET4b", nil},
		{"SET4c", nil},
		{"SET4d", nil},
		{"SET4e", nil},
		{"SET4h", nil},
		{"SET4l", nil},
		{"SET4m", nil},
		{"SET4a", nil},
		{"SET5b", nil},
		{"SET5c", nil},
		{"SET5d", nil},
		{"SET5e", nil},
		{"SET5h", nil},
		{"SET5l", nil},
		{"SET5m", nil},
		{"SET5a", nil},

		// CBF0
		{"SET6b", nil},
		{"SET6c", nil},
		{"SET6d", nil},
		{"SET6e", nil},
		{"SET6h", nil},
		{"SET6l", nil},
		{"SET6m", nil},
		{"SET6a", nil},
		{"SET7b", nil},
		{"SET7c", nil},
		{"SET7d", nil},
		{"SET7e", nil},
		{"SET7h", nil},
		{"SET7l", nil},
		{"SET7m", nil},
		{"SET7a", nil},
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

func (z z80) cb(op uint8) {
	cbinstr := cbinstructions[op]
	if cbinstr.f == nil {
		log.Panicf("z80: nil cbinstruction for op [0x%x] %q\n", op, cbinstr.name)
	}
	log.Printf("z80: cbop [0x%x] %q\n", op, cbinstr.name)
	cbinstr.f.(func())()
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
