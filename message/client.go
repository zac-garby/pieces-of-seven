package message

import (
	"time"

	"github.com/Zac-Garby/pieces-of-seven/geom"
)

// This file contains messages sent from
// the client to the server.

// ClientInfo tells the server information
// about the client.
type ClientInfo struct {
	Name string
}

// A Disconnect message tells the server
// that the client has left the game.
type Disconnect struct{}

// Moved tells the server that the client
// has moved, and where he moved to.
type Moved struct {
	Position geom.Coord
}

// A StateUpdate is sent peridocally from
// the client to the server to inform it
// of state changes, to ensure it's synced
// properly.
type StateUpdate struct {
	Position geom.Coord
}

// A ChatMessage tells the server that the
// client has sent a message.
type ChatMessage struct {
	Sender  string
	Content string
	Time    time.Time
}
