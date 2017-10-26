package ui

import (
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Textfield struct {
	Text        string
	Placeholder string
	Font        *ttf.Font
	Rect        *sdl.Rect
	Alignment   Alignment

	text *Text
}

func NewTextfield(placeholder string, font *ttf.Font, align Alignment) *Textfield {
	return &Textfield{
		Text:        "",
		Placeholder: placeholder,
		Alignment:   align,

		Rect: &sdl.Rect{
			X: 0, Y: 0,
			W: 0, H: 0,
		},

		text: NewText(
			placeholder,
			128, 128, 128,
			font,
			LeftAlign,
		),
	}
}

func (t *Textfield) SetRect(r *sdl.Rect) {
	t.Rect = r

	t.text.SetRect(&sdl.Rect{
		X: r.X + 10,
		Y: r.Y,
		W: r.W - 10,
		H: r.H,
	})
}

func (t *Textfield) GetRect() *sdl.Rect {
	return t.Rect
}

func (t *Textfield) Render(rend *sdl.Renderer) {
	rend.SetDrawColor(255, 255, 255, 255)
	rend.DrawRect(t.Rect)

	t.text.Render(rend)
}

func (t *Textfield) Update(float64, uint, uint) {
	t.Text = t.text.Text
}
