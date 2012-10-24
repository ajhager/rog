package rog

import (
	"image"
)

type Entity interface {
	Renderable
	Viewable
	Walkable
}

type Map struct {
	w, h          int
	tiles map[Point]Entity
	viewable ViewMap
}

func NewMap(width, height int) *Map {
	return &Map{
		width, height, 
		make(map[Point]Entity, 0),
		make(ViewMap, 0),
	}
}

// In returns whether the coordinate is inside the map bounds.
func (this *Map) In(x, y int) bool {
	return x >= 0 && x < this.w && y >= 0 && y < this.h
}

func (this *Map) SetTile(e Entity, x, y int) {
	this.tiles[Point{x, y}]	= e
}

func (this *Map) GetTile(x, y int) (Entity, bool) {
	e, ok := this.tiles[Point{x, y}]
	return e, ok
}

func (this *Map) Viewable(x, y int) bool {
	tile, ok := this.tiles[Point{x, y}]
	switch ok {
		case true:
			return tile.IsViewable()
		case false:
			return false
	}; panic(nil)
}

func (this *Map) ViewMap() ViewMap {
	return this.viewable
}

func (this *Map) Roughness(x, y int) int {
	tile, ok := this.tiles[Point{x, y}]
	switch ok {
		case true:
			return tile.GetRoughness()
		case false:
			return PATH_MAX
	}; panic(nil)
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

func (this *Map) Look(x, y int) bool {
	return this.Viewable(x, y)
}

// Clear resets the map to completely unblocked but unviewable.
func (this *Map) Clear() {
	this.viewable = make(ViewMap, 0)
}

// Width returns the width in cells of the map.
func (this *Map) Width() int {
	return this.w
}

// Height returns the height in cells of the map.
func (this *Map) Height() int {
	return this.h
}
