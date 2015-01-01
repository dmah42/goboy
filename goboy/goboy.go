package goboy

import (
	"log"
	"time"
)

var (
	MMU = makeMMU()
	Z80 = makeZ80()
	Key = makeKey()

	Timer	timer
)

func init() {
	Z80.Boot()
	MMU.Boot()
}

func LoadROM(rom string) {
	MMU.Load(rom)
}

func Run() {
	// run
	Z80.Stop = false
	for !Z80.Stop {
		select {
			case <-time.After(time.Millisecond):
				fclock := Z80.M + 17556
				for Z80.M < fclock {
					if Z80.Halt {
						Z80.Call("NOP")
					} else {
						//z80.r.r = (z80.r.r + 1) & 0xFF
						op := MMU.ReadByte(Z80.R.Pc)
						Z80.Call(Opcodes[op])
						Z80.R.Pc = (Z80.R.Pc + 1) & 0xFFFF
					}

					// check for interrupts
					if Z80.R.Ime != 0 && MMU.Ie != 0 && MMU.If != 0 {
						Z80.Halt = false
						Z80.R.Ime = 0
						ifired := MMU.Ie & MMU.If
						if (ifired & 1) != 0 {
							MMU.If &= 0xFE
							Z80.Call("RST40")
						} else if (ifired & 2) != 0 {
							MMU.If &= 0xFD
							Z80.Call("RST48")
						} else if (ifired & 4) != 0 {
							MMU.If &= 0xFB
							Z80.Call("RST50")
						} else if (ifired & 8) != 0 {
							MMU.If &= 0xF7
							Z80.Call("RST58")
						} else if (ifired & 16) != 0 {
							MMU.If &= 0xEF
							Z80.Call("RST60")
						} else {
							Z80.R.Ime = 1
						}
					}
					Z80.M += Z80.R.M
					// TODO: gpu.checkline()
					Timer.Inc()
					log.Printf("z80: %+v\n", Z80)
					log.Printf("mmu: %+v %+v %+v\n", MMU.inbios, MMU.Ie, MMU.If)
					log.Printf("timer: %+v\n", Timer)

					time.Sleep(time.Millisecond * 250)
				}
		}
	}

}
