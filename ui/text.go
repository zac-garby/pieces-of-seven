package ui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// A Text element displays a line of text.
type Text struct {
	Text string
	Font *ttf.Font
	Rect *sdl.Rect
}

func NewText(text string, font *ttf.Font) *Text {
	return &Text{
		Text: text,
		Font: font,
		Rect: &sdl.Rect{
			X: 0, Y: 0,
			W: 0, H: 0,
		},
	}
}

// SetRect only modifies the coordinates in this
// case.
func (t *Text) SetRect(r *sdl.Rect) {
	t.Rect.X = r.X
	t.Rect.Y = r.Y
}

// GetRect returns the text's rectangle.
func (t *Text) GetRect() *sdl.Rect {
	return t.Rect
}

// Render draws the text to the render, disregarding
// the width and height components of the rect.
func (t *Text) Render(rend *sdl.Renderer) {
	solid, err := t.Font.RenderUTF8_Solid(t.Text, sdl.Color{R: 255, G: 255, B: 255, A: 255})
	if err != nil {
		panic(err)
	}

	tex, err := rend.CreateTextureFromSurface(solid)
	if err != nil {
		panic(err)
	}

	dest := &solid.ClipRect
	dest.X = t.Rect.X
	dest.Y = t.Rect.Y

	rend.Copy(tex, nil, dest)

	solid.Free()
	tex.Destroy()
}

// Update does nothing, since text doesn't need to be updated.
func (t *Text) Update(float64, uint, uint) {}
