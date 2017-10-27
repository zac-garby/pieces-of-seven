package entity

import (
	"math"

	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/world"
	"github.com/veandco/go-sdl2/sdl"
)

// ShipSpeed is the speed of the ship's apparent movement,
// in pixels per second.
const ShipSpeed = 20

// A Ship (an Entity,) is a player's ship.
type Ship struct {
	Pos         geom.Coord  // The actual position
	Name        string      // The ship's name
	ApparentPos geom.Vector // The interpolated position, for rendering

	// 0 1 2
	// 7 x 3
	// 6 5 4
	direction uint8

	Path Path
}

// NewShip initialises a new Ship at the
// given coordinates.
func NewShip(x, y uint) *Ship {
	return &Ship{
		Pos:         geom.Coord{X: x, Y: y},
		ApparentPos: geom.Vector{X: float64(x), Y: float64(y)},
	}
}

// Move sets the ship's current movement path
// to one going to (x, y).
func (s *Ship) Move(to geom.Coord, in *world.World) {
	coords, found := in.Graph.FindPath(s.Pos, to)
	if !found {
		return
	}

	path := make(Path, len(coords))
	for i, coord := range coords {
		path[i] = coord
	}

	s.Path = path
}

// Render renders the ship on the given renderer.
func (s *Ship) Render(viewOffset *geom.Vector, ld *loader.Loader, rend *sdl.Renderer) {
	rend.SetDrawColor(255, 0, 0, 255)

	tex := ld.Textures["ship"]

	src := s.getSheetRect()

	dst := &sdl.Rect{
		X: int32(s.ApparentPos.X*world.TileSize) - int32(viewOffset.X),
		Y: int32(s.ApparentPos.Y*world.TileSize) - int32(viewOffset.Y),
		W: int32(world.TileSize),
		H: int32(world.TileSize),
	}

	rend.Copy(tex, src, dst)
}

// Step is called every tick
func (s *Ship) Step() {

}

// Update moves the ship's apparent position towards its destination.
func (s *Ship) Update(dt float64) {
	if len(s.Path) == 0 {
		return
	}

	NeXTStep := s.Path[0]

	diff := geom.Vector{
		X: float64(NeXTStep.X) - float64(s.ApparentPos.X),
		Y: float64(NeXTStep.Y) - float64(s.ApparentPos.Y),
	}

	dist := math.Sqrt(diff.X*diff.X + diff.Y*diff.Y)

	if dist < 0.1 {
		s.Path = s.Path[1:]

		s.Pos = NeXTStep

		s.ApparentPos.X = float64(s.Pos.X)
		s.ApparentPos.Y = float64(s.Pos.Y)
	} else {
		movement := geom.Vector{
			X: diff.X / dist / ShipSpeed,
			Y: diff.Y / dist / ShipSpeed,
		}

		s.ApparentPos.X += movement.X
		s.ApparentPos.Y += movement.Y
	}

	d := geom.Vector{
		X: float64(s.ApparentPos.X) - float64(s.Pos.X),
		Y: float64(s.ApparentPos.Y) - float64(s.Pos.Y),
	}

	s.direction = s.diffToDirection(d)
}

func (s *Ship) getSheetRect() *sdl.Rect {
	var (
		x = 0
		y = 0
	)

	switch s.direction {
	case 0:
		x = 3
		y = 1
	case 1:
		break
	case 2:
		y = 1
	case 3:
		x = 1
	case 4:
		x = 1
		y = 1
	case 5:
		x = 2
	case 6:
		x = 2
		y = 1
	case 7:
		x = 3
	}

	return &sdl.Rect{
		X: int32(x * 15),
		Y: int32(y * 15),
		W: 15,
		H: 15,
	}
}

func (s *Ship) diffToDirection(diff geom.Vector) uint8 {
	if diff.X > 0 {
		if diff.Y > 0 {
			return 4
		} else if diff.Y < 0 {
			return 2
		} else {
			return 3
		}
	} else if diff.X < 0 {
		if diff.Y > 0 {
			return 6
		} else if diff.Y < 0 {
			return 0
		} else {
			return 7
		}
	} else {
		if diff.Y > 0 {
			return 5
		} else if diff.Y < 0 {
			return 1
		}
	}

	// Default to current direction
	return s.direction
}
