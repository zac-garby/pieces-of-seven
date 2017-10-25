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

	Padding uint
}

// Get returns the component who's ID == id,
// and nil if there isn't such a component.
func (i *Interface) Get(id string) (*Component, int) {
	for i, comp := range i.components {
		if comp.ID == id {
			return &comp.Component, i
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

func (i *Interface) Update(dt float64, mx, my uint) {
	for _, comp := range i.components {
		comp.Update(dt, mx, my)
	}
}
