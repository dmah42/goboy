package goboy

import (
	"log"
)

type registers struct {
	a, b, c, d, e, h, l, f byte
	Pc, sp, i, r           int // 16-bit
	M                      int
	Ime                    int
}

type z80 struct {
	R	registers
	M	int
	Halt	bool
}

func makeZ80() z80 {
	var z z80
	z.R.Ime = 1
	return z
}

var (
	opmap = map[string]interface{}{
		"NOP": func() { Z80.R.M = 1 },
		"LDBCnn": func() {
			Z80.R.c = MMU.ReadByte(Z80.R.Pc)
			Z80.R.b = MMU.ReadByte(Z80.R.Pc + 1)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},
		"LDBCmA": func() {
			MMU.WriteByte((int(Z80.R.b)<<8)+int(Z80.R.c), Z80.R.a)
			Z80.R.M = 2
		},
		"INCBC": func() {
			Z80.R.c = (Z80.R.c + 1) & 0xFF
			if Z80.R.c == 0 {
				Z80.R.b = (Z80.R.b + 1) & 0xFF
			}
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
		"LDSPnn": func() {
			Z80.R.sp = MMU.ReadWord(Z80.R.Pc)
			Z80.R.Pc += 2
			Z80.R.M = 3
		},
		"JPnn": func() {
			Z80.R.Pc = MMU.ReadWord(Z80.R.Pc)
			Z80.R.M = 3
		},
	}
)

func (z z80) Call(opstr string) {
	if f, ok := opmap[opstr]; ok {
		f.(func())()
	} else {
		log.Panic("z80: No entry in opmap for op ", opstr)
	}
}
