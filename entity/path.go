package entity

import "github.com/Zac-Garby/pieces-of-seven/geom"

// A Path stores the movement path of an entity
// as a slice of coordinates.
type Path []geom.Coord

// Normalise improves a path by cutting of corners, such that:
//
//      o o o          o o
//          o    ->        o
//          o              o
//
// This means that the entity following the path will skip the
// corner and go diagonally.
//
func (p *Path) Normalise() *Path {
	rewrite := make(Path, 0, len(*p))

	for i, coord := range *p {
		if i >= len(*p)-2 {
			rewrite = append(rewrite, coord)

			continue
		}

		// The next 2 coordinates in the path
		next := (*p)[i+1]
		prev := (*p)[i+2]

		// If coord, next, and prev are in a straight line, add coord
		// to the rewritten path.
		if (coord.X == next.X && coord.X == prev.X) || (coord.Y == next.Y && coord.Y == prev.Y) {
			rewrite = append(rewrite, coord)
		}
	}

	rewrite = rewrite[1:]
	return &rewrite
}
