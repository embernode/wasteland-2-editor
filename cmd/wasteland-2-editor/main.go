package main

import (
	"fyne.io/fyne/v2/app"

	"github.com/embernode/wasteland-2-editor/internal/ui"
)

func main() {
	a := app.NewWithID("io.github.ncl8.wl2edit")
	w := a.NewWindow("Wasteland 2 Save Editor")
	ui.BuildMainWindow(w)
	w.ShowAndRun()
}
