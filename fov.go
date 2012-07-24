package rog

// FOVAlgo takes a FOVMap x,y vantage, radius of the view, whether to include walls and then marks in the map which cells are viewable.
type FOVAlgo func(*FOVMap, int, int, int, bool)

type FOVMap struct {
	w, h              int
	blocked, viewable [][]bool
}

func NewFOVMap(width, height int) *FOVMap {
	blocked := make([][]bool, height)
	viewable := make([][]bool, height)

	for y := 0; y < height; y++ {
		blocked[y] = make([]bool, width)
		viewable[y] = make([]bool, width)
	}

	return &FOVMap{width, height, blocked, viewable}
}

// In returns whether the coordinate is inside the map bounds.
func (this *FOVMap) In(x, y int) bool {
	return x >= 0 && x < this.w && y >= 0 && y < this.h
}

// Update runs the give fov alogrithm on the map.
func (this *FOVMap) Update(x, y, radius int, includeWalls bool, algo FOVAlgo) {
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			this.viewable[y][x] = false
		}
	}
	algo(this, x, y, radius, includeWalls)
}

// Block sets a cell as blocking or not.
func (this *FOVMap) Block(x, y int, blocked bool) {
	this.blocked[y][x] = blocked
}

// Look indicates if the cell at the coordinate can be seen.
func (this *FOVMap) Look(x, y int) bool {
	return this.viewable[y][x]
}

// Clear resets the map to completely unblocked but unviewable.
func (this *FOVMap) Clear() {
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			this.blocked[y][y] = false
			this.viewable[y][x] = false
		}
	}
}

// Width returns the width in cells of the map.
func (this *FOVMap) Width() int {
	return this.w
}

// Height returns the height in cells of the map.
func (this *FOVMap) Height() int {
	return this.h
}

func max(x, y int) (out int) {
	if x > y {
		out = x
	} else {
		out = y
	}
	return
}

func min(x, y int) (out int) {
	if x < y {
		out = x
	} else {
		out = y
	}
	return
}

// Circular Raycasting
func fovCircularCastRay(fov *FOVMap, xo, yo, xd, yd, r2 int, walls bool) {
	curx := xo
	cury := yo
	in := false
	blocked := false
	if fov.In(curx, cury) {
		in = true
		fov.viewable[cury][curx] = true
	}
	for _, p := range NewLine(xo, yo, xd, yd) {
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
			if !blocked && fov.blocked[cury][curx] {
				blocked = true
			} else if blocked {
				break
			}
			if walls || !blocked {
				fov.viewable[cury][curx] = true
			}
		} else if in {
			break
		}
	}
}

func fovCircularPostProc(fov *FOVMap, x0, y0, x1, y1, dx, dy int) {
	for cx := x0; cx <= x1; cx++ {
		for cy := y0; cy <= y1; cy++ {
			x2 := cx + dx
			y2 := cy + dy
			if fov.In(cx, cy) && fov.Look(cx, cy) && !fov.blocked[cy][cx] {
				if x2 >= x0 && x2 <= x1 {
					if fov.In(x2, cy) && fov.blocked[cy][x2] {
						fov.viewable[cy][x2] = true
					}
				}
				if y2 >= y0 && y2 <= y1 {
					if fov.In(cx, y2) && fov.blocked[y2][cx] {
						fov.viewable[y2][cx] = true
					}
				}
				if x2 >= x0 && x2 <= x1 && y2 >= y0 && y2 <= y1 {
					if fov.In(x2, y2) && fov.blocked[y2][x2] {
						fov.viewable[y2][x2] = true
					}
				}
			}
		}
	}
}

// FOVCicular raycasts out from the vantage in a circle.
func FOVCircular(fov *FOVMap, x, y, r int, walls bool) {
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
	for xo < xmax {
		fovCircularCastRay(fov, x, y, xo, yo, r2, walls)
		xo++
	}
	xo = xmax - 1
	yo = ymin + 1
	for yo < ymax {
		fovCircularCastRay(fov, x, y, xo, yo, r2, walls)
		yo++
	}
	xo = xmax - 2
	yo = ymax - 1
	for xo >= 0 {
		fovCircularCastRay(fov, x, y, xo, yo, r2, walls)
		xo--
	}
	xo = xmin
	yo = ymax - 2
	for yo > 0 {
		fovCircularCastRay(fov, x, y, xo, yo, r2, walls)
		yo--
	}
	if walls {
		fovCircularPostProc(fov, xmin, ymin, x, y, -1, -1)
		fovCircularPostProc(fov, x, ymin, xmax-1, y, 1, -1)
		fovCircularPostProc(fov, xmin, y, x, ymax-1, -1, 1)
		fovCircularPostProc(fov, x, y, xmax-1, ymax-1, 1, 1)
	}
}
