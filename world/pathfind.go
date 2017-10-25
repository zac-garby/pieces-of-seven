package world

import (
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/beefsack/go-astar"
)

// A Graph is a graph representation of a World,
// for use in pathfinding.
type Graph [Height][Width]*Node

// At returns the node located at (x, y)
func (g *Graph) At(x, y int) *Node {
	if x < 0 || x >= Width || y < 0 || y >= Height {
		return nil
	}

	return g[y][x]
}

// AtCoord returns the node located at the given coordinate
func (g *Graph) AtCoord(coord geom.Coord) *Node {
	return g.At(int(coord.X), int(coord.Y))
}

// FindPath finds a path from a coordinate, to another,
// and returns a slice of coordinates and a boolean saying
// whether a path exists or not.
func (g *Graph) FindPath(from, to geom.Coord) ([]geom.Coord, bool) {
	var (
		start = g.AtCoord(from)
		end   = g.AtCoord(to)
	)

	if start == nil || end == nil || !end.Tile.GetData().Passable {
		return []geom.Coord{}, false
	}

	path, _, found := astar.Path(end, start)
	if !found {
		return []geom.Coord{}, false
	}

	coords := []geom.Coord{}

	for _, node := range path {
		coords = append(coords, node.(*Node).Pos)
	}

	return coords, true
}

// A Node represents a Tile as a node on a graph,
// for use in pathfinding.
type Node struct {
	Tile  Tile
	Pos   geom.Coord
	Graph *Graph
}

// PathNeighbors returns the immediate neighbours
// of a Node.
func (n *Node) PathNeighbors() []astar.Pather {
	var (
		x          = int(n.Pos.X)
		y          = int(n.Pos.Y)
		neighbours = []astar.Pather{}
	)

	for _, offset := range [][]int{
		{0, -1},
		{0, +1},
		{-1, 0},
		{+1, 0},
	} {
		if node := n.Graph.At(x+offset[1], y+offset[0]); node != nil && node.Tile.GetData().Passable {
			neighbours = append(neighbours, node)
		}
	}

	return neighbours
}

// PathNeighborCost returns the cost to move from one
// node to another. At the moment, the cost is always 1.0.
func (n *Node) PathNeighborCost(to astar.Pather) float64 {
	return 1.0
}

// PathEstimatedCost returns the manhattan distance
// between 2 nodes, |x2 - x1| + |y2 - y1|
func (n *Node) PathEstimatedCost(to astar.Pather) float64 {
	other := to.(*Node)

	absX := int(other.Pos.X - n.Pos.X)
	if absX < 0 {
		absX = -absX
	}

	absY := int(other.Pos.Y - n.Pos.Y)
	if absY < 0 {
		absY = -absY
	}

	return float64(absX + absY)
}
