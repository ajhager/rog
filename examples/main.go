package main

import (
	"github.com/ajhager/rog"
	"runtime"
)

var (
	width  = 48
	height = 32

	wall = rog.HEX(0xffb4b4)
	floor = rog.HEX(0x804646)
    black = rog.HEX(0x000000)
    white = rog.HEX(0xffffff)
    lgrey = rog.HEX(0xc8c8c8)
    dgrey = rog.HEX(0x1e1e1e)

	fov   = rog.NewFOVMap(width, height)
	x     = 0
	y     = 0
	first = true
	stats runtime.MemStats

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
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

func movePlayer(w *rog.Window, xx, yy int) {
    if xx >= 0 && yy > 0 && xx < width && yy < height-1 && tmap[yy][xx] == ' ' {
	    w.Set(x, y, white, nil, " ")
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
    case "p":
        w.Screenshot("test")
    case "escape":
        w.Close()
    }

	for cy := 0; cy < fov.Height(); cy++ {
		for cx := 0; cx < fov.Width(); cx++ {
			w.Set(cx, cy, nil, black, " ")
			if fov.Look(cx, cy) {
                i := intensity(x, y, cx, cy, 20)
				if tmap[cy][cx] == '#' {
					w.Set(cx, cy, nil, wall.Scale(i), "")
				} else {
					w.Set(cx, cy, floor.Scale(i*1.5), floor.Scale(i), "✵")
				}
			}
		}
	}
	w.Set(x, y, lgrey, nil, "웃")

	runtime.ReadMemStats(&stats)
    w.Fill(0, 0, w.Width(), 1, ' ', lgrey, rog.Dodge(dgrey))
    w.Set(0, 0, nil, nil, "%vFS %vMB %vGC %vGR", w.Fps, stats.Sys/1000000, stats.NumGC, runtime.NumGoroutine())
    w.Fill(0, 31, w.Width(), 32, ' ', lgrey, rog.Dodge(dgrey))
	w.Set(0, 31, nil, nil, "Pos: %v %v Cell: %v %v", w.Mouse.Pos.X, w.Mouse.Pos.Y, w.Mouse.Cell.X, w.Mouse.Cell.Y)
}

func main() {
	rog.Open(width, height, "FOV Example", fovExample)
	rog.Start()
}
