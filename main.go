package main

import (
	"embed"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/Theodoree/music_player/internal/gui"
)

//go:embed resource
var resource embed.FS

func main() {
	a := app.NewWithID("io.fyne.music_player")
	a.SetIcon(theme.HomeIcon())
	gui.Run(a, &resource)
}
