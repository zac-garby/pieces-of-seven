package main

import (
	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/scene"
	"github.com/Zac-Garby/pieces-of-seven/scene/game"
	"github.com/Zac-Garby/pieces-of-seven/scene/mainmenu"
)

func makeScene(name string, ld *loader.Loader) scene.Scene {
	switch name {
	case "join":
		return game.New(ld, ":12358")
	case "mainmenu":
		return mainmenu.New(ld)
	default:
		panic("scene not found: " + name)
	}
}
