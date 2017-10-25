package world

import (
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/veandco/go-sdl2/sdl"
)

// Width is the width, in Tiles, of the World.
const Width = 32

// Height is the height, in Tiles, of the World.
const Height = 32

// The World is a 2d slice of Tiles.
// Coordinate (x, y) is at index [y-1][x-1].
// It also stores a path-finding graph.
type World struct {
	Tiles [Height][Width]Tile
	*Graph
}

// New creates a new World instance.
func New() *World {
	world := &World{}

	world.Tiles[5][5] = 1
	world.Tiles[6][6] = 1
	world.Tiles[7][7] = 1
	world.Tiles[8][8] = 1
	world.Tiles[9][9] = 1
	world.Tiles[10][10] = 1

	world.MakeGraph()

	return world
}

// Render renders the world to the given
// SDL renderer.
func (w *World) Render(rend *sdl.Renderer, viewOffset *geom.Vector, width, height int) {
	for tile, data := range tileData {
		rend.SetDrawColor(data.Colour[0], data.Colour[1], data.Colour[2], data.Colour[3])
		rects := w.getRectsOfType(tile, viewOffset, width, height)

		if len(rects) > 0 {
			rend.FillRects(rects)
		}
	}
}

// Tick steps the world by one tick.
func (w *World) Tick() {

}

func (w *World) getRectsOfType(t Tile, viewOffset *geom.Vector, width, height int) []sdl.Rect {
	rects := []sdl.Rect{}

	// Calculate the amount of visible tiles, with some
	// padding on the side just in case
	tilesWide := width/TileSize + 2
	tilesHigh := height/TileSize + 2
	startX := int(viewOffset.X)/TileSize - 1
	startY := int(viewOffset.Y)/TileSize - 1

	for y := startY; y < startY+tilesHigh; y++ {
		for x := startX; x < startX+tilesWide; x++ {
			if y < Height && x < Width && y >= 0 && x >= 0 && w.Tiles[y][x] == t {
				rects = append(rects, sdl.Rect{
					X: int32(x*TileSize) - int32(viewOffset.X),
					Y: int32(y*TileSize) - int32(viewOffset.Y),
					W: TileSize,
					H: TileSize,
				})
			}
		}
	}

	return rects
}

// MakeGraph creates a path-finding graph from the World.
func (w *World) MakeGraph() {
	w.Graph = &Graph{}

	for y := Height - 1; y >= 0; y-- {
		for x := 0; x < Width; x++ {
			node := &Node{
				Graph: w.Graph,
				Pos:   geom.Coord{X: uint(x), Y: uint(y)},
				Tile:  w.Tiles[y][x],
			}

			w.Graph[y][x] = node
		}
	}
}
