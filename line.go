package rog

import (
	"image"
	"math"
)

// Line returns the points on the grid that the line would pass through.
func Line(x0, y0, x1, y1 int) []image.Point {
	dx := int(math.Abs(float64(x1 - x0)))
	dy := int(math.Abs(float64(y1 - y0)))
	sx := 0
	sy := 0
	if x0 < x1 {
		sx = 1
	} else {
		sx = -1
	}

	if y0 < y1 {
		sy = 1
	} else {
		sy = -1
	}

	err := dx - dy

	ps := make([]image.Point, 0)
	for {
		ps = append(ps, image.Pt(x0, y0))
		if x0 == x1 && y0 == y1 {
			break
		}
		e2 := err * 2
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
	return ps
}
