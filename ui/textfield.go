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

// SetRect sets the text field's bounding rect to r
// and the inner text's bounding rect to a slightly
// smaller, translated version of r.
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

func (t *Textfield) Update(float64) {
	if len(t.Text) > 0 {
		t.text.Text = t.Text

		t.text.R = 255
		t.text.G = 255
		t.text.B = 255
	} else {
		t.text.Text = t.Placeholder

		t.text.R = 128
		t.text.G = 128
		t.text.B = 128
	}
}

func (t *Textfield) HandleEvent(event sdl.Event) {
	switch evt := event.(type) {
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

		t.Text += str
	}
}

func (t *Textfield) Activate() {
	sdl.StartTextInput()
}

func (t *Textfield) Deactivate() {
	sdl.StopTextInput()
}
