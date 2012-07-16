package main

import (
	"image/color"
	"github.com/ajhager/rog"
)

func popup(w *rog.Window) {
	w.Print("fps: %v", w.Fps)
}

func game(w *rog.Window) {
	grey := color.RGBA{20, 20, 20, 255}
	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			purple := color.RGBA{uint8(200), uint8(y * x), uint8((y * 4) % 255), 255}
			w.Set(x, y, rune(1000+(2*(x+(x*y)))+(y+x*y)), purple, grey)
		}
	}
	w.Print("fps: %v", w.Fps)
}

func main() {
	// rog.Open(48, 32, "Basic Example", game)
	rog.Open(48, 32, "Second Window", popup)
	rog.Start()
}
