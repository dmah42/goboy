package main

import (
	"log"
	"os"
	"time"

	"github.com/dominichamon/goboy/goboy"
)

const (
	// rate = time.Millisecond * 16
	rate = time.Second * 5
)

func main() {
	if len(os.Args) < 2 {
		log.Panic("no ROM file selected")
	}
	goboy.LoadROM(os.Args[1])

	goboy.Run()
}
