package loader

import (
	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// The Loader loads the assets required in the game.
type Loader struct {
	loadQueue map[string]Asset

	// Outputs
	Surfaces map[string]*sdl.Surface
	Textures map[string]*sdl.Texture
	Fonts    map[string]*ttf.Font
}

const (
	Texture = iota
	Font
)

// An Asset contains the necessary data for loading
// an asset.
type Asset struct {
	Path string
	Type int
	Data map[string]int
}

// New creates a new Loader.
func New() Loader {
	return Loader{
		loadQueue: make(map[string]Asset),

		Surfaces: make(map[string]*sdl.Surface),
		Textures: make(map[string]*sdl.Texture),
		Fonts:    make(map[string]*ttf.Font),
	}
}

// Queue adds the assets to the load queue.
func (l *Loader) Queue(assets map[string]Asset) {
	for name, asset := range assets {
		l.loadQueue[name] = asset
	}
}

// Load starts to load all the images in the
// load stack.
func (l *Loader) Load(rend *sdl.Renderer) {
	for name, asset := range l.loadQueue {
		switch asset.Type {
		case Texture:
			l.Surfaces[name], l.Textures[name] = l.loadImage(asset.Path, rend)

		case Font:
			l.Fonts[name] = l.loadFont(asset.Path, asset.Data["size"])
		}
	}
}

func (l *Loader) loadImage(path string, rend *sdl.Renderer) (*sdl.Surface, *sdl.Texture) {
	image, err := img.Load(path)
	if err != nil {
		panic(err)
	}

	texture, err := rend.CreateTextureFromSurface(image)
	if err != nil {
		panic(err)
	}

	return image, texture
}

func (l *Loader) loadFont(path string, size int) *ttf.Font {
	font, err := ttf.OpenFont(path, size)
	if err != nil {
		panic(err)
	}

	return font
}

// Free frees the resources used by the loaded
// assets.
func (l *Loader) Free() {
	for _, surface := range l.Surfaces {
		surface.Free()
	}

	for _, texture := range l.Textures {
		texture.Destroy()
	}
}
