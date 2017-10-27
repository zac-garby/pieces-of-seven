package main

import (
	"strings"

	"github.com/Zac-Garby/pieces-of-seven/loader"
	"github.com/Zac-Garby/pieces-of-seven/scene"
	"github.com/Zac-Garby/pieces-of-seven/scene/game"
	"github.com/Zac-Garby/pieces-of-seven/scene/joingame"
	"github.com/Zac-Garby/pieces-of-seven/scene/mainmenu"
)

func makeScene(name string, ld *loader.Loader) scene.Scene {
	split := strings.Split(name, "\n")

	switch split[0] {
	case "join":
		return game.New(ld, split[1], split[2])
	case "mainmenu":
		return mainmenu.New(ld)
	case "joingame":
		return joingame.New(ld)
	default:
		panic("scene not found: " + name)
	}
}
