package goboy

import (
	"log"
)

type clock struct {
	main, sub, div int
}

type timer struct {
	div, tma, tima, tac byte
	c clock
}

func (t *timer) step() {
	t.tima += 1
	t.c.main = 0
	if t.tima > 0xFF {
		t.tima = t.tma
		MMU.If |= 4
	}
}

func (t *timer) Inc() {
	t.c.sub += Z80.R.M
	if t.c.sub > 3 {
		t.c.main += 1
		t.c.sub -= 4

		t.c.div += 1
		if t.c.div == 0x10 {
			t.c.div = 0
			t.div = (t.div + 1) & 0xFF
		}
	}

	if (t.tac & 0x4) != 0{
		switch t.tac & 3 {
			case 0x0:
				if t.c.main >= 0x40 {
					t.step()
				}
			case 0x1:
				if t.c.main >= 0x1 {
					t.step()
				}
			case 0x2:
				if t.c.main >= 0x4 {
					t.step()
				}
			case 0x3:
				if t.c.main >= 0x10 {
					t.step()
				}
		}
	}
}

func (t timer) ReadByte(addr uint16) uint8 {
	switch addr {
		case 0xFF04:
			return t.div
		case 0xFF05:
			return t.tima
		case 0xFF06:
			return t.tma
		case 0xFF07:
			return t.tac
	}
	log.Panic("timer: Failed to read from address ", addr)
	return 0
}

func (t *timer) WriteByte(addr uint16, value uint8) {
	switch addr {
		case 0xFF04:
			t.div = 0
		case 0xFF05:
			t.tima = value
		case 0xFF06:
			t.tma = value
		case 0xFF07:
			t.tac = value & 0x7
	}
}
