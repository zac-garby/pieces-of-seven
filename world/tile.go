package world

// TileSize is the size of each Tile, in pixels.
const TileSize = 48

// A Tile represents the id of a single tile in a World.
type Tile int

// TileData stores information about each type of Tile.
type TileData struct {
	Name     string
	Passable bool
	Texture  string
	Frames   int32

	// If true, the renderer assumes the texture is 60x60,
	// with 16 tiles at 16x16 each.
	MarchSquares bool
}

const (
	// Water is the tile which the sea is made from.
	Water = iota

	// Land makes up islands.
	Land
)

// Tiles is a list of all rendered tiles, in their
// render order.
var Tiles = []Tile{
	Water,
	Land,
}

// Colour is an array of 4 ints, representing the 4
// colour channels.
type Colour [4]uint8

var tileData = map[Tile]*TileData{
	Water: {
		Name:         "water",
		Passable:     true,
		Texture:      "water",
		Frames:       1,
		MarchSquares: false,
	},

	Land: {
		Name:         "land",
		Passable:     false,
		Texture:      "sand",
		Frames:       1,
		MarchSquares: true,
	},
}

// GetData returns the data struct associated with the
// given Tile.
func (t Tile) GetData() *TileData {
	return tileData[t]
}
