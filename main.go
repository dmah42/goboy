package main

import (
	"log"
	"os"
	"time"
	
	"github.com/dominichamon/goboy/mmu"
	"github.com/dominichamon/goboy/z80"
)

func main() {
	z80.Reset()
	mmu.Reset()

	z80.Boot()
	mmu.Boot()

	if len(os.Args) < 2 {
		log.Panic("no ROM file selected")
	}
	mmu.Load(os.Args[1])

	// run
	z80.Stop = false
	select {
		case <-time.After(time.Millisecond):
			fclock := z80.Clock_m + 17556
			for z80.Clock_m < fclock {
				if z80.Halt {
					z80.Call("NOP")
				} else {
					//z80.r.r = (z80.r.r + 1) & 0xFF
					op := mmu.ReadByte(z80.R.Pc)
					z80.Call(z80.Opcodes[op])
					z80.R.Pc = (z80.R.Pc + 1) & 0xFFFF
				}

				// check for interrupts
				if z80.R.Ime != 0 && mmu.Ie != 0 && mmu.If != 0 {
					z80.Halt = false
					z80.R.Ime = 0
					ifired := mmu.Ie & mmu.If
					if (ifired & 1) != 0 {
						mmu.If &= 0xFE
						z80.Call("RST40")
					} else if (ifired & 2) != 0 {
						mmu.If &= 0xFD
						z80.Call("RST48")
					} else if (ifired & 4) != 0 {
						mmu.If &= 0xFB
						z80.Call("RST50")
					} else if (ifired & 8) != 0 {
						mmu.If &= 0xF7
						z80.Call("RST58")
					} else if (ifired & 16) != 0 {
						mmu.If &= 0xEF
						z80.Call("RST60")
					} else {
						z80.R.Ime = 1
					}
				}
				z80.Clock_m += z80.R.M
				// TODO: gpu.checkline()
				// TODO: timer.inc()
			}
	}
}
