// This program is an interpreter for the instruction set Chip-8, written in go.
// This project was developed as part of an A-Level Computer Science assignment,
// and is purely for educational purposes. CHIP-8 was implemented in on systems
// with 4kb of memory such as the Cosmac VIP, and was most used for game dev.
package chip8

//Required libraries
import (
	"fmt"
	"time"
)

// Format of the fontset
var fontSet = []uint8{
	0xF0, 0x90, 0x90, 0x90, 0xF0, //0
	0x20, 0x60, 0x20, 0x20, 0x70, //1
	0xF0, 0x10, 0xF0, 0x80, 0xF0, //2
	0xF0, 0x10, 0xF0, 0x10, 0xF0, //3
	0x90, 0x90, 0xF0, 0x10, 0x10, //4
	0xF0, 0x80, 0xF0, 0x10, 0xF0, //5
	0xF0, 0x80, 0xF0, 0x90, 0xF0, //6
	0xF0, 0x10, 0x20, 0x40, 0x40, //7
	0xF0, 0x90, 0xF0, 0x90, 0xF0, //8
	0xF0, 0x90, 0xF0, 0x10, 0xF0, //9
	0xF0, 0x90, 0xF0, 0x90, 0x90, //A
	0xE0, 0x90, 0xE0, 0x90, 0xE0, //B
	0xF0, 0x80, 0x80, 0x80, 0xF0, //C
	0xE0, 0x90, 0x90, 0x90, 0xE0, //D
	0xF0, 0x80, 0xF0, 0x80, 0xF0, //E
	0xF0, 0x80, 0xF0, 0x80, 0x80, //F
}

//Format of the emulator
type chip8 struct {
	//System memory. See diagram in documentation
	mem [4096]byte

	//General purpose registers
	v [16]byte

	//Index register
	I uint16

	//Program counter
	pc uint16

	//Current opcode
	op uint16
	
	//Program counter stack
	stack [16]uint16

	//Stack pointer 
	sp uint16

	//Delay timer
	delayTime byte

	//CPU clock
	clock *time.Ticker

	// TODO: timers, keypad, graphics, audio, shutdown
}

func Init() {
	//Check that the package is imported by main
	fmt.Println("Chip8 is initialised")
}
//Initialise emulator instance
func NewVM()  chip8 {
	vm := chip8{
		mem:	[4096]byte{},
		v:	[16]byte{},
		pc:	0x200,
		stack: [16]uint16{},
		clock: time.NewTicker(time.Second / 700),
		}
	//Load fontset
	for i := 0; i < 80; i++ {
		vm.mem[i] = fontSet[i]
		}	

	//Print	empty variables test
	fmt.Printf("V1-16 contain:\n")
	for i := 0; i < 16; i++ {
		fmt.Printf("0x%X ", vm.v[i])
	}
	fmt.Printf("\nThe stack contains:\n")
	for i := 0; i < 16; i++ {
		fmt.Printf("0x%X\n", vm.stack[i])
	}
	fmt.Printf("Program counter contains:\n0x%X\n", vm.pc)
	fmt.Printf("Index register contains:\n0x%X\n", vm.I) 
	fmt.Printf("Opcode contains:\n0x%X\n", vm.op)

	//Print mem test
	fmt.Printf("Memory contains:\n")
	for i := 0; i < 80; i++ {	
		fmt.Printf("0x%X ", vm.mem[i])
	}
	
	//Assign variables test
	fmt.Printf("\nV1-16 have been assigned:\n")
	for i := 0; i < 16; i++ {
		vm.v[i] = byte(i)
		fmt.Printf("0x%X ", vm.v[i])
	}
	fmt.Printf("\nStack has been assigned:\n")
	for i := 0; i < 16; i++ {
		vm.stack[i] = uint16(i)
		fmt.Printf("0x%X ", vm.stack[i])
	}
	vm.I = 65535
	fmt.Printf("\nIndex counter has been assigned:\n0x%X\n", vm.I)
	vm.op = 0x00e0
	fmt.Printf("Opcode has been assigned:\n0x%X\n", vm.op)
	
	//Overflow variable test VERY UNEXPECTED RESULT (resolved)
	vm.v[15] += 255
	fmt.Printf("Variable overflow:\n0x%X\n", vm.v[15])
	vm.I += 1
	fmt.Printf("Index overflow:\n0x%X\n", vm.I)

	//Check clockspeed is accurate
	
	//Doesn't work
	//done := make(chan bool)
	//go func() {
	//	for {
	//		select {
	//		case <-done:
	//			return
	//		case tick:= <-vm.clock.C:
	//			tick += 1
	//		}
	//	}
	//}()
	//time.Sleep(1 * time.Second)
	//vm.clock.Stop
	//done <- true
	//fmt.Printf("\nFinal tick was:\n%v", tick)
	
	//Doesn't end loop
	//for start := time.Now(); time.Since(start) < time.Second; {
	//	for _ = range vm.clock.C {
	//		t++
	//	}
	//}

	t := 0
	loop:
		for timeout := time.After(5 * time.Second); ; {
			select {
			case <-timeout: 
				break loop
			case <-vm.clock.C:
				t++
			}
		}
	fmt.Printf("Final tick was:\n%v\n", t)

	return vm
}




