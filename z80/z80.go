package z80

import (
	"log"

	"github.com/dominichamon/goboy/mmu"
)

type Registers struct {
	a, b, c, d, e, h, l, f byte
	Pc, sp, i, r           int // 16-bit
	M                      int
	Ime                    int
}

var (
	R          Registers
	Clock_m    int
	Halt, Stop bool

	opmap = map[string]interface{}{
		"NOP": func() { R.M = 1 },
		"LDBCnn": func() {
			R.c = mmu.ReadByte(R.Pc)
			R.b = mmu.ReadByte(R.Pc + 1)
			R.Pc += 2
			R.M = 3
		},
		"LDBCmA": func() {
			mmu.WriteByte((int(R.b)<<8)+int(R.c), R.a)
			R.M = 2
		},
		"INCBC": func() {
			R.c = (R.c + 1) & 0xFF
			if R.c == 0 {
				R.b = (R.b + 1) & 0xFF
			}
			R.M = 1
		},
		"INCr_b": func() {
			R.b = (R.b + 1) & 0xFF
			R.f = 0
			if R.b == 0 {
				R.f = 0x80
			}
			R.M = 1
		},
		"LDSPnn": func() {
			R.sp = mmu.ReadWord(R.Pc)
			R.Pc += 2
			R.M = 3
		},
		"JPnn": func() {
			R.Pc = mmu.ReadWord(R.Pc)
			R.M = 3
		},
	}
)

func Reset() {
	R.a = 0
	R.b = 0
	R.c = 0
	R.d = 0
	R.e = 0
	R.h = 0
	R.l = 0
	R.f = 0

	R.sp = 0
	R.Pc = 0
	R.i = 0
	R.r = 0

	R.M = 0
	Clock_m = 0
	R.Ime = 1

	Halt = false
	Stop = false

	log.Println("z80: Reset")
}

func Boot() {
	R.Pc = 0x100
	R.sp = 0xFFFE
	R.a = 1
	R.c = 0x13
	R.e = 0xD8
}

func Call(opstr string) {
	log.Printf("z80: Call %q\n", opstr)
	if f, ok := opmap[opstr]; ok {
		f.(func())()
	} else {
		log.Panic("z80: No entry in opmap for op ", opstr)
	}
}
