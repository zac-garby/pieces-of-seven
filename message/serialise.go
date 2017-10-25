package message

import (
	"encoding/json"
	"fmt"
	"reflect"
)

// Serialize serializes a message by marshalling
// it into JSON, then prepending a prefix unique
// to the type of message.
func Serialize(msg interface{}) ([]byte, error) {
	var prefix byte

	switch msg.(type) {
	case *GameInfo, GameInfo:
		prefix = 'g'
	case *PlayerMoved, PlayerMoved:
		prefix = 'p'
	case *NewPlayer, NewPlayer:
		prefix = 'n'
	case *PlayerLeft, PlayerLeft:
		prefix = 'l'

	case *ClientInfo, ClientInfo:
		prefix = 'c'
	case *Moved, Moved:
		prefix = 'm'
	case *Disconnect, Disconnect:
		prefix = 'x'

	default:
		return []byte{}, fmt.Errorf("invalid message type: %s", reflect.TypeOf(msg).String())
	}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return []byte{}, err
	}

	bytes = append([]byte{prefix}, bytes...)

	return bytes, nil
}

// Deserialize does the opposite of Serialize:
// takes some JSON data with an appropriate
// prefix, and returns the message.
func Deserialize(data []byte) (interface{}, error) {
	var template interface{}

	switch data[0] {
	case 'g':
		template = &GameInfo{}
	case 'p':
		template = &PlayerMoved{}
	case 'n':
		template = &NewPlayer{}
	case 'l':
		template = &PlayerLeft{}

	case 'c':
		template = &ClientInfo{}
	case 'm':
		template = &Moved{}
	case 'x':
		template = &Disconnect{}

	default:
		return nil, fmt.Errorf("invalid message prefix: %s", string(data[0]))
	}

	if err := json.Unmarshal(data[1:], template); err != nil {
		return nil, err
	}

	return template, nil
}
