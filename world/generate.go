package world

import (
	"math/rand"
	"time"
)

// Some constants related to world generation
const (
	genThreshold  = 0.5
	genSampleDist = 3
	genIterations = 5
)

type genData [Height][Width]int
type genCoord [2]int

func init() {
	// Randomly seed the rng
	rand.Seed(int64(time.Now().Nanosecond()))
}

// Generate creates a world with randomly generated
// terrain. The way it works is first filling the
// world with random ints, either 1 or 0. Then, it
// will go through each cell and find the average
// value of all cells in a certain radius. If that
// average is above the threshold the cell is set
// to 1, and otherwise it's set to 0. This is repeated
// a number of times.
func Generate() *World {
	var data [Height][Width]int

	// Random initial data
	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			data[y][x] = rand.Int() % 2
		}
	}

	// Iterate the data
	for i := 0; i < genIterations; i++ {
		data = iterate(data)
	}

	w := New()

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			if data[y][x] > 0 {
				w.Tiles[y][x] = Land
			}
		}
	}

	return w
}

func iterate(data genData) genData {
	var newData genData

	for y := 0; y < Height; y++ {
		for x := 0; x < Width; x++ {
			coords := coordsInCircle(genCoord{x, y})
			total := 0.0

			for _, coord := range coords {
				var (
					cx = coord[0]
					cy = coord[1]
				)

				if cx >= 0 && cy >= 0 && cx < Width && cy < Height {
					total += float64(data[cy][cx])
				} else {
					// Prevent tiles being added around the edges of the map
					total -= 10
				}
			}

			avg := total / float64(len(coords))

			if avg > genThreshold {
				newData[y][x] = 1
			} else {
				newData[y][x] = 0
			}
		}
	}

	return newData
}

func coordsInCircle(origin genCoord) []genCoord {
	var (
		ox     = origin[0]
		oy     = origin[1]
		coords []genCoord
	)

	for x := ox - genSampleDist; x <= ox+genSampleDist; x++ {
		for y := oy - genSampleDist; y <= oy+genSampleDist; y++ {
			dist := (x-ox)*(x-ox) + (y-oy)*(y-oy)

			if dist <= genSampleDist*genSampleDist {
				coords = append(coords, genCoord{x, y})
			}
		}
	}

	return coords
}
