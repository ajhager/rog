package rog

import (
	"fmt"
	"image"
)

type Point image.Point 

type Viewable interface {
	IsViewable() bool
}

type ViewableMap interface {
	In(x, y int) bool
	Width() int
	Height() int
	Viewable(x, y int) bool
	ViewMap() ViewMap
}

type ViewMap map[Point]bool

func (self ViewMap) Update(other ViewMap) {
	for k, v := range other{
		self[k] = v
	}
}

// FOVAlgo takes a FOVMap x,y vantage, radius of the view, whether to include walls and then marks in the map which cells are viewable.
type FOVAlgo func(ViewableMap, int, int, int, bool) ViewMap

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Circular Raycasting
func fovCircularCastRay(fov ViewableMap, xo, yo, xd, yd, r2 int, walls bool) ViewMap {
	data := make(ViewMap, 0)
	curx := xo
	cury := yo
	in := false
	blocked := false
	if fov.In(curx, cury) {
		in = true
		data[Point{curx, cury}] = true
	}
	for _, p := range Line(xo, yo, xd, yd) {
		curx = p.X
		cury = p.Y
		if r2 > 0 {
			curRadius := (curx-xo)*(curx-xo) + (cury-yo)*(cury-yo)
			if curRadius > r2 {
				break
			}
		}
		if fov.In(curx, cury) {
			in = true
			if !blocked && !fov.Viewable(curx, cury) {
				blocked = true
			} else if blocked {
				break
			}
			if walls || !blocked {
				data[Point{curx, cury}] = true
			}
		} else if in {
			break
		}
	}
	return data
}

func fovCircularPostProc(fov ViewableMap, vdata ViewMap, x0, y0, x1, y1, dx, dy int) {
	data := make(ViewMap, 0)
	for cx := x0; cx <= x1; cx++ {
		for cy := y0; cy <= y1; cy++ {
			x2 := cx + dx
			y2 := cy + dy
			seen := vdata[Point{cx, cy}]
			if fov.In(cx, cy) && seen && !fov.Viewable(cx, cy) {
				if x2 >= x0 && x2 <= x1 {
					if fov.In(x2, cy) && fov.Viewable(x2, cy) {
						fmt.Print('.')
						data[Point{x2, cy}] = true
					}
				}
				if y2 >= y0 && y2 <= y1 {
					if fov.In(cx, y2) && fov.Viewable(cx, y2) {
						fmt.Print('.')
						data[Point{cx, y2}] = true
					}
				}
				if x2 >= x0 && x2 <= x1 && y2 >= y0 && y2 <= y1 {
					if fov.In(x2, y2) && fov.Viewable(x2, y2) {
						fmt.Print('.')
						data[Point{x2, y2}] = true
					}
				}
			}
		}
	}
	vdata.Update(data)
}

// FOVCicular raycasts out from the vantage in a circle.
func FOVCircular(fov ViewableMap, x, y, r int, walls bool) ViewMap {
	xo := 0
	yo := 0
	xmin := 0
	ymin := 0
	xmax := fov.Width()
	ymax := fov.Height()
	r2 := r * r
	if r > 0 {
		xmin = max(0, x-r)
		ymin = max(0, y-r)
		xmax = min(fov.Width(), x+r+1)
		ymax = min(fov.Height(), y+r+1)
	}
	xo = xmin
	yo = ymin

	data := make(ViewMap, 0)

	for xo < xmax {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		xo++
	}
	xo = xmax - 1
	yo = ymin + 1
	for yo < ymax {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		yo++
	}
	xo = xmax - 2
	yo = ymax - 1
	for xo >= 0 {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		xo--
	}
	xo = xmin
	yo = ymax - 2
	for yo > 0 {
		data.Update(fovCircularCastRay(fov, x, y, xo, yo, r2, walls))
		yo--
	}
	// if walls {
	// 	fovCircularPostProc(fov, data, xmin, ymin, x, y, -1, -1)
	// 	fovCircularPostProc(fov, data, x, ymin, xmax-1, y, 1, -1)
	// 	fovCircularPostProc(fov, data, xmin, y, x, ymax-1, -1, 1)
	// 	fovCircularPostProc(fov, data, x, y, xmax-1, ymax-1, 1, 1)
	// }

	return data
}
