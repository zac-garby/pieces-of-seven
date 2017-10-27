package world

import (
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/veandco/go-sdl2/sdl"
)

// Width is the width, in Tiles, of the World.
const Width = 256

// Height is the height, in Tiles, of the World.
const Height = 256

// The World is a 2d slice of Tiles.
// Coordinate (x, y) is at index [y-1][x-1].
// It also stores a path-finding graph.
type World struct {
	Tiles [Height][Width]Tile
	*Graph

	frame int32
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
	world.Tiles[5][6] = 1
	world.Tiles[5][7] = 1
	world.Tiles[6][8] = 1

	world.MakeGraph()

	return world
}

// Render renders the world to the given
// SDL renderer.
func (w *World) Render(rend *sdl.Renderer, ld *loader.Loader, viewOffset *geom.Vector, width, height int) {
	for _, tile := range Tiles {
		data := tile.GetData()
		tex := ld.Textures[data.Texture]
		frame := w.frame % data.Frames

		srcs, dsts := w.getRectsOfType(tile, viewOffset, width, height, frame)

		for i, rect := range dsts {
			rend.Copy(tex, &srcs[i], &rect)
		}
	}
}

func (w *World) getTexRectForMarchingSquares(x, y int) sdl.Rect {
	bits := 0

	// The if statements below cover the neighbours
	// in this order:
	//
	// 1 2 3
	// 4   5
	// 6 7 8

	matches := func(x, y int) bool {
		if x < 0 || y < 0 || x >= Width || y >= Height {
			return true
		}

		return w.Tiles[y][x] == Land
	}

	if w.Tiles[y][x] == Land {
		bits = 0xF
	}

	if matches(x-1, y-1) {
		bits |= 1 << 3
	}

	if matches(x, y-1) {
		bits |= 1<<3 | 1<<2
	}

	if matches(x+1, y-1) {
		bits |= 1 << 2
	}

	if matches(x-1, y) {
		bits |= 1<<3 | 1<<0
	}

	if matches(x+1, y) {
		bits |= 1<<2 | 1<<1
	}

	if matches(x-1, y+1) {
		bits |= 1 << 0
	}

	if matches(x, y+1) {
		bits |= 1<<0 | 1<<1
	}

	if matches(x+1, y+1) {
		bits |= 1 << 1
	}

	var texx, texy int32

	for i := 0; i < bits; i++ {
		texx++

		if texx > 3 {
			texx = 0
			texy++
		}
	}

	rect := sdl.Rect{
		X: texx * 15,
		Y: texy * 15,
		W: 15,
		H: 15,
	}

	return rect
}

// Tick steps the world by one tick.
func (w *World) Tick() {
	w.frame += 1
}

func (w *World) getRectsOfType(t Tile, viewOffset *geom.Vector, width, height int, frame int32) ([]sdl.Rect, []sdl.Rect) {
	dests := []sdl.Rect{}
	srcs := []sdl.Rect{}

	// Calculate the amount of visible tiles, with some
	// padding on the side just in case
	tilesWide := width/TileSize + 2
	tilesHigh := height/TileSize + 3
	startX := int(viewOffset.X)/TileSize - 1
	startY := int(viewOffset.Y)/TileSize - 1

	if t.GetData().MarchSquares {
		for y := startY; y < startY+tilesHigh; y++ {
			for x := startX; x < startX+tilesWide; x++ {
				if y < Height && x < Width && y >= 0 && x >= 0 {
					dests = append(dests, sdl.Rect{
						X: int32(x*TileSize) - int32(viewOffset.X),
						Y: int32(y*TileSize) - int32(viewOffset.Y),
						W: TileSize,
						H: TileSize,
					})

					srcs = append(srcs, w.getTexRectForMarchingSquares(x, y))
				}
			}
		}
	} else {
		texRect := sdl.Rect{
			X: 15 * frame,
			Y: 0,
			W: 15,
			H: 15,
		}

		for y := startY; y < startY+tilesHigh; y++ {
			for x := startX; x < startX+tilesWide; x++ {
				if y < Height && x < Width && y >= 0 && x >= 0 && (w.Tiles[y][x] == t || t == Water) {
					dests = append(dests, sdl.Rect{
						X: int32(x*TileSize) - int32(viewOffset.X),
						Y: int32(y*TileSize) - int32(viewOffset.Y),
						W: TileSize,
						H: TileSize,
					})

					srcs = append(srcs, texRect)
				}
			}
		}
	}

	return srcs, dests
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
