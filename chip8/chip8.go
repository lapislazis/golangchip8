// This program is an interpreter for the instruction set Chip-8, written in go.
// This project was developed as part of an A-Level Computer Science assignment,
// and is purely for educational purposes. CHIP-8 was implemented in on systems
// with 4kb of memory such as the Cosmac VIP, and was most used for game dev.
package chip8

//Required libraries
import (
	"fmt"
	"time"
	"os"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
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

	//Sound timer
	soundTime byte

	//CPU clock
	clock *time.Ticker

	//Graphics
	gfx [32][64]uint8 //size of display

	//Draw flag
	drawFlag bool

	//Keyboard
	key [16]uint8 //CHIP-8 keypad had 16 keys

	//Channel to check for audio events
	audioChan chan struct{}

	// TODO: timers, keypad, graphics, audio, shutdown
}

func Init() {
	//Check that the package is imported by main
	fmt.Println("Chip8 is initialised")
}
//Initialise emulator instance
func NewVM()  chip8 {
	vm := chip8{
		mem:		[4096]byte{},
		v:		[16]byte{},
		pc:		0x200,
		stack: 		[16]uint16{},
		clock: 		time.NewTicker(time.Second / 700),
		drawFlag:	true,
		audioChan:	make(chan struct{}),
		}
	//Load fontset
	for i := 0; i < 80; i++ {
		vm.mem[i] = fontSet[i]
		}
	vm.v[1] = 1

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

func (vm *chip8) Buffer() [32][64]uint8 {
	return vm.gfx
}

func (vm *chip8) Draw() bool {
	df := vm.drawFlag
	vm.drawFlag = false
	return df
}

//Checks a key is pressed (check format in main)
func (vm *chip8) Key(num uint8, down bool) {
	if down {
		vm.key[num] = 1
	} else {
		vm.key[num] = 0
	}
}

func (vm *chip8) delayTimeTick() {
	if vm.delayTime > 0 {
		vm.delayTime --
	}
}

func (vm *chip8) Audio() {
	f, err := os.Open("beep.mp3")
	if err != nil {
		return
	}
	
	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return
	}
	defer streamer.Close()

	speaker.Init(
		format.SampleRate,
		format.SampleRate.N(time.Second/10),
	)

	for range vm.audioChan {
		speaker.Play(streamer)
		fmt.Printf("Beep!!")
	}
}

func (vm *chip8) soundTimeTick() {
	if vm.soundTime > 0 {
		if vm.soundTime == 1 {
			vm.audioChan <- struct{}{}
		}
		vm.soundTime--
	}
}			

