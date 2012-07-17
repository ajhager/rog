package main

import (
	"image/color"
	"github.com/ajhager/rog"
)

func popup(w *rog.Window) {
	w.Print("fps: %v", w.Fps)

	w.Set(19, 5, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 5, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 5, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Set)
	w.Set(21, 5, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Set)

	w.Set(19, 6, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 6, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 6, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Multiply)
	w.Set(21, 6, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Set)

	w.Set(19, 7, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 7, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 7, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Screen)
	w.Set(21, 7, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Set)

	w.Set(19, 8, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 8, -1, nil, color.RGBA{100, 100, 200, 255}, rog.Set)
	w.Set(20, 8, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Overlay)
	w.Set(21, 8, -1, nil, color.RGBA{200, 100, 100, 255}, rog.Set)
}

func game(w *rog.Window) {
	grey := color.RGBA{20, 20, 20, 255}
	for x := 0; x < w.Width(); x++ {
		for y := 0; y < w.Height(); y++ {
			purple := color.RGBA{uint8(200), uint8(y * x), uint8((y * 4) % 255), 255}
			w.Set(x, y, rune(1000+(2*(x+(x*y)))+(y+x*y)), purple, grey, rog.Set)
		}
	}
	w.Print("fps: %v", w.Fps)
}

func main() {
	// rog.Open(48, 32, "Basic Example", game)
	rog.Open(48, 32, "Second Window", popup)
	rog.Start()
}
