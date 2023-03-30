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
	gui "alex/chip8/gui"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	sdl "github.com/veandco/go-sdl2/sdl"
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

	//Channel to check for audio events
	audioChan chan struct{}

	//CPU clock
	Clock *time.Ticker

	//Graphics
	gfx [64 * 32]byte //size of display

	//Draw flag
	drawFlag bool

	//Keyboard
	key [16]uint8 //CHIP-8 keypad had 16 keys

	//Channel to check for shutdown signal
	ShutdownChan chan struct{}

	//SDL window
	win *gui.Window

}

//Initialise emulator instance
func NewVM(filePath string, clockSpeed int)  (*chip8, error) {
	win, err := gui.NewWindow()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	vm := chip8{
		mem:			[4096]byte{},
		v:				[16]byte{},
		pc:				0x200,
		stack: 			[16]uint16{},
		gfx: 			[64 * 32]byte{},
		audioChan:		make(chan struct{}),
		Clock:			time.NewTicker(time.Second / time.Duration(clockSpeed)),
		ShutdownChan:	make(chan struct{}),
		win:			win,
		}

	//Load fontset
	for i := 0; i < 80; i++ {
		vm.mem[i] = fontSet[i]
		}


	if loadErr := vm.LoadProgram(filePath); loadErr != nil {
		panic(loadErr)
		return nil, loadErr
	}

	return &vm, nil
}

func (vm *chip8) Run() {
	for {
		select {
		case <-vm.Clock.C:
			if !vm.win.Closed() {
				if vm.delayTime > 0 {
					vm.delayTimeTick()
				} else {
					vm.FDE()
					vm.drawOrUpdate()
					vm.KeyPoll()
					vm.delayTimeTick()
					vm.soundTimeTick()
				}
				continue
			}
			break
		case <-vm.ShutdownChan:
			break
		}
		break
	}
	vm.signalShutdown("\nShutting down...")
}

func (vm *chip8) drawOrUpdate() {
	if vm.drawFlag {
		vm.win.DrawGraphics(vm.graphicsBuffer())
	} else {
		vm.win.UpdateInput()
	}
}

