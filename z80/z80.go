package z80

import (
	"log"

	"github.com/dominichamon/goboy/mmu"
)

type registers struct {
	a, b, c, d, e, h, l, f byte
	pc, sp, i, r           int // 16-bit
	m                      int
	ime                    int
}

var (
	r          registers
	clock_m    int
	halt, stop bool

	opmap = map[opcode]interface{}{
		NOP: func() { r.m = 1 },
		LDBCnn: func() {
			r.c = mmu.ReadByte(r.pc)
			r.b = mmu.ReadByte(r.pc + 1)
			r.pc += 2
			r.m = 3
		},
		LDBCmA: func() {
			mmu.WriteByte((int(r.b)<<8)+int(r.c), r.a)
			r.m = 2
		},
		INCBC: func() {
			r.c = (r.c + 1) & 0xFF
			if r.c == 0 {
				r.b = (r.b + 1) & 0xFF
			}
			r.m = 1
		},
		INCr_b: func() {
			r.b = (r.b + 1) & 0xFF
			r.f = 0
			if r.b == 0 {
				r.f = 0x80
			}
			r.m = 1
		},
	}
)

func Reset() {
	r.a = 0
	r.b = 0
	r.c = 0
	r.d = 0
	r.e = 0
	r.h = 0
	r.l = 0
	r.f = 0

	r.sp = 0
	r.pc = 0
	r.i = 0
	r.r = 0

	r.m = 0
	clock_m = 0
	r.ime = 1

	halt = false
	stop = false

	log.Println("z80: Reset")
}

func Exec() {
	r.r = (r.r + 1) & 0xFF
	op := opcode(mmu.ReadByte(r.pc))
	if f, ok := opmap[op]; ok {
		f.(func())()
	} else {
		log.Panic("z80: No entry in opmap for op ", op)
	}
	r.pc = (r.pc + 1) & 0xFFFF
	clock_m += r.m
}
