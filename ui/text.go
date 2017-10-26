package ui

import (
	"math"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// A Text element displays a line of text.
type Text struct {
	Text      string
	Font      *ttf.Font
	Rect      *sdl.Rect
	Alignment Alignment

	R, G, B uint8
}

func NewText(text string, r, g, b uint8, font *ttf.Font, align Alignment) *Text {
	return &Text{
		Text: text,
		Font: font,
		Rect: &sdl.Rect{
			X: 0, Y: 0,
			W: 0, H: 0,
		},
		Alignment: align,
		R:         r,
		G:         g,
		B:         b,
	}
}

// SetRect only modifies the coordinates in this
// case.
func (t *Text) SetRect(r *sdl.Rect) {
	t.Rect = r
}

// GetRect returns the text's rectangle.
func (t *Text) GetRect() *sdl.Rect {
	return t.Rect
}

// Render draws the text to the render, disregarding
// the width and height components of the rect.
func (t *Text) Render(rend *sdl.Renderer) {
	solid, err := t.Font.RenderUTF8_Solid(t.Text, sdl.Color{R: t.R, G: t.G, B: t.B, A: 255})
	if err != nil {
		panic(err)
	}

	tex, err := rend.CreateTextureFromSurface(solid)
	if err != nil {
		panic(err)
	}

	srect := solid.ClipRect
	dest := &solid.ClipRect
	dest.X = t.Rect.X
	dest.Y = t.Rect.Y

	// Center vertically
	dest.Y += (t.Rect.H - srect.H) / 2

	if t.Alignment == LeftAlign {
		// Do nothing, since it's already left aligned
	}

	if t.Alignment == CenterAlign {
		dest.X += (t.Rect.W - srect.W) / 2
	}

	if t.Alignment == RightAlign {
		dest.X += t.Rect.W - srect.W
	}

	// Draw the texture
	rend.Copy(tex,
		&sdl.Rect{
			0,
			0,
			t.Rect.W,
			t.Rect.H,
		},

		&sdl.Rect{
			X: dest.X,
			Y: dest.Y,
			W: int32(math.Min(float64(dest.W), float64(t.Rect.W))),
			H: int32(math.Min(float64(dest.H), float64(t.Rect.H))),
		},
	)

	solid.Free()
	tex.Destroy()
}

// Update does nothing, since text doesn't need to be updated.
func (t *Text) Update(float64, uint, uint) {}
