package game

import (
	"bufio"
	"errors"
	"net"

	"fmt"

	"github.com/Zac-Garby/pieces-of-seven/entity"
	"github.com/Zac-Garby/pieces-of-seven/message"
)

// EOT is the end of transmission character
const EOT byte = 4

var ErrNoConnection = errors.New("no connection established")

type Client struct {
	Address string
	Name    string
	Game    *Game

	conn net.Conn
}

func NewClient(addr string, game *Game, name string) *Client {
	c := &Client{
		Address: addr,
		Game:    game,
		Name:    name,
	}

	return c
}

func (c *Client) Listen() {
	conn, err := net.Dial("tcp", c.Address)
	if err != nil {
		c.Game.shouldQuit = true
		return
	}

	c.conn = conn

	c.SendClientInfo()

	reader := bufio.NewReader(conn)

	for {
		reply, err := reader.ReadBytes(EOT)

		// An error here will most likely be because
		// the connection to the server was dropped.
		if err != nil {
			c.LeaveGame()
			break
		}

		reply = reply[:len(reply)-1]

		msg, err := message.Deserialize(reply)
		if err != nil {
			fmt.Println(err)
			break
		}

		c.handleMessage(msg)
	}
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Send(msg interface{}) error {
	b, err := message.Serialize(msg)
	if err != nil {
		return err
	}

	b = append(b, EOT)

	if c.conn != nil {
		c.conn.Write(b)
	} else {
		return ErrNoConnection
	}

	return nil
}

func (c *Client) LeaveGame() error {
	return c.Send(&message.Disconnect{})
}

func (c *Client) handleMessage(msg interface{}) {
	switch m := msg.(type) {
	case *message.GameInfo:
		// Initialise the game's world with
		// the provided tiles.
		c.Game.World.Tiles = m.Tiles
		c.Game.World.MakeGraph()

		// Add the existing players to the game.
		for id, apl := range m.Players {
			ship := entity.NewShip(
				apl.Position.X,
				apl.Position.Y,
			)

			ship.Move(apl.Destination, c.Game.World)

			if id == m.ID {
				c.Game.Player = ship
			}

			c.Game.Entities = append(c.Game.Entities, ship)
			c.Game.Players[id] = ship
		}

	case *message.NewPlayer:
		ship := entity.NewShip(
			m.Player.Position.X,
			m.Player.Position.Y,
		)

		ship.Move(m.Player.Destination, c.Game.World)

		if _, exists := c.Game.Players[m.ID]; !exists {
			c.Game.Players[m.ID] = ship
			c.Game.Entities = append(c.Game.Entities, ship)
		}

	case *message.PlayerLeft:
		for i, ent := range c.Game.Entities {
			if ent == c.Game.Players[m.ID] {
				var start []entity.Entity
				if i > 0 {
					start = c.Game.Entities[:i-1]
				}

				var end []entity.Entity
				if i+1 < len(c.Game.Entities) {
					end = c.Game.Entities[i+1:]
				}

				c.Game.Entities = append(start, end...)

				break
			}
		}

		delete(c.Game.Players, m.ID)

	case *message.PlayerMoved:
		c.Game.Players[m.ID].Move(m.Position, c.Game.World)
	}
}

func (c *Client) SendClientInfo() error {
	info := &message.ClientInfo{
		Name: c.Name,
	}

	return c.Send(info)
}
