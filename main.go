package main

import (
	//"fmt"
	emu "alex/CSProject/chip8"
	//"time"
)

func main() {
	emu.Init()
	c8 := emu.NewVM()
	//c8 := emu.NewVM()
	//fmt.Printf("%v", c8.vm.mem)
//	for range emu.Clock(time.Second/700) {
		c8.FDE()
		c8.FDE()
		c8.FDE()
//	}
}

