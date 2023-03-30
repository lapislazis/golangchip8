package gui

import(
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"fmt"
)

//Instead of using SDL, which was not playing nice with Cobra, I 
//decided to use another library called pixel

const (
	winX			float64 = 64 
	winY			float64 = 32 
	screenWidth		float64 = 1024
	screenHeight	float64 = 768
)

type Window struct {
	*pixelgl.Window
}

func NewWindow() (*Window, error) {
	config := pixelgl.WindowConfig{
		Title: 		"CHIP-8",
		Bounds:		pixel.R(0, 0, 1024, 768),
		VSync:		true,
	}
	win, err := pixelgl.NewWindow(config)
	if err != nil {
		return nil, fmt.Errorf("Error creating new window: %v", err)
	}

	return &Window{win}, nil
}

func (win *Window) DrawGraphics(gfx ([64 * 32]uint8)) {
	win.Clear(colornames.Black)
	imDraw := imdraw.New(nil)
	imDraw.Color = pixel.RGB(1, 1, 1)
	w, h := float64(screenWidth/winX), float64(screenHeight/winY)

	for i := 0; i < 64; i++ {
		for j := 0; j < 32; j++ {
			// If the gfx byte in question is turned off,
			// continue and skip drawing the rectangle
			if gfx[(31-j)*64+i] == 0 {
				continue
			}
			imDraw.Push(pixel.V(w*float64(i), h*float64(j)))
			imDraw.Push(pixel.V(w*float64(i)+w, h*float64(j)+h))
			imDraw.Rectangle(0)
		}
	}

	imDraw.Draw(win)
	win.Update()
}

















//SDL was no worky :(
// type Render struct {
// 	*sdl.Window
// 	*sdl.Renderer
// }

// func NewWindow() (*Render, error) {
// 	//Initialise SDL2 (window and keyboard library)
// 	if sdlError := sdl.Init(sdl.INIT_EVERYTHING); sdlError != nil {
// 		panic(sdlError)
// 	}

// 	//Make window with CHIP-8 resolution of 64*32
// 	win, winError := sdl.CreateWindow("CHIP-8", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 1280, 640, sdl.WINDOW_SHOWN)
// 	if winError != nil {
// 		panic(winError)
// 	}

// 	//Make renderer
// 	render, renderError := sdl.CreateRenderer(win, -1, 0)
// 	if renderError != nil {
// 		panic(renderError)
// 	}

// 	return &Render{
// 		Window:		win,
// 		Renderer:	render,
// 	}, nil
// }

// func (render *Render) DrawOnScreen(gfx ([64 * 32]uint8)) {
// 	render.SetDrawColor(0, 0, 0, 255)
// 	render.Clear()
// 	for j := 0; j < len(gfx); j++ {
// 		for i := 0; i < len(gfx[j]); i++ {
// 			if gfx[j][i] != 0 { //If the pixel isn't empty
// 				render.SetDrawColor(255, 255, 255, 255) //Draw it as white
// 			} else { //If it is empty
// 				render.SetDrawColor(0, 0, 0, 255) //Draw it as black (red for testing)
// 			}
// 			render.FillRect(&sdl.Rect{ //Draws the display * 20 to fit the expanded resolution
// 				Y: int32(j) * 20,
// 				X: int32(i) * 20,
// 				W: 20,
// 				H: 20,
// 			})
// 		}
// 	}
// 	render.Present()
// }