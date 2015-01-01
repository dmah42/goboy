package goboy

import (
	"log"
	"time"
)

var (
	MMU = makeMMU()
	Z80 = makeZ80()
	Key = makeKey()
	GPU = makeGPU()

	Timer	timer

	Run = false
)

func frame() {
	fclock := Z80.M + 17556
	// t0 := time.Now()
	for Z80.M < fclock {
		if Z80.Halt {
			Z80.Call("NOP")
		} else {
			Z80.R.r = (Z80.R.r + 1) & 0xFF
			op := MMU.ReadByte(Z80.R.Pc)
			//log.Printf("op [%x] %q\n", op, Opcodes[op])
			Z80.Call(Opcodes[op])
			Z80.R.Pc = (Z80.R.Pc + 1) & 0xFFFF
		}

		// check for interrupts
		if Z80.R.Ime != 0 && MMU.Ie != 0 && MMU.If != 0 {
			log.Printf("int %x %x %x\n", Z80.R.Ime, MMU.Ie, MMU.If)
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
		GPU.Checkline()
		Timer.Inc()

		// log.Printf("z80: %+v\n", Z80)
		// log.Printf("mmu: %+v %+v %+v\n", MMU.inbios, MMU.Ie, MMU.If)
		// log.Printf("timer: %+v\n", Timer)
		// log.Printf("key: %+v\n", Key)
		// log.Printf("gpu: %+v %+v %+v\n", GPU.linemode, GPU.modeclocks, GPU.curline)
		// log.Println("-----------------")

		//time.Sleep(250 * time.Millisecond)
	}
	//log.Printf("fps: %.3f\n", 1.0 / time.Since(t0).Seconds())
}

func Loop(rom string) {
	Z80.R.Pc = 0x100
	Z80.R.sp = 0xFFFE
	Z80.R.a = 1
	Z80.R.c = 0x13
	Z80.R.e = 0xD8
	Z80.R.h = 0x01
	Z80.R.l = 0x4D

	MMU.inbios = false

	MMU.Load(rom)

	log.Println("goboy: Starting loop")
	for {
		select {
			case <-time.After(time.Millisecond):
				if Run { frame() }
		}
	}
}
