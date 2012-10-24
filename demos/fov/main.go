package main

import (
	"github.com/ajhager/rog"
	"image"
	"runtime"
)

type Floor struct { rog.Tile }
func NewFloor(x, y int) *Floor {
    return &Floor{
        rog.Tile {
            X: x, Y: y,
            Fg: rog.Hex(0x442020), 
            Bg: rog.Hex(0x885040),
            Glyph: '.', 
            Roughness: rog.PATH_MIN, 
            Viewable: true,
        },
    }
}

type Wall struct { rog.Tile }
func NewWall(x, y int) *Wall {
    return &Wall{
    	rog.Tile {
		    X: x, Y: y,
		    Fg: rog.Hex(0x885544), 
		    Bg: rog.Hex(0xffbb99),
		    Glyph: '#', 
		    Roughness: rog.PATH_MAX, 
		    Viewable: false,
		},
	}
}


var (
	width  = 40
	height = 20

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
					pmap.SetTile(NewWall(x, y), x, y)
				} else {
					pmap.SetTile(NewFloor(x, y), x, y)
				}
			}
		}
		movePlayer(23, 9)
	}

	if rog.Mouse().Left.Released {
		path = pmap.Path(x, y, rog.Mouse().Cell.X, rog.Mouse().Cell.Y)
	}

	switch rog.Key() {
	case 'k':
		movePlayer(x, y-1)
	case 'j':
		movePlayer(x, y+1)
	case 'h':
		movePlayer(x-1, y)
	case 'l':
		movePlayer(x+1, y)
	case rog.Esc:
		rog.Close()
	}

	//floor := NewFloor(0, 0)

	for cy := 0; cy < pmap.Height(); cy++ {
		for cx := 0; cx < pmap.Width(); cx++ {
			rog.Set(cx, cy, nil, rog.Black, " ")
		}
	}

	for point, _ := range pmap.ViewMap() {
		tile, _ := pmap.GetTile(point.X, point.Y)
		i := intensity(x, y, point.X, point.Y, 20)
		rog.Set(point.X, point.Y, tile.GetFg().Scale(i), tile.GetBg().Scale(i), string(tile.GetGlyph()))
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
	rog.Open(width, height, 1, false, "Example", nil)
	for rog.Running() {
		example()
		rog.Flush()
	}
}
