// Copyright 2014 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package rog

// Should gracefully handle out of bounds coordinates.
type Grid2d interface {
	Get(x, y int) bool
	Set(x, y int, on bool)
	Bounds() (int, int)
}

// Dense
type DenseGrid struct {
	width  int
	height int
	data   [][]bool
}

// Sparse
type SparseGrid struct {
	width  int
	height int
	data   map[int]bool
}

func NewSparseGrid(width, height int) *SparseGrid {
	return &SparseGrid{width, height, make(map[int]bool)}
}

func (m *SparseGrid) index(x, y int) int {
	y = y % m.height
	if y < 0 {
		y = m.height - y
	}
	x = x % m.width
	if x < 0 {
		x = m.width - x
	}
	return y*m.width + x
}

func (m *SparseGrid) Get(x, y int) bool {
	open, _ := m.data[m.index(x, y)]
	return open
}

func (m *SparseGrid) Set(x, y int, open bool) {
	if open {
		m.data[m.index(x, y)] = open
	} else {
		delete(m.data, m.index(x, y))
	}
}

func (m *SparseGrid) Bounds() (int, int) {
	return m.width, m.height
}

func DigArena(grid Grid2d, x, y, width, height int) {
	for i := x; i < x+width; i++ {
		for j := y; j < y+height; j++ {
			if i != x && j != y && i != x+width-1 && j != y+height-1 {
				grid.Set(i, j, true)
			} else {
				grid.Set(i, j, false)
			}
		}
	}
}
