package message

import (
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/Zac-Garby/pieces-of-seven/world"
	"github.com/satori/go.uuid"
)

// This file contains messages which the server
// can send to the client.

// An AbstractPlayer is a slightly compressed
// struct which can be expanded again to create
// an entity.Ship.
type AbstractPlayer struct {
	Position    geom.Coord
	Destination geom.Coord
}

// GameInfo is the information initially sent
// to a new client.
type GameInfo struct {
	Tiles   [world.Height][world.Width]world.Tile
	Players map[uuid.UUID]AbstractPlayer

	ID uuid.UUID // The UUID of the receiving client
}

// NewPlayer tells existing players that a
// new client has joined.
type NewPlayer struct {
	ID     uuid.UUID
	Player AbstractPlayer
}

// PlayerLeft tells players that a client
// has left the game.
type PlayerLeft struct {
	ID uuid.UUID
}

// PlayerMoved tells a client that a client
// has moved, and where he moved to.
type PlayerMoved struct {
	ID       uuid.UUID
	Position geom.Coord
}
