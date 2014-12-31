package main

import (
	"github.com/dominichamon/goboy/mmu"
	"github.com/dominichamon/goboy/z80"
)

func main() {
	z80.Reset()
	mmu.Reset()

	z80.Exec()
	z80.Exec()
}
