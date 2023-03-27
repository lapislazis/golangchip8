package main

import (
	chip8 "alex/CSProject/chip8"
	"time"
	sdl "github.com/veandco/go-sdl2/sdl"
	"os"
	"fmt"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		panic("Please provide a file/path/to/load and optionally a clockspeed (cycles per second)")
	}

	filePath := os.Args[1]
	clockSpeed := int32(700) 

	if len(os.Args) == 3 {
		if inp, inpErr := strconv.ParseInt(os.Args[2], 10, 32); inpErr != nil {
			panic(inpErr)
		} else {
			if inp > 0 {
				clockSpeed = int32(inp)
			}
		}
	}

	//Initialise emulator 
	chip8.Init()
	c8 := chip8.NewVM()
	if loadErr := c8.LoadProgram(filePath); loadErr != nil {
		panic(loadErr)
	}
	fmt.Printf("\nFile path is: %v\n", filePath)
	fmt.Printf("\nClockspeed is: %v\n", clockSpeed)

	//Initialise SDL2 (window and keyboard library)
	if sdlError := sdl.Init(sdl.INIT_EVERYTHING); sdlError != nil {
		panic(sdlError)
	}
	defer sdl.Quit()

	//Make window with CHIP-8 resolution of 64*32
	win, winError := sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1280, 640, sdl.WINDOW_SHOWN)
	if winError != nil {
		panic(winError)
	}
	defer win.Destroy()

	//Make renderer
	render, renderError := sdl.CreateRenderer(win, -1, 0)
	if renderError != nil {
		panic(renderError)
	}
	defer render.Destroy()

	for range chip8.Clock(time.Second/time.Duration(clockSpeed)) {
		c8.FDE()
		if c8.Draw() {
			render.SetDrawColor(0, 0, 0, 255)
			render.Clear()
			//Gpl gfx buffer
			buffer := c8.Buffer()
			for j := 0; j < len(buffer); j++ {
				for i := 0; i < len(buffer[j]); i++ {
					if buffer[j][i] != 0 { //If the pixel isn't empty
						render.SetDrawColor(255, 255, 255, 255) //Draw it as white
					} else { //If it is empty
						render.SetDrawColor(0, 0, 0, 255) //Draw it as black (red for testing)
					}
					render.FillRect(&sdl.Rect{ //Draws the display * 20 to fit the expanded resolution
						Y: int32(j) * 20,
						X: int32(i) * 20,
						W: 20,
						H: 20,
					})
				}
			}
			render.Present()
		}
		//The code below is for polling keyboard events. CHIP-8 keypads were very different from modern keyboards,
		//and so the key bindings have been changed from a normal keyboard layout to match. 
		//The format is the following:
		//Chip8 keypad         Keyboard mapping
		//1 | 2 | 3 | C        1 | 2 | 3 | 4
		//4 | 5 | 6 | D   =>   Q | W | E | R
		//7 | 8 | 9 | E   =>   A | S | D | F
		//A | 0 | B | F        Z | X | C | V
		//Check for keyboard events
		for poll := sdl.PollEvent(); poll != nil; poll = sdl.PollEvent() {
			switch pl := poll.(type) {
			case *sdl.QuitEvent:
				fmt.Printf("\nProgram quitting...\n")
				os.Exit(0)
			case *sdl.KeyboardEvent:
				if pl.Type == sdl.KEYUP {
					switch pl.Keysym.Sym {
					case sdl.K_1:
						c8.Key(0x1, false)
					case sdl.K_2:
						c8.Key(0x2, false)
					case sdl.K_3:
						c8.Key(0x3, false)
					case sdl.K_4:
						c8.Key(0xC, false)
					case sdl.K_q:
						c8.Key(0x4, false)
					case sdl.K_w:
						c8.Key(0x5, false)
					case sdl.K_e:
						c8.Key(0x6, false)
					case sdl.K_r:
						c8.Key(0xD, false)
					case sdl.K_a:
						c8.Key(0x7, false)
					case sdl.K_s:
						c8.Key(0x8, false)
					case sdl.K_d:
						c8.Key(0x9, false)
					case sdl.K_f:
						c8.Key(0xE, false)
					case sdl.K_z:
						c8.Key(0xA, false)
					case sdl.K_x:
						c8.Key(0x0, false)
					case sdl.K_c:
						c8.Key(0xB, false)
					case sdl.K_v:
						c8.Key(0xF, false)
					}
				} else if pl.Type == sdl.KEYDOWN {
					switch pl.Keysym.Sym {
					case sdl.K_1:
						c8.Key(0x1, true)
					case sdl.K_2:
						c8.Key(0x2, true)
					case sdl.K_3:
						c8.Key(0x3, true)
					case sdl.K_4:
						c8.Key(0xC, true)
					case sdl.K_q:
						c8.Key(0x4, true)
					case sdl.K_w:
						c8.Key(0x5, true)
					case sdl.K_e:
						c8.Key(0x6, true)
					case sdl.K_r:
						c8.Key(0xD, true)
					case sdl.K_a:
						c8.Key(0x7, true)
					case sdl.K_s:
						c8.Key(0x8, true)
					case sdl.K_d:
						c8.Key(0x9, true)
					case sdl.K_f:
						c8.Key(0xE, true)
					case sdl.K_z:
						c8.Key(0xA, true)
					case sdl.K_x:
						c8.Key(0x0, true)
					case sdl.K_c:
						c8.Key(0xB, true)
					case sdl.K_v:
						c8.Key(0xF, true)
					}
				}
			}
		}
	}

}	

