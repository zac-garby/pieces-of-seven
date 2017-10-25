package entity

import (
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/world"
	"github.com/veandco/go-sdl2/sdl"
)

// An Entity is a movable thing which is updated
// every tick and can be rendered.
type Entity interface {
	Move(to geom.Coord, in *world.World)
	Render(viewOffset *geom.Vector, ld *loader.Loader, rend *sdl.Renderer)

	Step()             // Called every tick
	Update(dt float64) // Called every frame
}
