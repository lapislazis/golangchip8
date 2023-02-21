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
	//FDE Testing (setting memory)
	//vm.mem[512] = 0x15
	//vm.mem[513] = 0x00
	//vm.mem[1280] = 0x12
	//vm.mem[1281] = 0x00
	//vm.mem[512] = 0x61
	//vm.mem[513] = 0xFD
	//vm.mem[514] = 0x62
	//vm.mem[515] = 0x03
	//vm.mem[516] = 0x81
	//vm.mem[517] = 0x24
	

	return vm
}
func Clock(d time.Duration) <-chan time.Time {
	ch := make(chan time.Time, 1)
	go func() {
		for {
		time.Sleep(d)
		ch <- time.Now()
	}
	close(ch)
	}()
	return ch
}
	

//Fetch-Decode-Execute Cycle
func (vm *chip8) FDE() {

	//Sets opcode variable to whats in mem, shift left, OR whats in mem+1
	vm.op = (uint16(vm.mem[vm.pc]) << 8) | uint16(vm.mem[vm.pc+1])

	//Match opcode to first nibble
	switch vm.op & 0xF000 {
	//First nibble is 0001
	case 0x1000: //0x1NNN jumps to NNN on memory
		//Program counter = last 3 nibbles of opcode
		vm.pc = vm.op & 0x0FFF
		fmt.Printf("Current address in memory is %v\n", vm.pc)
	case 0x3000: //0x3XNN Sets the value of VX to NN
		//Skips if VX is already NN
		if uint16(vm.v[(vm.op & 0x0F00) >> 8]) == vm.op & 0x00FF {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x6000: //0x6XNN sets value of VX to NN
		vm.v[(vm.op & 0x0F00) >> 8] = byte(vm.op & 0x00FF)
		fmt.Printf("Value of V%d ", ((vm.op & 0x0F00) >> 8))
		fmt.Println("is set to ", (vm.v[(vm.op & 0x0F00) >> 8]))
		vm.pc += 2
	case 0x8000: 
		switch vm.op & 0x000F { //Check the last nibble 
		case 0x0004: //0x8XY4 Adds VY to VX, sets VF to 1 if overflow
			if vm.v[(vm.op & 0x00F0) >> 4] > (0xFF - vm.v[(vm.op & 0x0F00) >> 8]) {
				vm.v[15] = 1 //Set VF to 1
				fmt.Println("Overflow! VF is set to ", vm.v[15])
			} else {
				vm.v[15] = 0
			}
			vm.v[(vm.op & 0x0F00) >> 8] += vm.v[(vm.op & 0x00F0) >> 4] //Otherwise, add VX and VY
			fmt.Println("The value of V1 after addition is ", vm.v[1])

		}
	default:
		fmt.Printf("Invalid opcode 0x%X\n", vm.op)
	}
	
}

	





