package scene

import "github.com/veandco/go-sdl2/sdl"

// A Scene - such as Game - is a particular state
// of the game.
type Scene interface {
	Enter()
	Exit()
	Update(dt float64) string
	Render(rend *sdl.Renderer, width, height int)
	HandleEvent(event sdl.Event) string
}
