package joingame

import (
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// JoinGame is entered when a player wants to join
// a server. It prompts him for the ip to connect to.
type JoinGame struct {
	ld    *loader.Loader
	inter *ui.Interface
}

// New creates a new JoinGame scene.
func New(ld *loader.Loader) *JoinGame {
	join := &JoinGame{
		ld: ld,
		inter: &ui.Interface{
			Padding: 5,
		},
	}

	join.inter.Add("text", ui.NewText(
		"Enter an IP to connect to:",
		255, 255, 255,
		ld.Fonts["body"],
		ui.CenterAlign,
	))

	join.inter.Add("addr", ui.NewTextfield(
		"127.0.0.1:12358",
		ld.Fonts["body"],
		ui.CenterAlign,
	))

	join.inter.Layout(300, 100, 600, 40)

	return join
}

// Enter is called when a JoinGame scene is entered.
func (j *JoinGame) Enter() {}

// Exit is called when a JoinGame scene is exited.
func (j *JoinGame) Exit() {}

// Update updates the scene by 'dt' seconds.
func (j *JoinGame) Update(dt float64) string {
	j.inter.Update(dt)

	return ""
}

// Render renders the scene to an SDL renderer.
func (j *JoinGame) Render(rend *sdl.Renderer, width, height int) {
	j.inter.Render(rend)
}

// HandleEvent handles an SDL event. If it returns a non-empty
// string, the game changes to that scene.
func (j *JoinGame) HandleEvent(event sdl.Event) string {
	switch evt := event.(type) {
	case *sdl.KeyUpEvent:
		switch evt.Keysym.Sym {
		case sdl.K_RETURN:
			return "join"
		}

	default:
		j.inter.HandleEvent(event)
	}

	return ""
}
