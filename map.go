package rog

type Map struct {
	w, h          int
	blocked, seen [][]bool
}

func NewMap(width, height int) *Map {
	blocked := make([][]bool, height)
	seen := make([][]bool, height)

	for y := 0; y < height; y++ {
		blocked[y] = make([]bool, width)
		seen[y] = make([]bool, width)
	}

	return &Map{width, height, blocked, seen}
}

// In returns whether the coordinate is inside the map bounds.
func (this *Map) In(x, y int) bool {
	return x >= 0 && x < this.w && y >= 0 && y < this.h
}

// Update runs the give fov alogrithm on the map.
func (this *Map) Fov(x, y, radius int, includeWalls bool, algo FOVAlgo) {
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			this.seen[y][x] = false
		}
	}
	algo(this, x, y, radius, includeWalls)
}

// Block sets a cell as blocking or not.
func (this *Map) Block(x, y int, blocked bool) {
	this.blocked[y][x] = blocked
}

// Look indicates if the cell at the coordinate can be seen.
func (this *Map) Look(x, y int) bool {
	return this.seen[y][x]
}

// Clear resets the map to completely unblocked but unviewable.
func (this *Map) Clear() {
	for y := 0; y < this.h; y++ {
		for x := 0; x < this.w; x++ {
			this.blocked[y][y] = false
			this.seen[y][x] = false
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
