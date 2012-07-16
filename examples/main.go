package main

import (
	"image/color"
	"github.com/ajhager/rog"
)

func popup(w *rog.Window) {
	w.Clear()
	red := color.RGBA{200, 50, 0, 255}
	w.Print("POPUP")
	w.Set(10, 10, '*', red, nil)
	w.SetTitle("Popup")
}

func game(w *rog.Window) {
	grey := color.RGBA{20, 20, 20, 255}
	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			purple := color.RGBA{uint8(200), uint8(y * x), uint8((y * 4) % 255), 255}
			w.Set(x, y, rune(1000+(2*(x+(x*y)))+(y+x*y)), purple, grey)
		}
	}
}

func main() {
	rog.Open(48, 32, "Basic Example", game)
	rog.Open(48, 32, "Second Window", popup)
	rog.Start()
}
