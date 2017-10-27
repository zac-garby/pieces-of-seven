package main

import (
	"time"

	"github.com/Zac-Garby/pieces-of-seven/scene/mainmenu"

	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/scene"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	width  = 1200
	height = 800
)

var (
	scn scene.Scene
	ld  loader.Loader
)

func main() {
	sdl.Init(sdl.INIT_EVERYTHING)
	defer sdl.Quit()

	ttf.Init()
	defer ttf.Quit()

	window, renderer, _ := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN|sdl.RENDERER_ACCELERATED)
	defer window.Destroy()
	defer renderer.Destroy()

	ld = loader.New()

	ld.Queue(map[string]loader.Asset{
		// Textures
		"ship":  {Path: "assets/sprites/ship.png", Type: loader.Texture, Data: make(map[string]int)},
		"icon":  {Path: "assets/icon.png", Type: loader.Texture, Data: make(map[string]int)},
		"water": {Path: "assets/tiles/water.png", Type: loader.Texture, Data: make(map[string]int)},
		"sand":  {Path: "assets/tiles/sand.png", Type: loader.Texture, Data: make(map[string]int)},

		// Fonts
		"body":    {Path: "assets/fonts/Montserrat-Regular.ttf", Type: loader.Font, Data: map[string]int{"size": 25}},
		"body-sm": {Path: "assets/fonts/Montserrat-Regular.ttf", Type: loader.Font, Data: map[string]int{"size": 20}},
	})

	ld.Load(renderer)
	defer ld.Free()

	// Set the icon to 'assets/icon.png'
	window.SetIcon(ld.Surfaces["icon"])

	// Set the window title to Game
	window.SetTitle("Game")

	// Ensure the window has focus
	window.Raise()

	// Enable VSync
	sdl.GL_SetSwapInterval(1)

	// Initialise the game scene
	scn = mainmenu.New(&ld)
	scn.Enter()
	defer scn.Exit()

	last := time.Now()

	running := true
main:
	for running {
		dt := time.Since(last).Seconds()
		last = time.Now()

		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false

			default:
				next := scn.HandleEvent(event)
				if next != "" {
					scn.Exit()
					scn = makeScene(next, &ld)
					scn.Enter()

					continue main
				}
			}
		}

		next := scn.Update(dt)
		if next != "" {
			scn.Exit()
			scn = makeScene(next, &ld)
			scn.Enter()

			continue main
		}

		renderer.SetDrawColor(0, 0, 0, 255)
		renderer.Clear()

		scn.Render(renderer, width, height)

		renderer.Present()
	}
}
