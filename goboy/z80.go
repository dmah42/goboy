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

	instructions = map[string]interface{}{
		"NOP": func() { Z80.R.M = 1 },
		"HALT": func() {
			Z80.Halt = true
			Z80.R.M = 1
		},

		"DI": func() {
			Z80.R.Ime = false
			Z80.R.M = 1
		},
		"EI": func() {
			Z80.R.Ime = true
			Z80.R.M = 1
		},

		"LDrr_aa": func() {
			Z80.R.a = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_ab": func() {
			Z80.R.a = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_ac": func() {
			Z80.R.a = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_ad": func() {
			Z80.R.a = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_ae": func() {
			Z80.R.a = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_ah": func() {
			Z80.R.a = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_al": func() {
			Z80.R.a = Z80.R.l
			Z80.R.M = 1
		},
		"LDrr_ba": func() {
			Z80.R.b = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_bb": func() {
			Z80.R.b = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_bc": func() {
			Z80.R.b = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_bd": func() {
			Z80.R.b = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_be": func() {
			Z80.R.b = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_bh": func() {
			Z80.R.b = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_bl": func() {
			Z80.R.b = Z80.R.l
			Z80.R.M = 1
		},
		"LDrr_ca": func() {
			Z80.R.c = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_cb": func() {
			Z80.R.c = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_cc": func() {
			Z80.R.c = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_cd": func() {
			Z80.R.c = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_ce": func() {
			Z80.R.c = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_ch": func() {
			Z80.R.c = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_cl": func() {
			Z80.R.c = Z80.R.l
			Z80.R.M = 1
		},
		"LDrr_da": func() {
			Z80.R.d = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_db": func() {
			Z80.R.d = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_dc": func() {
			Z80.R.d = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_dd": func() {
			Z80.R.d = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_de": func() {
			Z80.R.d = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_dh": func() {
			Z80.R.d = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_dl": func() {
			Z80.R.d = Z80.R.l
			Z80.R.M = 1
		},
		"LDrr_ea": func() {
			Z80.R.e = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_eb": func() {
			Z80.R.e = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_ec": func() {
			Z80.R.e = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_ed": func() {
			Z80.R.e = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_ee": func() {
			Z80.R.e = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_eh": func() {
			Z80.R.e = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_el": func() {
			Z80.R.e = Z80.R.l
			Z80.R.M = 1
		},
		"LDrr_ha": func() {
			Z80.R.h = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_hb": func() {
			Z80.R.h = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_hc": func() {
			Z80.R.h = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_hd": func() {
			Z80.R.h = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_he": func() {
			Z80.R.h = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_hh": func() {
			Z80.R.h = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_hl": func() {
			Z80.R.h = Z80.R.l
			Z80.R.M = 1
		},
		"LDrr_la": func() {
			Z80.R.l = Z80.R.a
			Z80.R.M = 1
		},
		"LDrr_lb": func() {
			Z80.R.l = Z80.R.b
			Z80.R.M = 1
		},
		"LDrr_lc": func() {
			Z80.R.l = Z80.R.c
			Z80.R.M = 1
		},
		"LDrr_ld": func() {
			Z80.R.l = Z80.R.d
			Z80.R.M = 1
		},
		"LDrr_le": func() {
			Z80.R.l = Z80.R.e
			Z80.R.M = 1
		},
		"LDrr_lh": func() {
			Z80.R.l = Z80.R.h
			Z80.R.M = 1
		},
		"LDrr_ll": func() {
			Z80.R.l = Z80.R.l
			Z80.R.M = 1
		},

		"LDrHLm_a": func() {
			Z80.R.a = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},
		"LDrHLm_b": func() {
			Z80.R.b = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},
		"LDrHLm_c": func() {
			Z80.R.c = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},
		"LDrHLm_d": func() {
			Z80.R.d = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},
		"LDrHLm_e": func() {
			Z80.R.e = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},
		"LDrHLm_h": func() {
			Z80.R.h = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},
		"LDrHLm_l": func() {
			Z80.R.l = MMU.ReadByte(uint16(Z80.R.h << 8) + uint16(Z80.R.l))
			Z80.R.M = 2
		},

		"LDBCnn": func() {
			Z80.R.c = MMU.ReadByte(Z80.R.Pc)
			Z80.R.b = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},

		"LDBCmA": func() {
			MMU.WriteByte((uint16(Z80.R.b)<<8)+uint16(Z80.R.c), Z80.R.a)
			Z80.R.M = 2
		},
		"LDDEmA": func() {
			MMU.WriteByte((uint16(Z80.R.d)<<8)+uint16(Z80.R.e), Z80.R.a)
			Z80.R.M = 2
		},

		"LDmmA": func() {
			MMU.WriteByte(MMU.ReadWord(Z80.R.Pc), Z80.R.a)
			Z80.R.Pc += 2
			Z80.R.M = 4
		},

		"LDBCnn": func() {
			Z80.R.b = MMU.ReadByte(Z80.R.Pc)
			Z80.R.c = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},
		"LDDEnn": func() {
			Z80.R.d = MMU.ReadByte(Z80.R.Pc)
			Z80.R.e = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},
		"LDHLnn": func() {
			Z80.R.l = MMU.ReadByte(Z80.R.Pc)
			Z80.R.h = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},
		"LDSPnn": func() {
			Z80.R.sp = MMU.ReadWord(Z80.R.Pc)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},

		"LDrn_a": func() {
			Z80.R.a = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},
		"LDrn_b": func() {
			Z80.R.b = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},
		"LDrn_c": func() {
			Z80.R.c = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},
		"LDrn_d": func() {
			Z80.R.d = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},
		"LDrn_e": func() {
			Z80.R.e = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},
		"LDrn_h": func() {
			Z80.R.h = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},
		"LDrn_l": func() {
			Z80.R.l = MMU.ReadByte(Z80.R.Pc)
			Z80.R.Pc += 1
			Z80.R.M = 2
		},

		"LDSPnn": func() {
			Z80.R.sp = MMU.ReadWord(Z80.R.Pc)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},

		"LDAIOn": func() {
			Z80.R.a = MMU.ReadByte(uint16(0xFF00) + uint16(MMU.ReadByte(Z80.R.Pc)))
			Z80.R.Pc += 1
			Z80.R.M = 3
		},
		"LDIOnA": func() {
			MMU.WriteByte(uint16(0xFF00) + uint16(MMU.ReadByte(Z80.R.Pc)), Z80.R.a)
			Z80.R.Pc += 1
			Z80.R.M = 3
		},
		"LDAIOC": func() {
			Z80.R.a = MMU.ReadByte(uint16(0xFF00) + uint16(Z80.R.c))
			Z80.R.M = 2
		},
		"LDIOCA": func() {
			MMU.WriteByte(uint16(0xFF00) + uint16(Z80.R.c), Z80.R.a)
			Z80.R.M = 2
		},

		"ADDr_a": func() {
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
		},

		"ADDHLBC": func() {
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
		},
		"ADDHLDE": func() {
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
		},
		"ADDHLHL": func() {
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
		},
		"ADDHLSP": func() {
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
		},
		"ADDSPn": func() {
			i := int16(MMU.ReadByte(Z80.R.Pc))
			if i > 127 {
				i = -((^i+1)&0xFF)
			}
			Z80.R.Pc += 1
			sp := int16(Z80.R.sp) + i
			Z80.R.sp = uint16(sp)
			Z80.R.M = 4
		},

		"SUBr_a": func() {
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
		},

		"CPn": func() {
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
		},

		"INCBC": func() {
			Z80.R.c = (Z80.R.c + 1) & 0xFF
			if Z80.R.c == 0 {
				Z80.R.b = (Z80.R.b + 1) & 0xFF
			}
			Z80.R.M = 1
		},
		"INCDE": func() {
			Z80.R.e = (Z80.R.e + 1) & 0xFF
			if Z80.R.e == 0 {
				Z80.R.d = (Z80.R.d + 1) & 0xFF
			}
			Z80.R.M = 1
		},
		"INCHL": func() {
			Z80.R.l = (Z80.R.l + 1) & 0xFF
			if Z80.R.l == 0 {
				Z80.R.h = (Z80.R.h + 1) & 0xFF
			}
			Z80.R.M = 1
		},
		"INCSP": func() {
			Z80.R.sp = (Z80.R.sp + 1) & 0xFFFF
			Z80.R.M = 1
		},

		"INCr_b": func() {
			Z80.R.b = (Z80.R.b + 1) & 0xFF
			Z80.R.f = 0
			if Z80.R.b == 0 {
				Z80.R.f = 0x80
			}
			Z80.R.M = 1
		},

		"PUSHBC": func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.b)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.c)
			Z80.R.M = 3
		},
		"PUSHDE": func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.d)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.e)
			Z80.R.M = 3
		},
		"PUSHHL": func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.h)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.l)
			Z80.R.M = 3
		},
		"PUSHAF": func() {
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.a)
			Z80.R.sp -= 1
			MMU.WriteByte(Z80.R.sp, Z80.R.f)
			Z80.R.M = 3
		},

		"POPBC": func() {
			Z80.R.c = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.b = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		},
		"POPDE": func() {
			Z80.R.e = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.d = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		},
		"POPHL": func() {
			Z80.R.l = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.h = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		},
		"POPAF": func() {
			Z80.R.f = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.a = MMU.ReadByte(Z80.R.sp)
			Z80.R.sp += 1
			Z80.R.M = 3
		},

		"JPnn": func() {
			Z80.R.Pc = MMU.ReadWord(Z80.R.Pc)
			Z80.R.M = 3
		},

		"CALLnn": func() {
			Z80.R.sp -= 2
			MMU.WriteWord(Z80.R.sp, Z80.R.Pc + 2)
			Z80.R.Pc = MMU.ReadWord(Z80.R.Pc)
			Z80.R.M = 5
		},

		"RET": func() {
			Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
			Z80.R.sp += 2
			Z80.R.M = 3
		},
		"RETI": func() {
			Z80.R.Ime = true
			Z80.loadRegisters()
			Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
			Z80.R.sp += 2
			Z80.R.M = 3
		},
		"RETNZ": func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x80) == 0 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		},
		"RETZ": func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x80) == 0x80 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		},
		"RETNC": func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x10) == 0 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		},
		"RETC": func() {
			Z80.R.M = 1
			if (Z80.R.f & 0x10) == 0x10 {
				Z80.R.Pc = MMU.ReadWord(Z80.R.sp)
				Z80.R.sp += 2
				Z80.R.M += 2
			}
		},
	}
)

func (z z80) Call(opstr string) {
	if f, ok := instructions[opstr]; ok {
		f.(func())()
	} else {
		log.Panic("z80: No instruction for op ", opstr)
	}
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
