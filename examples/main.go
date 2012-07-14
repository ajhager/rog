package main

import (
	"image/color"
	"github.com/ajhager/rog"
)

func popup(w *rog.Window) {
	red := color.RGBA{200, 50, 0, 255}
	w.SetChFg(2, 5, 'P', red)
	w.SetChFg(3, 4, 'O', red)
	w.SetChFg(4, 3, 'P', red)
	w.SetChFg(5, 4, 'U', red)
	w.SetChFg(6, 5, 'P', red)
}

func game(w *rog.Window) {
	grey := color.RGBA{20, 20, 20, 255}
	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			purple := color.RGBA{uint8(200), uint8(y * x), uint8((y * 4) % 255), 255}
			w.Set(x, y, rune(1000+(2*(x+(x*y)))+(y+x*y)), grey, purple)
		}
	}
}

func main() {
	rog.Open(10, 10, "Popup", popup)
	rog.Open(48, 32, "Basic Example", game)
	rog.Start()
}
