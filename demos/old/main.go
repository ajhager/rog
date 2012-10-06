package main

import (
	"github.com/ajhager/rog"
	_ "github.com/ajhager/rog/glfw"
	"image"
	"runtime"
)

var (
	width  = 40
	height = 20

	wall   = rog.Hex(0xffbb99)
	floorc = rog.Hex(0x885040)
	lgrey  = rog.Hex(0xc8c8c8)
	dgrey  = rog.Hex(0x1e1e1e)

	pmap  = rog.NewMap(width, height)
	path  []image.Point
	x     = 0
	y     = 0
	first = true
	stats runtime.MemStats

	tmap = [][]rune{
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("################    ####################"),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("    #           ####        #  #  #     "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("##############              #######     "),
		[]rune("#                                       "),
		[]rune("#            #                          "),
		[]rune("#            #                          "),
		[]rune("################ ## ## ## ## ## ## ## ##"),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
	}
)

func movePlayer(xx, yy int) {
	if xx >= 0 && yy > 0 && xx < width && yy < height-1 && tmap[yy][xx] == ' ' {
		rog.Set(x, y, rog.White, nil, " ")
		x = xx
		y = yy
		pmap.Fov(x, y, 20, true, rog.FOVCircular)
	}
}

func intensity(px, py, cx, cy, r int) float64 {
	r2 := float64(r * r)
	squaredDist := float64((px-cx)*(px-cx) + (py-cy)*(py-cy))
	coef1 := 1.0 / (1.0 + squaredDist/20)
	coef2 := coef1 - 1.0/(1.0+r2)
	return coef2 / (1.0 - 1.0/(1.0+r2))
}

func example() {
	if first {
		first = false
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if tmap[y][x] == '#' {
					pmap.Block(x, y, true)
				}
			}
		}
		movePlayer(23, 9)
	}

	if rog.Mouse().Left.Released {
		path = pmap.Path(x, y, rog.Mouse().Cell.X, rog.Mouse().Cell.Y)
	}

	switch rog.Key() {
	case 'K':
		movePlayer(x, y-1)
	case 'J':
		movePlayer(x, y+1)
	case 'H':
		movePlayer(x-1, y)
	case 'L':
		movePlayer(x+1, y)
	case rog.Escape:
		rog.Close()
	}

	for cy := 0; cy < pmap.Height(); cy++ {
		for cx := 0; cx < pmap.Width(); cx++ {
			rog.Set(cx, cy, nil, rog.Black, " ")
			if pmap.Look(cx, cy) {
				i := intensity(x, y, cx, cy, 20)
				if tmap[cy][cx] == '#' {
					rog.Set(cx, cy, nil, wall.Scale(i), "")
				} else {
					rog.Set(cx, cy, rog.Scale(1.5), floorc.Scale(i), "✵")
				}
			}
		}
	}
	rog.Set(x, y, lgrey, nil, "웃")

	for _, p := range path {
		if p.X != x || p.Y != y {
			rog.Set(p.X, p.Y, lgrey, nil, "*")
		}
	}

	runtime.ReadMemStats(&stats)
	rog.Fill(0, 0, rog.Width(), 1, lgrey, rog.Dodge(dgrey), ' ')
	rog.Set(0, 0, nil, nil, "%vFPS %vMB %vGC %vGR", rog.Fps(), stats.HeapAlloc/1000000, stats.NumGC, runtime.NumGoroutine())
	rog.Fill(0, height-1, rog.Width(), 1, lgrey, rog.Dodge(dgrey), ' ')
	rog.Set(0, height-1, nil, nil, "Pos: %v %v Cell: %v %v", rog.Mouse().Pos.X, rog.Mouse().Pos.Y, rog.Mouse().Cell.X, rog.Mouse().Cell.Y)
}

func main() {
	rog.Open(width, height, 1, "Example", "../../data/font.png")
	for rog.IsOpen() {
		example()
		rog.Flush()
	}
}
