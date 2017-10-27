package game

import "time"

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
