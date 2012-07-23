package main

import (
	"image/color"
	"github.com/ajhager/rog"
	"runtime"
)

var (
	width  = 48
	height = 32

	darkWall   = color.RGBA{0, 0, 100, 255}
	lightWall  = color.RGBA{130, 110, 50, 255}
	darkFloor  = color.RGBA{50, 50, 150, 255}
	lightFloor = color.RGBA{200, 180, 50, 255}

	fov   = rog.NewFOVMap(width, height)
	x     = 0
	y     = 16
	first = true
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
	w.Set(x, y, ' ', color.White, nil, rog.Normal)
	x = xx
	y = yy
	w.Set(x, y, '@', color.White, nil, rog.Normal)
	fov.Update(x, y, 10, true, rog.FOVCircular)
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
	}

	runtime.ReadMemStats(&stats)

	w.Print(0, 0, "                                                ")
	w.Print(0, 0, "%vFS %vMB %vGC %vGR", w.Fps, stats.Sys/1000000, stats.NumGC, runtime.NumGoroutine())

	w.Print(0, 31, "                                                ")
	w.Print(0, 31, "Pos: %v %v Cell: %v %v", w.Mouse.Pos.X, w.Mouse.Pos.Y, w.Mouse.Cell.X, w.Mouse.Cell.Y)

	if w.Mouse.Left.Released {
        movePlayer(w, w.Mouse.Cell.X, w.Mouse.Cell.Y)
	}

    switch w.Key {
    case "k":
        if tmap[y-1][x] != '#' {
            movePlayer(w, x, y - 1)
        }
    case "j":
        if tmap[y+1][x] != '#' {
            movePlayer(w, x, y + 1)
        }
    case "h":
        if tmap[y][x-1] != '#' {
            movePlayer(w, x - 1, y)
        }
    case "l":
        if tmap[y][x+1] != '#' {
            movePlayer(w, x + 1, y)
        }
    }

	for y := 0; y < fov.Height(); y++ {
		for x := 0; x < fov.Width(); x++ {
			if fov.Look(x, y) {
				if tmap[y][x] == '#' {
					w.Set(x, y, -1, nil, lightWall, rog.Normal)
				} else {
					w.Set(x, y, -1, nil, lightFloor, rog.Normal)
				}
			} else {
				if tmap[y][x] == '#' {
					w.Set(x, y, -1, nil, darkWall, rog.Normal)
				} else {
					w.Set(x, y, -1, nil, darkFloor, rog.Normal)
				}
			}
		}
	}
}

func main() {
	rog.Open(width, height, "FOV Example", fovExample)
	rog.Start()
}