func (vm *chip8) KeyPoll() {
	if sdlError := sdl.Init(sdl.INIT_EVERYTHING); sdlError != nil {
		panic(sdlError)
	}
	for poll := sdl.PollEvent(); poll != nil; poll = sdl.PollEvent() {
			switch pl := poll.(type) {
			case *sdl.QuitEvent:
				fmt.Printf("Quit event detected")		
			case *sdl.KeyboardEvent:
				if pl.Type == sdl.KEYUP {
					switch pl.Keysym.Sym {
					case sdl.K_1:
						vm.Key(0x1, false)
					case sdl.K_2:
						vm.Key(0x2, false)
					case sdl.K_3:
						vm.Key(0x3, false)
					case sdl.K_4:
						vm.Key(0xC, false)
					case sdl.K_q:
						vm.Key(0x4, false)
					case sdl.K_w:
						vm.Key(0x5, false)
					case sdl.K_e:
						vm.Key(0x6, false)
					case sdl.K_r:
						vm.Key(0xD, false)
					case sdl.K_a:
						vm.Key(0x7, false)
					case sdl.K_s:
						vm.Key(0x8, false)
					case sdl.K_d:
						vm.Key(0x9, false)
					case sdl.K_f:
						vm.Key(0xE, false)
					case sdl.K_z:
						vm.Key(0xA, false)
					case sdl.K_x:
						vm.Key(0x0, false)
					case sdl.K_c:
						vm.Key(0xB, false)
					case sdl.K_v:
						vm.Key(0xF, false)
					}
				} else if pl.Type == sdl.KEYDOWN {
					switch pl.Keysym.Sym {
					case sdl.K_1:
						vm.Key(0x1, true)
					case sdl.K_2:
						vm.Key(0x2, true)
					case sdl.K_3:
						vm.Key(0x3, true)
					case sdl.K_4:
						vm.Key(0xC, true)
					case sdl.K_q:
						vm.Key(0x4, true)
					case sdl.K_w:
						vm.Key(0x5, true)
					case sdl.K_e:
						vm.Key(0x6, true)
					case sdl.K_r:
						vm.Key(0xD, true)
					case sdl.K_a:
						vm.Key(0x7, true)
					case sdl.K_s:
						vm.Key(0x8, true)
					case sdl.K_d:
						vm.Key(0x9, true)
					case sdl.K_f:
						vm.Key(0xE, true)
					case sdl.K_z:
						vm.Key(0xA, true)
					case sdl.K_x:
						vm.Key(0x0, true)
					case sdl.K_c:
						vm.Key(0xB, true)
					case sdl.K_v:
						vm.Key(0xF, true)
					}
				}
			}
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

	fileBuffer := make([]byte, fStat.Size())
	//Checks file can be read properly
	if _, readErr := file.Read(fileBuffer); readErr != nil {
		return readErr
	}

	//If there are no errors, replace every byte in memory from 0x200 onward
	for i := 0; i < len(fileBuffer); i++ {
		vm.mem[i+512] = fileBuffer[i]
	}

	return nil
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

func (vm *chip8) soundTimeTick() {
	if vm.soundTime > 0 {
		if vm.soundTime == 1 {
			vm.audioChan <- struct{}{}
		}
		vm.soundTime--
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
		fmt.Printf("\nBeep!!")
	}
}	

func (vm chip8) graphicsBuffer() [64 * 32]byte {
	return vm.gfx
}	

func (vm *chip8) signalShutdown(msg string) {
	fmt.Println(msg)
	close(vm.audioChan)
	vm.ShutdownChan <- struct{}{}
}

func (vm *chip8) drawSprite(x, y uint16) {
	height := vm.op & 0x000F
	vm.v[0xF] = 0
	var pix uint16

	for yLine := uint16(0); yLine < height; yLine++ {
		pix = uint16(vm.mem[vm.I+yLine])

		for xLine := uint16(0); xLine < 8; xLine++ {
			ind := (x + xLine + ((y + yLine) * 64))
			if ind >= uint16(len(vm.graphicsBuffer())) {
				continue
			}
			if (pix & (0x80 >> xLine)) != 0 {
				if vm.graphicsBuffer()[ind] == 1 {
					vm.v[0xF] = 1
				}
				vm.gfx[ind] ^= 1
			}
		}
	}

	vm.drawFlag = true
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
			case 0x0000: //0x00E0 clears screen
				vm.gfx = [64 * 32]byte{}
				vm.pc += 2
			case 0x000E: //0x00EE returns from a subroutine
				vm.pc = vm.stack[vm.sp] + 2
				vm.sp--
			default: 
				fmt.Printf("Invalid opcode 0x%X\n", vm.op)
			}
		
	//First nibble is 0001
	case 0x1000: //0x1NNN jumps to NNN on memory
		//Program counter = last 3 nibbles of opcode
		vm.pc = vm.op & 0x0FFF

	//First nibble is 0002
	case 0x2000: //0x2nnn calls subroutine at nnn
		vm.sp += 1 
		vm.stack[vm.sp] = vm.pc 
		vm.pc = vm.op & 0x0FFF

	case 0x3000: //0x3XKK Skips next instruction if VX = KK
		if uint16(vm.v[(vm.op & 0x0F00) >> 8]) == uint16(vm.op & 0x00FF) {
			vm.pc += 4
		} else { //Otherwise does nothing
			vm.pc += 2
		}

	case 0x4000: //0x4XKK skips next instruction if VX != KK
		if uint16(vm.v[(vm.op & 0x0F00) >> 8]) != uint16(vm.op & 0x00FF) {
			vm.pc += 4
		} else { //Otherwise does nothing
			vm.pc += 2 
		}

	case 0x5000: //0x5XY0 skips next instruction if VX = VY
		if uint16(vm.v[(vm.op & 0x0F00) >> 8]) == uint16(vm.v[(vm.op & 0x00F0) >> 4]) {
			vm.pc += 4
		} else {
			vm.pc += 2
		}

	case 0x6000: //0x6XNN sets value of VX to NN
		vm.v[(vm.op & 0x0F00) >> 8] = byte(vm.op & 0x00FF)
		vm.pc += 2

	case 0x7000: //0x7XNN adds value NN to value VX, doesn't affect carry
		vm.v[(vm.op & 0x0F00) >> 8] += byte(vm.op & 0x00FF)
		vm.pc += 2

	case 0x8000: 
		switch vm.op & 0x000F { //Check the last nibble 
			
			case 0x0001: //0x8XY1 Sets VX to VX OR VY
				vm.v[(vm.op & 0x0F00) >> 8] = (vm.v[(vm.op & 0x0F00) >> 8] | vm.v[(vm.op & 0x00F0) >> 4])
				vm.pc += 2
			
			case 0x0002: //0x8XY2 Sets VX to VX AND VY
				vm.v[(vm.op & 0x0F00) >> 8] = (vm.v[(vm.op & 0x0F00) >> 8] & vm.v[(vm.op & 0x00F0) >> 4])
				vm.pc += 2

			case 0x0003: //0x8XY3 Sets VX to VX XOR (exclusive OR) VY
				vm.v[(vm.op & 0x0F00) >> 8] = (vm.v[(vm.op & 0x0F00) >> 8] ^ vm.v[(vm.op & 0x00F0) >> 4])
				vm.pc += 2

			case 0x0004: //0x8XY4 Adds VY to VX, sets VF to 1 if result overflows
				if vm.v[(vm.op & 0x00F0) >> 4] > (0xFF - vm.v[(vm.op & 0x0F00) >> 8]) {
					vm.v[15] = 1 //Set VF to 1
				} else {
					vm.v[15] = 0
				}
				vm.v[(vm.op & 0x0F00) >> 8] += vm.v[(vm.op & 0x00F0) >> 4] //Otherwise, add VX and VY
			
			case 0x0005: //0x8XY5 Sets VX to VX - VY, sets VF to 1 if a borrow occurs
				if vm.v[(vm.op & 0x00F0) >> 4] > vm.v[(vm.op & 0x0F00) >> 8] {
					vm.v[0xF] = 0
				} else {
					vm.v[0xF] = 1
				}
				vm.v[(vm.op & 0x0F00) >> 8] -= vm.v[(vm.op & 0x00F0) >> 4]
				vm.pc += 2

			//case 0x0006: //Sets VX to VY right shifted by 1, setting VF to the bit lost in the shift



			}

	case 0xA000: //0xANNN Sets I to address NNN
		vm.I = vm.op & 0x0FFF
		vm.pc += 2

	case 0xD000: //0xDXYN Draws sprite of length N in memory starting at I at co-ords (VX, VY)
		x := uint16(vm.v[(vm.op & 0x0F00) >> 8])
		y := uint16(vm.v[(vm.op & 0x00F0) >> 4])
		vm.drawSprite(x, y)
		vm.pc += 2

	case 0xE000:
		switch vm.op & 0x00FF {
			case 0x009E:
				if vm.key[vm.v[(vm.op & 0x0F00) >> 8]] == 1 {
					vm.pc += 4
				} else {
					vm.pc += 2
				}
		//0x00A1 here
			default:
				fmt.Printf("Invalid opcode 0x%x\n", vm.op)
				vm.pc += 2
		}	

	case 0xF000:
		switch vm.op & 0x00FF {
			case 0x0015: //0xFX15 sets delay timer to VX
				vm.delayTime = vm.v[(vm.op & 0x0F00) >> 8]
				vm.pc += 2
			case 0x0018: //0xFX18 sets sound timer to VX
				vm.soundTime = vm.v[(vm.op & 0x0F00) >> 8]
				vm.pc += 2
			default:
				fmt.Printf("Invalid opcode 0x%x\n", vm.op)
				vm.pc += 2
			}

	default:
		fmt.Printf("Invalid opcode 0x%X\n", vm.op)
		vm.pc += 2
	}
}
	







