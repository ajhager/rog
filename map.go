package rog

import (
	"image"
)


type Map struct {
	w, h          int
	blocked [][]bool
	viewable ViewMap
}

func NewMap(width, height int) *Map {
	blocked := make([][]bool, height)
	for y := 0; y < height; y++ {
		blocked[y] = make([]bool, width)		
	}

	seen := make(ViewMap, 0)

	return &Map{width, height, blocked, seen}
}

// In returns whether the coordinate is inside the map bounds.
func (this *Map) In(x, y int) bool {
	return x >= 0 && x < this.w && y >= 0 && y < this.h
}

func (this *Map) Viewable(x, y int) bool {
	return this.blocked[y][x]
}

func (this *Map) Roughness(x, y int) int {
	if this.blocked[y][x] {
		return PATH_MAX
	}
	return PATH_MIN
}

func (this *Map) ViewMap() ViewMap {
	data := make(ViewMap, 0)
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			data[Point{x, y}] = this.blocked[y][x]
		}
	}
	return data
}

// Update runs the give fov alogrithm on the map.
func (this *Map) Fov(x, y, radius int, includeWalls bool, algo FOVAlgo) {
	this.viewable = algo(this, x, y, radius, includeWalls)
}

// Performs astar and returns the list of cells on the path.
func (this *Map) Path(x0, y0, x1, y1 int) []image.Point {
	nodes := Astar(this, x0, y0, x1, y1, true)
	points := make([]image.Point, len(nodes))
	for i := 0; i < len(nodes); i++ {
		points[i] = image.Pt(nodes[i].X, nodes[i].Y)
	}
	return points
}

// Block sets a cell as blocking or not.
func (this *Map) Block(x, y int, blocked bool) {
	this.blocked[y][x] = blocked
}

// Look indicates if the cell at the coordinate can be seen.
func (this *Map) Look(x, y int) bool {
	return this.viewable[Point{x, y}]
}

// Clear resets the map to completely unblocked but unviewable.
func (this *Map) Clear() {
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			this.blocked[y][y] = false
			this.viewable = make(ViewMap, 0)
		}
	}
}

// Width returns the width in cells of the map.
func (this *Map) Width() int {
	return this.w
}

// Height returns the height in cells of the map.
func (this *Map) Height() int {
	return this.h
}
