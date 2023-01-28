package main

import (
	//"fmt"
	emu "alex/CSProject/chip8"
)

func main() {
	emu.Init()
	emu.NewVM()
	fmt.Println("Testing")
	//c8 := emu.NewVM()
	//fmt.Printf("%v", c8.vm.mem)
}

