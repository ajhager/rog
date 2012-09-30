package main

import (
	"github.com/ajhager/rog"
	"github.com/ajhager/rog/wde"
	"math/rand"
)

const (
	topLeft  = 0
	topRight = 1
	botLeft  = 2
	botRight = 3
	width    = 40
	height   = 20
)

var (
	sampleConsole *rog.Console = rog.NewConsole(width, height)
	colors        []rog.RGB    = []rog.RGB{
		rog.RGB{50, 40, 150},
		rog.RGB{240, 85, 5},
		rog.RGB{50, 35, 240},
		rog.RGB{10, 200, 130}}
	dirR  []int   = []int{1, -1, 1, 1}
	dirG  []int   = []int{1, -1, -1, 1}
	dirB  []int   = []int{1, 1, 1, -1}
	black rog.RGB = rog.RGB{0, 0, 0}
)

func render() {
	for c := 0; c < 4; c++ {
		switch rand.Int31n(3) {
		case 0:
			colors[c].R += uint8(5 * dirR[c])
			if colors[c].R == 255 {
				dirR[c] = -1
			} else if colors[c].R == 0 {
				dirR[c] = 1
			}
		case 1:
			colors[c].G += uint8(5 * dirG[c])
			if colors[c].G == 255 {
				dirG[c] = -1
			} else if colors[c].G == 0 {
				dirG[c] = 1
			}
		case 2:
			colors[c].B += uint8(5 * dirB[c])
			if colors[c].B == 255 {
				dirB[c] = -1
			} else if colors[c].B == 0 {
				dirB[c] = 1
			}
		}
	}

	for x := 0; x < width; x++ {
		xcoef := float64(x) / float64(width-1)
		top := colors[topLeft].Alpha(colors[topRight], xcoef)
		bot := colors[botLeft].Alpha(colors[botRight], xcoef)
		for y := 0; y < height; y++ {
			ycoef := float64(y) / float64(height-1)
			cur := top.Alpha(bot, ycoef)
			sampleConsole.Set(x, y, nil, cur, "")
		}
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			_, col, _ := sampleConsole.Get(x, y)
			col = col.Alpha(black, 0.5)
			sampleConsole.Set(x, y, col, nil, string(rand.Int31n(26)+97))
		}
	}
}

func main() {
	rog.Open(width, height+2, 1, "tcod true color", wde.Backend())
	for rog.IsOpen() {
		render()
		rog.Blit(sampleConsole, 0, 1)
		rog.Set(0, height+1, nil, nil, "%v", rog.Fps())
		if rog.Key() == rog.Escape {
			rog.Close()
		}
		rog.Flush()
	}
}
