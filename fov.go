package rog

// FOVAlgo takes a FOVMap x,y vantage, radius of the view, whether to include walls and then marks in the map which cells are viewable.
type FOVAlgo func(*Map, int, int, int, bool)

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
func fovCircularCastRay(fov *Map, xo, yo, xd, yd, r2 int, walls bool) {
	curx := xo
	cury := yo
	in := false
	blocked := false
	if fov.In(curx, cury) {
		in = true
		fov.seen[cury][curx] = true
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
			if !blocked && fov.blocked[cury][curx] {
				blocked = true
			} else if blocked {
				break
			}
			if walls || !blocked {
				fov.seen[cury][curx] = true
			}
		} else if in {
			break
		}
	}
}

func fovCircularPostProc(fov *Map, x0, y0, x1, y1, dx, dy int) {
	for cx := x0; cx <= x1; cx++ {
		for cy := y0; cy <= y1; cy++ {
			x2 := cx + dx
			y2 := cy + dy
			if fov.In(cx, cy) && fov.Look(cx, cy) && !fov.blocked[cy][cx] {
				if x2 >= x0 && x2 <= x1 {
					if fov.In(x2, cy) && fov.blocked[cy][x2] {
						fov.seen[cy][x2] = true
					}
				}
				if y2 >= y0 && y2 <= y1 {
					if fov.In(cx, y2) && fov.blocked[y2][cx] {
						fov.seen[y2][cx] = true
					}
				}
				if x2 >= x0 && x2 <= x1 && y2 >= y0 && y2 <= y1 {
					if fov.In(x2, y2) && fov.blocked[y2][x2] {
						fov.seen[y2][x2] = true
					}
				}
			}
		}
	}
}

// FOVCicular raycasts out from the vantage in a circle.
func FOVCircular(fov *Map, x, y, r int, walls bool) {
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