//Fetch-Decode-Execute Cycle
func (vm *chip8) FDE() {
	//Sets opcode variable to whats in mem, shift left, OR whats in mem+1
	vm.op = (uint16(vm.mem[vm.pc]) << 8) | uint16(vm.mem[vm.pc+1])

	//Match opcode to first nibble
	switch vm.op & 0xF000 {
	//First nibble is 0000
	case 0x0000: 
		switch vm.op & 0x000F {
		case 0x0000: //0x0000 clears screen
			for i := 0; i < len(vm.gfx); i++ {
				for j := 0; j < len(vm.gfx[i]); j++ {
					vm.gfx[i][j] = 0x0
				}
			}
			vm.drawFlag = true
			vm.pc = vm.pc + 2
		//0x00EE here
		default: 
			fmt.Printf("Invalid opcode 0x%X\n", vm.op)
		}	
	//First nibble is 0001
	case 0x1000: //0x1NNN jumps to NNN on memory
		//Program counter = last 3 nibbles of opcode
		vm.pc = vm.op & 0x0FFF
	case 0x3000: //0x3XNN Sets the value of VX to NN
		//Skips if VX is already NN
		if uint16(vm.v[(vm.op & 0x0F00) >> 8]) == vm.op & 0x00FF {
			vm.pc += 4
		} else {
			vm.pc += 2
		}
	case 0x6000: //0x6XNN sets value of VX to NN
		vm.v[(vm.op & 0x0F00) >> 8] = byte(vm.op & 0x00FF)
		vm.pc += 2
	case 0x8000: 
		switch vm.op & 0x000F { //Check the last nibble 
		case 0x0004: //0x8XY4 Adds VY to VX, sets VF to 1 if overflow
			if vm.v[(vm.op & 0x00F0) >> 4] > (0xFF - vm.v[(vm.op & 0x0F00) >> 8]) {
				vm.v[15] = 1 //Set VF to 1
			} else {
				vm.v[15] = 0
			}
			vm.v[(vm.op & 0x0F00) >> 8] += vm.v[(vm.op & 0x00F0) >> 4] //Otherwise, add VX and VY
		}
	case 0xA000: //0xANNN Sets I to address NNN
		vm.I = vm.op & 0x0FFF
		vm.pc = vm.pc + 2
	case 0xD000: //0xDXYN Draws sprite of length N in memory starting at I at co-ords (VX, VY)
		x := vm.v[(vm.op & 0x0F00) >> 8] 
		y := vm.v[(vm.op & 0x00F0) >> 4]
		h := (vm.op & 0x000F)
		vm.v[0xF] = 0
		var j uint16 = 0
		var i uint16 = 0
		for j = 0; j < h; j++ {
			pixel := vm.mem[vm.I+j] //Pixel = current y-axis in memory
			for i = 0; i < 8; i++ { //For each x up to 8 pixels wide
				if (pixel & (0x80 >> i)) != 0 {	//If current byte isn't 0
					if vm.gfx[(y + uint8(j))][x + uint8(i)] == 1 { //And if that byte in display is already on
						vm.v[0xF] = 1 //Collision flag on
					}
					vm.gfx[(y + uint8(j))][(x + uint8(i))] ^= 1 //Either way, set the pixel to its value OR 1 
				}
			}
		}
		vm.drawFlag = true
		vm.pc = vm.pc + 2

	case 0xE000:
		switch vm.op & 0x00FF {
		case 0x009E:
			if vm.key[vm.v[(vm.op & 0x0F00) >> 8]] == 1 {
				vm.pc = vm.pc + 4
			} else {
				vm.pc = vm.pc + 2
			}
		//0x00A1 here
		default:
			fmt.Printf("Invalid opcode 0x%x\n", vm.op)
			vm.pc = vm.pc + 2
		}	

	case 0xF000:
		switch vm.op * 0x00FF {
		case 0x0018: //0xFX18 sets sound timer to X
			vm.soundTime = vm.v[(vm.op & 0x0F00) >> 8]
			vm.pc = vm.pc + 2
		default:
			fmt.Printf("Invalid opcode 0x%x\n", vm.op)
			vm.pc = vm.pc + 2
		}

	default:
		fmt.Printf("Invalid opcode 0x%X\n", vm.op)
		vm.pc = vm.pc + 2
	}
	
}

func (vm *chip8) LoadProgram(filePath string) error {
	//Reads file using os library
	file, fileErr := os.OpenFile(filePath, os.O_RDONLY, 0777)
	if fileErr != nil {
		return fileErr
	}
	defer file.Close()

	//Reads file information, also using os
	fStat, fStatErr := file.Stat()
	if fStatErr != nil {
		return fStatErr
	}
	if int64(len(vm.mem)-512) < fStat.Size() { //Checks file size doesn't exceed memory space
		return fmt.Errorf("File size is greater than memory")
	}

	buffer := make([]byte, fStat.Size())
	//Checks file can be read properly
	if _, readErr := file.Read(buffer); readErr != nil {
		return readErr
	}

	//If there are no errors, replace every byte in memory from 0x200 onward
	for i := 0; i < len(buffer); i++ {
		vm.mem[i+512] = buffer[i]
	}

	return nil
}





