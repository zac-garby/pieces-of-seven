package lib

import (
	"bufio"
	"fmt"
	"net"

	"github.com/Zac-Garby/pieces-of-seven/entity"
	"github.com/Zac-Garby/pieces-of-seven/message"
	"github.com/Zac-Garby/pieces-of-seven/world"
	"github.com/satori/go.uuid"
)

// EOT is the end of transmission character
const EOT byte = 4

type Server struct {
	World   *world.World
	Players map[uuid.UUID]*entity.Ship
	Address string

	closed      map[uuid.UUID]bool
	connections map[uuid.UUID]net.Conn
}

func New(addr string) *Server {
	s := &Server{
		World:       world.Generate(),
		Address:     addr,
		Players:     make(map[uuid.UUID]*entity.Ship),
		connections: make(map[uuid.UUID]net.Conn),
		closed:      make(map[uuid.UUID]bool),
	}

	s.World.MakeGraph()

	return s
}

func (s *Server) Listen() error {
	ln, err := net.Listen("tcp", s.Address)
	if err != nil {
		return err
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		go s.handleConnection(conn)
	}

	return nil
}

func (s *Server) handleConnection(conn net.Conn) {
	id := uuid.NewV4()
	s.connections[id] = conn

	pos := s.World.FindFreeSpace()

	// Create a new Ship for the connected player
	player := entity.NewShip(pos.X, pos.Y)
	s.Players[id] = player

	if err := s.Send(id, &message.GameInfo{
		Tiles:   s.World.Tiles,
		Players: s.abstractPlayers(),
		ID:      id,
	}); err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)

	for {
		bytes, err := reader.ReadBytes(EOT)

		// An error will most likely occur because the
		// connection was dropped, in which case the
		// read loop is ended.
		if err != nil {
			s.handleDisconnect(id)
			break
		}

		// Message received
		if len(bytes) > 0 {
			bytes = bytes[:len(bytes)-1]

			msg, err := message.Deserialize(bytes)
			if err != nil {
				fmt.Printf("deserializing: %s\n", err.Error())
			}

			s.handleMessage(id, msg)
		}
	}
}

func (s *Server) abstractPlayers() map[uuid.UUID]message.AbstractPlayer {
	ap := make(map[uuid.UUID]message.AbstractPlayer)

	for id, ship := range s.Players {
		dest := ship.Pos

		if len(ship.Path) > 0 {
			dest = ship.Path[len(ship.Path)-1]
		}

		ap[id] = message.AbstractPlayer{
			Position:    ship.Pos,
			Destination: dest,
		}
	}

	return ap
}

func (s *Server) Send(id uuid.UUID, msg interface{}) error {
	b, err := message.Serialize(msg)
	if err != nil {
		return err
	}

	b = append(b, EOT)
	s.connections[id].Write(b)

	return nil
}

func (s *Server) Broadcast(msg interface{}) error {
	for id := range s.connections {
		if err := s.Send(id, msg); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) handleMessage(id uuid.UUID, msg interface{}) {
	switch m := msg.(type) {
	case *message.ClientInfo:
		pos := s.Players[id].Pos
		s.Players[id].Name = m.Name

		err := s.Broadcast(&message.NewPlayer{
			ID: id,
			Player: message.AbstractPlayer{
				Position:    pos,
				Destination: pos,
			},
		})

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("%s joined the game\n", m.Name)

	case *message.Disconnect:
		s.handleDisconnect(id)

	case *message.Moved:
		s.Players[id].Move(m.Position, s.World)

		err := s.Broadcast(&message.PlayerMoved{
			ID:       id,
			Position: m.Position,
		})

		if err != nil {
			fmt.Println(err)
		}

	case *message.StateUpdate:
		s.Players[id].Pos = m.Position
	}
}

func (s *Server) handleDisconnect(id uuid.UUID) {
	if s.closed[id] {
		return
	}

	s.closed[id] = true

	err := s.Broadcast(&message.PlayerLeft{
		ID: id,
	})

	if err != nil {
		fmt.Println("in handleDisconnect:", err)
	}

	fmt.Printf("%s left the game\n", s.Players[id].Name)

	// Close the connection and delete the player
	s.connections[id].Close()
	delete(s.Players, id)
}
