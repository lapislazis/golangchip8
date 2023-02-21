package main

import (
	emu "alex/CSProject/chip8"
	"time"
)

func main() {
	emu.Init()
	c8 := emu.NewVM()
	for range emu.Clock(time.Second/700) {
		c8.FDE()
	}
}

