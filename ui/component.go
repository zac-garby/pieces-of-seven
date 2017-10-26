package ui

import (
	"github.com/veandco/go-sdl2/sdl"
)

// A Component is the interface for all UI
// elements, such as buttons.
type Component interface {
	SetRect(*sdl.Rect)
	GetRect() *sdl.Rect

	Render(rend *sdl.Renderer)
	Update(dt float64, mx, my uint)
}
