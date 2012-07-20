package main

import (
	"image/color"
	"github.com/ajhager/rog"
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
	dx    = 1
	time  = float64(0)
	first = true
	i     = 0
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

func fovExample(w *rog.Window) {
	i += 1
	time += w.Dt

	w.Print("%v    ", w.Fps)

	if first {
		first = false
		w.Set(x, y, '@', color.White, nil, rog.Normal)
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				if tmap[y][x] == '#' {
					fov.Block(x, y, true)
				}
			}
		}
	}

	if time >= .25 {
		time = 0
		w.Set(x, y, ' ', color.White, nil, rog.Normal)
		x += dx
		if x == (width-1) || x == 0 {
			dx = -dx
		}
		w.Set(x, y, '@', color.White, nil, rog.Normal)
		fov.Update(x, y, 10, true, rog.FOVCircular)
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
