package mainmenu

import (
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/ui"
	"github.com/veandco/go-sdl2/sdl"
)

// The MainMenu is the first scene to be loaded,
// and renders buttons to go the various other
// scenes.
type MainMenu struct {
	ld    *loader.Loader
	inter *ui.Interface
}

// New creates a new MainMenu instance.
func New(ld *loader.Loader) *MainMenu {
	menu := &MainMenu{
		ld: ld,
		inter: &ui.Interface{
			Padding: 5,
		},
	}

	menu.inter.Add("host", ui.NewText(
		"Host a server via the server package.",
		ld.Fonts["body"],
		ui.LeftAlign,
	))

	menu.inter.Add("join", ui.NewText(
		"Press [C] to join one.",
		ld.Fonts["body"],
		ui.LeftAlign,
	))

	menu.inter.Layout(30, 30, 1200-60, 40)

	return menu
}

// Enter is called when a MainMenu is entered.
func (m *MainMenu) Enter() {}

// Exit is called when a MainMenu is exited.
func (m *MainMenu) Exit() {}

// Update updates the main menu by 'dt' seconds. The returned
// scene will be changed to in the main loop, or, if nil is
// returned, the scene won't be changed.
func (m *MainMenu) Update(dt float64) string {
	m.inter.Update(dt, 0, 0)

	return ""
}

// Render renders the main menu to an SDL renderer.
func (m *MainMenu) Render(rend *sdl.Renderer, width, height int) {
	m.inter.Render(rend)
}

// HandleEvent handles a single event. If it returns
// a non-empty string, the scene switches to that scene.
func (m *MainMenu) HandleEvent(event sdl.Event) string {
	switch evt := event.(type) {
	case *sdl.KeyUpEvent:
		switch evt.Keysym.Sym {
		case sdl.K_c:
			return "joingame"
		}
	}

	return ""
}
