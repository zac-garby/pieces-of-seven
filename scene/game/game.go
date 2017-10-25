package game

import (
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

	ld         *loader.Loader
	nextTick   float64
	shouldQuit bool
}

// New creates a new Game instance.
func New(ld *loader.Loader, addr string) *Game {
	game := &Game{
		World:      &world.World{},
		ViewOffset: &geom.Vector{X: 0, Y: 0},
		nextTick:   1.0 / TickRate,
		ld:         ld,
		shouldQuit: false,
		Players:    make(map[uuid.UUID]*entity.Ship),
	}

	game.Client = NewClient(addr, game)

	return game
}

// Enter is called when a Game scene is entered.
func (g *Game) Enter() {
	go g.Client.Listen()
}

// Exit is called when a Game scene is exited.
func (g *Game) Exit() {
	if !g.shouldQuit {
		g.Client.LeaveGame()
	}
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

	if g.nextTick <= 0 {
		g.nextTick = 1.0 / TickRate
		g.tick()
	}

	for _, e := range g.Entities {
		e.Update(dt)
	}

	return ""
}

// Render renders a game (i.e. the objects inside it)
// onto an SDL renderer.
func (g *Game) Render(rend *sdl.Renderer, width, height int) {
	g.World.Render(rend, g.ViewOffset, width, height)

	for _, e := range g.Entities {
		e.Render(g.ViewOffset, g.ld, rend)
	}
}

// HandleEvent handles a window event, such as a mouse
// click or a key release.
func (g *Game) HandleEvent(event sdl.Event) string {
	switch evt := event.(type) {
	case *sdl.MouseButtonEvent:

		// If the left mouse button was clicked
		if evt.Type == sdl.MOUSEBUTTONDOWN && evt.Button == sdl.BUTTON_LEFT {
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
		return "mainmenu"
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
	for _, e := range g.Entities {
		e.Step()
	}
}
