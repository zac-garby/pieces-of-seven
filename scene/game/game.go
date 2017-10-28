package game

import (
	"time"

	"github.com/Zac-Garby/pieces-of-seven/entity"
	"github.com/Zac-Garby/pieces-of-seven/geom"
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/message"
	"github.com/Zac-Garby/pieces-of-seven/world"
	"github.com/satori/go.uuid"
	"github.com/veandco/go-sdl2/sdl"
)

// TickRate is the amount of ticks per second
const TickRate = 5

// ServerUpdateRate is the amount of times a client
// reports its state to the server every second
const ServerUpdateRate = 0.5

// ChatLogWidth is the width of the chat log,
// which is positioned to the side of the screen,
// in pixels.
const ChatLogWidth = 400

// A Game stores all the components of the game,
// to abstract the game data/logic from the main
// loop.
type Game struct {
	World      *world.World
	Entities   []entity.Entity
	Players    map[uuid.UUID]*entity.Ship
	Player     *entity.Ship
	ViewOffset *geom.Vector
	Client     *Client
	ChatLog    *ChatLog

	ld           *loader.Loader
	nextTick     float64
	nextUpdate   float64
	shouldQuit   bool
	shouldSetCam bool
}

// New creates a new Game instance.
func New(ld *loader.Loader, addr, name string) *Game {
	game := &Game{
		World:        &world.World{},
		ViewOffset:   &geom.Vector{X: 0, Y: 0},
		ChatLog:      NewChatLog(),
		nextTick:     1.0 / TickRate,
		nextUpdate:   1.0 / ServerUpdateRate,
		ld:           ld,
		shouldQuit:   false,
		shouldSetCam: true,
		Players:      make(map[uuid.UUID]*entity.Ship),
	}

	game.Client = NewClient(addr, game, name)

	return game
}

// Enter is called when a Game scene is entered.
func (g *Game) Enter() {
	go g.Client.Listen()

	sdl.StartTextInput()
}

// Exit is called when a Game scene is exited.
func (g *Game) Exit() {
	if !g.shouldQuit {
		g.Client.LeaveGame()
	}

	sdl.StopTextInput()
}

// Update updates the game by 'dt' seconds. The returned
// scene will be changed to in the main loop, or, if nil
// is returned, the scene won't be changed.
func (g *Game) Update(dt float64) string {
	if g.shouldQuit {
		g.Client.LeaveGame()
		return "mainmenu"
	}

	g.nextTick -= dt
	g.nextUpdate -= dt

	if g.nextTick <= 0 {
		g.nextTick = 1.0 / TickRate
		g.tick()
	}

	if g.nextUpdate <= 0 {
		g.nextUpdate = 1.0 / ServerUpdateRate
		g.serverUpdate()
	}

	for _, e := range g.Entities {
		e.Update(dt)
	}

	return ""
}

// Render renders a game (i.e. the objects inside it)
// onto an SDL renderer.
func (g *Game) Render(rend *sdl.Renderer, width, height int) {
	g.World.Render(rend, g.ld, g.ViewOffset, width-ChatLogWidth, height)

	for _, e := range g.Entities {
		e.Render(g.ViewOffset, g.ld, rend)
	}

	if g.Player != nil {
		ppos := geom.Vector{
			X: (g.Player.ApparentPos.X * world.TileSize) + (ChatLogWidth / 2),
			Y: g.Player.ApparentPos.Y * world.TileSize,
		}

		if g.shouldSetCam {
			g.ViewOffset.X = float64(ppos.X) - float64(width/2)
			g.ViewOffset.Y = float64(ppos.Y) - float64(height/2)

			g.shouldSetCam = false

		} else {
			g.ViewOffset.X = lerp(g.ViewOffset.X, float64(ppos.X)-float64(width/2), 0.01)
			g.ViewOffset.Y = lerp(g.ViewOffset.Y, float64(ppos.Y)-float64(height/2), 0.01)
		}

		if g.ViewOffset.X < 0 {
			g.ViewOffset.X = 0
		}

		if g.ViewOffset.Y < 0 {
			g.ViewOffset.Y = 0
		}

		if g.ViewOffset.X+float64(width) > (world.Width*world.TileSize)+ChatLogWidth {
			g.ViewOffset.X = float64(world.Width*world.TileSize - width + ChatLogWidth)
		}

		if g.ViewOffset.Y+float64(height) > world.Height*world.TileSize {
			g.ViewOffset.Y = float64(world.Height*world.TileSize - height)
		}
	}

	g.ChatLog.Render(rend, g.ld, width-ChatLogWidth, 0, ChatLogWidth, height)
}

// HandleEvent handles a window event, such as a mouse
// click or a key release.
func (g *Game) HandleEvent(event sdl.Event) string {
	switch evt := event.(type) {
	case *sdl.MouseButtonEvent:

		// If the left mouse button was clicked
		if evt.Type == sdl.MOUSEBUTTONDOWN && evt.Button == sdl.BUTTON_LEFT {
			if evt.X >= 1200-ChatLogWidth {
				break
			}

			x, y := g.ViewportToRelative(evt.X, evt.Y)
			tx, ty := g.PositionToTile(x, y)

			coord := geom.Coord{
				X: uint(tx),
				Y: uint(ty),
			}

			// Move the player locally
			g.Player.Move(coord, g.World)

			// Tell the server the player's moved
			g.Client.Send(&message.Moved{Position: coord})
		}

	case *sdl.KeyUpEvent:
		if evt.Keysym.Sym == sdl.K_ESCAPE {
			return "mainmenu"
		}

	case *sdl.KeyDownEvent:
		switch evt.Keysym.Sym {
		case sdl.K_BACKSPACE:
			if len(g.ChatLog.Input) > 0 {
				g.ChatLog.Input = g.ChatLog.Input[:len(g.ChatLog.Input)-1]
			}

		case sdl.K_RETURN:
			if g.Player == nil {
				break
			}

			g.Client.Send(&message.ChatMessage{
				Content: g.ChatLog.Input,
				Sender:  g.Player.Name,
				Time:    time.Now(),
				Type:    GlobalMessage,
			})

			g.ChatLog.Input = ""
		}

	case *sdl.TextInputEvent:
		str := ""

		// evt.Text is a null terminated c-string
		// str is the normal Go string
		for _, ch := range evt.Text {
			if ch == 0 {
				break
			}

			str += string(ch)
		}

		g.ChatLog.Input += str
	}

	return ""
}

// ViewportToRelative maps viewport positions to
// a coordinate relative to the World.
func (g *Game) ViewportToRelative(x, y int32) (int32, int32) {
	var (
		newX = x + int32(g.ViewOffset.X)
		newY = y + int32(g.ViewOffset.Y)
	)

	return newX, newY
}

// PositionToTile maps world coordinates to the
// tile at that position.
func (g *Game) PositionToTile(x, y int32) (int32, int32) {
	var (
		newX = x / world.TileSize
		newY = y / world.TileSize
	)

	return newX, newY
}

func (g *Game) tick() {
	g.World.Tick()

	for _, e := range g.Entities {
		e.Step()
	}
}

func (g *Game) serverUpdate() {
	if g.Player != nil {
		g.Client.Send(&message.StateUpdate{
			Position: g.Player.Pos,
		})
	}
}

func lerp(a, b, t float64) float64 {
	return (1-t)*a + t*b
}
