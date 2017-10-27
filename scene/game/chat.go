package game

import (
	"time"

	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/veandco/go-sdl2/sdl"
)

type MessageType int

const (
	NoneMessage   MessageType = 0
	GlobalMessage             = 1 << (iota - 1)
	PrivateMessage
	ServerMessage
	DebugMessage
)

// By default, players see global, private, and server messages
const DefaultMessageMask = GlobalMessage | PrivateMessage | ServerMessage

// A Message contains the sender, text content,
// and time sent, of a chat message.
type Message struct {
	Sender  string
	Content string
	Type    MessageType
	Time    time.Time
}

func (m Message) IsVisible(mask MessageType) bool {
	return m.Type&mask > 0
}

// A ChatLog contains all the incoming messages,
// as well as
type ChatLog struct {
	Messages []*Message
	Mask     MessageType
}

// NewChatLog creates a new ChatLog
func NewChatLog() *ChatLog {
	return &ChatLog{
		Messages: []*Message{
			{
				Sender:  "server",
				Content: "hello, world",
				Type:    ServerMessage,
				Time:    time.Now(),
			},

			{
				Sender:  "server",
				Content: "hello, world. hello, world. hello, world. hello, world. testing wrapping. hello, world. hello, world. hello, world.",
				Type:    ServerMessage,
				Time:    time.Now(),
			},
		},

		Mask: DefaultMessageMask,
	}
}

// Log adds a new message to a chat log
func (c *ChatLog) Log(msg *Message) {
	c.Messages = append(c.Messages, msg)
}

// GetVisible returns a slice of all the messages
// visible under the current mask.
func (c *ChatLog) GetVisible() []*Message {
	var visible []*Message

	for _, msg := range c.Messages {
		if msg.IsVisible(c.Mask) {
			visible = append(visible, msg)
		}
	}

	return visible
}

// Render renders a chat log on an SDL renderer.
func (c *ChatLog) Render(rend *sdl.Renderer, ld *loader.Loader, x, y, width, height int) {
	bg := &sdl.Rect{
		X: int32(x),
		Y: int32(y),
		W: int32(width),
		H: int32(height),
	}

	rend.SetDrawColor(0, 0, 0, 255)
	rend.FillRect(bg)

	var (
		font    = ld.Fonts["body-sm"]
		msgs    = c.GetVisible()
		nextPos = geom.Coord{uint(x + 10), uint(y + 10)}
	)

	for _, msg := range msgs {
		solid, err := font.RenderUTF8_Blended_Wrapped(msg.Content, sdl.Color{R: 255, G: 255, B: 255, A: 255}, width-20)
		if err != nil {
			panic(err)
		}

		tex, err := rend.CreateTextureFromSurface(solid)
		if err != nil {
			panic(err)
		}

		var (
			src = &solid.ClipRect

			dest = &sdl.Rect{
				X: int32(nextPos.X),
				Y: int32(uint(height) - nextPos.Y - uint(src.H)),
				W: src.W,
				H: src.H,
			}

			separator = &sdl.Rect{
				X: int32(nextPos.X),
				Y: int32(uint(height) - nextPos.Y),
				W: int32(width - 20),
				H: 1,
			}
		)

		rend.Copy(tex, src, dest)

		rend.SetDrawColor(50, 50, 50, 255)
		rend.FillRect(separator)

		nextPos.Y += uint(src.H + 10)
	}
}
