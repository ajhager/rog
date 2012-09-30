package rog

import (
	"fmt"
)

// Console is a double buffered grid of unicode characters that can be rendered to an image.Image.
type Console struct {
	bg, fg [][]RGB
	ch     [][]rune
	w, h   int
}

// NewConsole creates an empty console.
func NewConsole(width, height int) *Console {
	bg := make([][]RGB, height)
	fg := make([][]RGB, height)
	ch := make([][]rune, height)

	for y := 0; y < height; y++ {
		bg[y] = make([]RGB, width)
		fg[y] = make([]RGB, width)
		ch[y] = make([]rune, width)
	}

	con := &Console{bg, fg, ch, width, height}

	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.bg[y][x] = Black
			con.fg[y][x] = White
			con.ch[y][x] = ' '
		}
	}

	return con
}

func (con *Console) put(x, y int, fg, bg Blender, ch rune) {
	if ch > 0 {
		con.ch[y][x] = ch
	}

	if bg != nil {
		con.bg[y][x] = bg.Blend(con.bg[y][x])
	}

	if fg != nil {
		con.fg[y][x] = fg.Blend(con.bg[y][x])
	}
}

func (con *Console) set(i, j, x, y, w, h int, fg, bg Blender, data string, rest ...interface{}) {
	runes := []rune(fmt.Sprintf(data, rest...))
	if len(runes) > 0 {
		if h == 0 {
			h = con.h - y
		}
		for k := 0; k < len(runes); k++ {
			if i == x+w {
				j += 1
				i = x
			}
			if j == y+h {
				break
			}
			con.put(i, j, fg, bg, runes[k])
			i += 1
		}
	} else {
		con.put(i, j, fg, bg, -1)
	}
}

// Set draws a string starting at x,y onto the console, wrapping at the bounds if needed.
func (con *Console) Set(x, y int, fg, bg Blender, data string, rest ...interface{}) {
	con.set(x, y, 0, 0, con.w, con.h, fg, bg, data, rest...)
}

// SetR draws a string starting at x,y onto the console, wrapping at the bounds created by x, y, w, h if needed.
// If h is 0, the text will cut off at the bottom of the console, otherwise it will cut off after the y+h row.
func (con *Console) SetR(x, y, w, h int, fg, bg Blender, data string, rest ...interface{}) {
	con.set(x, y, x, y, w, h, fg, bg, data, rest...)
}

// Get returns the fg, bg colors and rune of the cell.
func (con *Console) Get(x, y int) (RGB, RGB, rune) {
	return con.fg[y][x], con.bg[y][x], con.ch[y][x]
}

// Fill draws a rect on the root console using ch.
func (con *Console) Fill(x, y, w, h int, fg, bg Blender, ch rune) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			con.Set(i, j, fg, bg, string(ch))
		}
	}
}

// Clear is a short hand to fill the entire screen with the given colors and rune.
func (con *Console) Clear(fg, bg Blender, ch rune) {
	con.Fill(0, 0, con.w, con.h, fg, bg, ch)
}

// Width returns the width of the console in cells.
func (con *Console) Width() int {
	return con.w
}

// Height returns the height of the console in cells.
func (con *Console) Height() int {
	return con.h
}
