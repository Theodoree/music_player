package main

import (
	"embed"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"github.com/Theodoree/music_player/internal/gui"
	"math/rand"
	"time"
)

//go:embed resource
var resource embed.FS

func main() {
	rand.Seed(time.Now().UnixNano())
	a := app.NewWithID("io.fyne.music_player")
	a.SetIcon(theme.HomeIcon())
	gui.Run(a, &resource)
}
