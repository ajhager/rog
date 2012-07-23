package main

import (
	"image/color"
	"github.com/ajhager/rog"
	"runtime"
)

var (
	width  = 48
	height = 32

	darkWall   = color.RGBA{40, 40, 40, 255}
	lightWall = color.RGBA{165, 120, 150, 255}
	darkFloor  = color.RGBA{20, 15, 17, 255}
	lightFloor  = color.RGBA{100, 70, 90, 255}
	lightFloor2  = color.RGBA{90, 60, 80, 255}
    grey  = color.RGBA{200, 200, 200, 255}
    umbra = color.RGBA{30, 20, 10, 255}

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

	for y := 0; y < fov.Height(); y++ {
		for x := 0; x < fov.Width(); x++ {
			if fov.Look(x, y) {
				if tmap[y][x] == '#' {
					w.Set(x, y, "", lightWall, lightWall)
				} else {
					w.Set(x, y, "✵", lightFloor2, lightFloor)
				}
			} else {
				if tmap[y][x] == '#' {
					w.Set(x, y, "", nil, darkWall)
				} else {
					w.Set(x, y, " ", nil, darkFloor)
				}
			}
		}
	}
	w.Set(x, y, "@", grey, nil)

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
