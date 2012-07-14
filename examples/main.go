package main

import (
	"image/color"
	"github.com/ajhager/rog"
)

func game(w *rog.Window) {
	grey := color.RGBA{20, 20, 20, 255}
	for x := 0; x < 40; x++ {
		for y := 0; y < 20; y++ {
			purple := color.RGBA{uint8(150 + x%255), uint8(y * x), uint8((y * 4) % 255), 255}
			w.Set(x, y, rune(1000+(2*(x+(x*y)))+(y+x*y)), grey, purple)
		}
	}
}

func main() {
	rog.Open(40, 20, "Basic Example", game)
	rog.Start()
}
