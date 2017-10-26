package ui

import (
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/veandco/go-sdl2/sdl"
)

type idComponent struct {
	Component
	ID string
}

// An Interface manages a set of UI
// Components.
type Interface struct {
	components []idComponent

	// The ID of the active component.
	// "" means no component is active.
	active string

	Padding uint
}

// Get returns the component whose ID == id,
// and nil if there isn't such a component.
func (i *Interface) Get(id string) (Component, int) {
	for i, comp := range i.components {
		if comp.ID == id {
			return comp.Component, i
		}
	}

	return nil, -1
}

// Add appends a component to the Interface.
// If a component already exists with the same
// ID, it's replaced.
func (i *Interface) Add(id string, comp Component) {
	idcomp := idComponent{
		Component: comp,
		ID:        id,
	}

	if old, idx := i.Get(id); old != nil {
		i.components[idx] = idcomp
	}

	i.components = append(i.components, idcomp)
}

// GetActive returns the active component,
// or nil if no component is active.
func (i *Interface) GetActive() Component {
	comp, _ := i.Get(i.active)
	return comp
}

// Layout recalculates the layout of the
// components, relative to (x, y).
func (i *Interface) Layout(x, y, width, height uint) {
	nextCoord := geom.Coord{X: x, Y: y}

	for _, comp := range i.components {
		comp.SetRect(&sdl.Rect{
			X: int32(nextCoord.X),
			Y: int32(nextCoord.Y),
			W: int32(width),
			H: int32(height),
		})

		nextCoord.Y += height + i.Padding
	}
}

func (i *Interface) Render(rend *sdl.Renderer) {
	for _, comp := range i.components {
		comp.Render(rend)
	}
}

func (i *Interface) Update(dt float64) {
	xi, yi, _ := sdl.GetMouseState()

	var (
		x = int32(xi)
		y = int32(yi)
	)

	cursor := sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_ARROW)

	for _, comp := range i.components {
		comp.Update(dt)

		rect := comp.GetRect()

		if x >= rect.X && y >= rect.Y && x < rect.X+rect.W && y < rect.Y+rect.H {
			cursor = sdl.CreateSystemCursor(comp.Cursor())
		}
	}

	sdl.SetCursor(cursor)
}

func (i *Interface) HandleEvent(event sdl.Event) {
	// If the left mouse button was pressed
	if evt, ok := event.(*sdl.MouseButtonEvent); ok &&
		evt.Type == sdl.MOUSEBUTTONDOWN && evt.Button == sdl.BUTTON_LEFT {

		for _, comp := range i.components {
			rect := comp.GetRect()

			// If the component was clicked on
			if evt.X >= rect.X && evt.Y >= rect.Y && evt.X < rect.X+rect.W && evt.Y < rect.Y+rect.H {
				if old := i.GetActive(); old != nil {
					old.Deactivate()
				}

				i.active = comp.ID

				comp.Activate()

				break
			}
		}
	}

	if len(i.active) > 0 {
		if comp := i.GetActive(); comp != nil {
			comp.HandleEvent(event)
		}
	}
}
