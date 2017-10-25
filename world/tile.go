package world

// TileSize is the size of each Tile, in pixels.
const TileSize = 48

// A Tile represents the id of a single tile in a World.
type Tile int

// TileData stores information about each type of Tile.
type TileData struct {
	Name     string
	Passable bool
	Colour   Colour
}

const (
	// Water is the tile which the sea is made from.
	Water = iota

	// Land makes up islands.
	Land
)

// Colour is an array of 4 ints, representing the 4
// colour channels.
type Colour [4]uint8

var tileData = map[Tile]*TileData{
	0: &TileData{
		Name:     "water",
		Passable: true,
		Colour:   Colour{102, 156, 232, 255},
	},

	1: &TileData{
		Name:     "land",
		Passable: false,
		Colour:   Colour{232, 216, 102, 255},
	},
}

// GetData returns the data struct associated with the
// given Tile.
func (t Tile) GetData() *TileData {
	return tileData[t]
}
