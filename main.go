package main

import (
	chip8 "alex/CSProject/chip8"
	"time"
	sdl "github.com/veandco/go-sdl2/sdl"
)

func main() {
	chip8.Init()
	c8 := chip8.NewVM()

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

	for range chip8.Clock(time.Second/700) {
		c8.FDE()
	}
}

