package main

import (
	"image/color"
	"github.com/ajhager/rog"
	"runtime"
)

var (
	width  = 48
	height = 32

	darkWall   = color.RGBA{0, 0, 0, 255}
	lightWall = color.RGBA{255, 180, 180, 255}
	darkFloor  = color.RGBA{0, 0, 0, 255}
	lightFloor  = color.RGBA{128, 70, 70, 255}
    grey  = color.RGBA{200, 200, 200, 255}
    umbra = color.RGBA{30, 30, 30, 255}

	fov   = rog.NewFOVMap(width, height)
	x     = 0
	y     = 16
	first = true
	row = "                                                "
	tmap  = [][]rune{
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("####################    ########################"),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("    #               ####        #  #  #         "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("##################              #######         "),
		[]rune("#                                               "),
		[]rune("#                #                              "),
		[]rune("#                #                              "),
		[]rune("#################### ## ## ## ## ## ## ## ## ###"),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
		[]rune("                                                "),
	}

	stats runtime.MemStats
)

func movePlayer(w *rog.Window, xx, yy int) {
    if xx >= 0 && yy > 0 && xx < width && yy < height-1 && tmap[yy][xx] == ' ' {
	    w.Set(x, y, " ", color.White, nil)
	    x = xx
	    y = yy
	    fov.Update(x, y, 20, true, rog.FOVCircular)
    }
}

func intensity(px, py, cx, cy, r int) float64 {
    r2 := float64(r * r)
    squaredDist := float64((px-cx)*(px-cx)+(py-cy)*(py-cy))
    coef1 := 1.0 / (1.0 + squaredDist / 20)
    coef2 := coef1 - 1.0 / (1.0 + r2)
    return coef2 / (1.0 - 1.0 / (1.0 + r2))
}

func fovExample(w *rog.Window) {
	if first {
		first = false
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if tmap[y][x] == '#' {
					fov.Block(x, y, true)
				}
			}
		}
        movePlayer(w, 27, 16)
	}

	if w.Mouse.Left.Released {
        movePlayer(w, w.Mouse.Cell.X, w.Mouse.Cell.Y)
	}

    switch w.Key {
    case "k":
        movePlayer(w, x, y - 1)
    case "j":
        movePlayer(w, x, y + 1)
    case "h":
        movePlayer(w, x - 1, y)
    case "l":
        movePlayer(w, x + 1, y)
    }

	for cy := 0; cy < fov.Height(); cy++ {
		for cx := 0; cx < fov.Width(); cx++ {
            w.Set(cx, cy, "", nil, darkWall)
			w.Set(cx, cy, " ", nil, darkFloor)
			if fov.Look(cx, cy) {
                i := intensity(x, y, cx, cy, 20)
				if tmap[cy][cx] == '#' {
					w.Set(cx, cy, "", nil, rog.ColorMul(lightWall, i))
				} else {
					w.Set(cx, cy, "✵", rog.ColorMul(lightFloor, i*1.5), rog.ColorMul(lightFloor, i), rog.Screen)
				}
			}
		}
	}
	w.Set(x, y, "웃", grey, nil)

	runtime.ReadMemStats(&stats)
    w.Fill(0, 0, w.Width(), 1, ' ', grey, umbra, rog.Dodge)
    w.Set(0, 0, w.P("%vFS %vMB %vGC %vGR", w.Fps, stats.Sys/1000000, stats.NumGC, runtime.NumGoroutine()), nil, nil)
	w.Set(0, 31, row, grey, umbra, rog.Dodge)
	w.Set(0, 31, w.P("Pos: %v %v Cell: %v %v", w.Mouse.Pos.X, w.Mouse.Pos.Y, w.Mouse.Cell.X, w.Mouse.Cell.Y), nil, nil)
}

func main() {
	rog.Open(width, height, "FOV Example", fovExample)
	rog.Start()
}
