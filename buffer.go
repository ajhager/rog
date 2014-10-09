package rog

import (
	"strings"

	"github.com/ajhager/engi"
)

type Buffer struct {
	bg, fg [][]uint32
	ch     [][]rune
	w, h   int
}

func NewBuffer(width, height int) *Buffer {
	bg := make([][]uint32, height)
	fg := make([][]uint32, height)
	ch := make([][]rune, height)

	for y := 0; y < height; y++ {
		bg[y] = make([]uint32, width)
		fg[y] = make([]uint32, width)
		ch[y] = make([]rune, width)
	}

	buf := &Buffer{bg, fg, ch, width, height}

	for x := 0; x < buf.w; x++ {
		for y := 0; y < buf.h; y++ {
			buf.bg[y][x] = 0x000000
			buf.fg[y][x] = 0xffffff
			buf.ch[y][x] = ' '
		}
	}

	return buf
}

func (buf *Buffer) put(x, y int, fg, bg uint32, ch rune) {
	if x < 0 || x >= buf.w || y < 0 || y >= buf.h {
		return
	}

	if ch > 0 {
		buf.ch[y][x] = ch
	}

	buf.bg[y][x] = bg
	buf.fg[y][x] = fg
}

func (buf *Buffer) Set(x, y int, fg, bg uint32, data string) {
	t := strings.Count(data, "") - 1
	if t > 0 {
		for _, r := range data {
			buf.put(x, y, fg, bg, r)
			x += 1
		}
	} else {
		buf.put(x, y, fg, bg, -1)
	}
}

// Blit draws buf onto this buffer with top left starting at x, y.
func (buf *Buffer) Blit(o *Buffer, x, y int) {
	for i := 0; i < o.Width(); i++ {
		for j := 0; j < o.Height(); j++ {
			fg, bg, ch := o.Get(i, j)
			buf.Set(x+i, y+j, fg, bg, string(ch))
		}
	}
}

// Clear is a short hand to fill the entire screen with the given colors and rune.
func (buf *Buffer) Clear(fg, bg uint32, ch rune) {
	buf.Fill(0, 0, buf.w, buf.h, fg, bg, ch)
}

// Fill draws a rect on the root buffer using ch.
func (buf *Buffer) Fill(x, y, w, h int, fg, bg uint32, ch rune) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			buf.Set(i, j, fg, bg, string(ch))
		}
	}
}

// Get returns the fg, bg colors and rune of the cell.
func (buf *Buffer) Get(x, y int) (uint32, uint32, rune) {
	return buf.fg[y][x], buf.bg[y][x], buf.ch[y][x]
}

// Width returns the width of the buffer in cells.
func (buf *Buffer) Width() int {
	return buf.w
}

// Height returns the height of the buffer in cells.
func (buf *Buffer) Height() int {
	return buf.h
}

func (buf *Buffer) Draw(batch *engi.Batch, font *engi.Font, cellW, cellH int) {
	for x := 0; x < buf.w; x++ {
		for y := 0; y < buf.h; y++ {
			fg, bg, ch := buf.Get(x, y)
			fx := float32(x * cellW)
			fy := float32(y * cellH)
			font.Put(batch, 0, fx, fy, bg)
			font.Put(batch, ch, fx, fy, fg)
		}
	}
}
