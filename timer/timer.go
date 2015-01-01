package timer

import (
	"log"

	"github.com/dominichamon/goboy/mmu"
	"github.com/dominichamon/goboy/z80"
)

type clock struct {
	main, sub, div int
}

var (
	div, tma, tima, tac byte
	c clock
)

func Reset() {
	div = 0
	tma = 0
	tima = 0
	tac = 0
	c.main = 0
	c.sub = 0
	c.div = 0
	log.Println("timer: Reset")
}

func step() {
	tima += 1
	c.main = 0
	if tima > 0xFF {
		tima = tma
		mmu.If |= 4
	}
}

func Inc() {
	c.sub += z80.R.M
	if c.sub > 3 {
		c.main += 1
		c.sub -= 4

		c.div += 1
		if c.div == 0x10 {
			c.div = 0
			div = (div + 1) & 0xFF
		}
	}

	if (tac & 0x4) != 0{
		switch tac & 3 {
			case 0x0:
				if c.main >= 0x40 {
					step()
				}
			case 0x1:
				if c.main >= 0x1 {
					step()
				}
			case 0x2:
				if c.main >= 0x4 {
					step()
				}
			case 0x3:
				if c.main >= 0x10 {
					step()
				}
		}
	}
}

func ReadByte(addr int) byte {
	switch addr {
		case 0xFF04:
			return div
		case 0xFF05:
			return tima
		case 0xFF06:
			return tma
		case 0xFF07:
			return tac
	}
	log.Panic("timer: Failed to read from address ", addr)
	return 0
}

func WriteByte(addr int, value byte) {
	switch addr {
		case 0xFF04:
			div = 0
		case 0xFF05:
			tima = value
		case 0xFF06:
			tma = value
		case 0xFF07:
			tac = value & 0x7
	}
}
