package rog

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Console is a double buffered grid of unicode characters that can be rendered to an image.Image.
type Console struct {
	bg, bgbuf, fg, fgbuf [][]color.Color
	ch, chbuf            [][]rune
	w, h                 int
	font                 image.Image
}

// NewConsole creates an empty console.
func NewConsole(width, height int) *Console {
	bg := make([][]color.Color, height)
	bgbuf := make([][]color.Color, height)
	fg := make([][]color.Color, height)
	fgbuf := make([][]color.Color, height)
	ch := make([][]rune, height)
	chbuf := make([][]rune, height)

	for y := 0; y < height; y++ {
		bg[y] = make([]color.Color, width)
		bgbuf[y] = make([]color.Color, width)
		fg[y] = make([]color.Color, width)
		fgbuf[y] = make([]color.Color, width)
		ch[y] = make([]rune, width)
		chbuf[y] = make([]rune, width)
	}

	mask, _, err := image.Decode(bytes.NewBuffer(font()))
	if err != nil {
		panic(err)
	}

	con := &Console{bg, bgbuf, fg, fgbuf, ch, chbuf, width, height, mask}

	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.bg[y][x] = color.Black
			con.bgbuf[y][x] = color.Black
			con.fg[y][x] = color.White
			con.fgbuf[y][x] = color.White
			con.ch[y][x] = ' '
			con.chbuf[y][x] = ' '
		}
	}

	return con
}

func (con *Console) put(x, y int, fg, bg interface{}, ch rune) {
	if ch > 0 {
		con.ch[y][x] = ch
	}

	switch bgcolor := bg.(type) {
	case color.Color:
		con.bg[y][x] = bgcolor
	case Blender:
		con.bg[y][x] = bgcolor(con.bg[y][x])
	default:
	}

	switch fgcolor := fg.(type) {
	case color.Color:
		con.fg[y][x] = fgcolor
	case Blender:
		con.fg[y][x] = fgcolor(con.bg[y][x])
	default:
	}
}

func (con *Console) set(i, j, x, y, w, h int, fg, bg interface{}, data string, rest ...interface{}) {
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
func (con *Console) Set(x, y int, fg, bg interface{}, data string, rest ...interface{}) {
	con.set(x, y, 0, 0, con.w, con.h, fg, bg, data, rest...)
}

// Set draws a string starting at x,y onto the console, wrapping at the bounds created by x, y, w, h if needed.
// If h is 0, the text will cut off at the bottom of the console, otherwise it will cut off after the y+h row.
func (con *Console) SetR(x, y, w, h int, fg, bg interface{}, data string, rest ...interface{}) {
	con.set(x, y, x, y, w, h, fg, bg, data, rest...)
}

// Get returns the fg, bg colors and rune of the cell.
func (con *Console) Get(x, y int) (color.Color, color.Color, rune) {
	return con.fg[y][x], con.bg[y][x], con.ch[y][x]
}

// Fill draws a rect on the root console using ch.
func (con *Console) Fill(x, y, w, h int, fg, bg interface{}, ch rune) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			con.Set(i, j, fg, bg, string(ch))
		}
	}
}

// Clear is a short hand to fill the entire screen with the given colors and rune.
func (con *Console) Clear(fg, bg interface{}, ch rune) {
	con.Fill(0, 0, con.w, con.h, fg, bg, ch)
}

// Render draws the console onto an image.
func (c *Console) Render(im draw.Image) {
	maskRect := image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
	for y := 0; y < c.h; y++ {
		for x := 0; x < c.w; x++ {
			bg := c.bg[y][x]
			fg := c.fg[y][x]
			ch := c.ch[y][x]
			if bg != c.bgbuf[y][x] || fg != c.fgbuf[y][x] || ch != c.chbuf[y][x] {
				c.bgbuf[y][x] = bg
				c.fgbuf[y][x] = fg
				c.chbuf[y][x] = ch
				rect := maskRect.Add(image.Point{x * 16, y * 16})
				src := &image.Uniform{bg}
				draw.Draw(im, rect, src, image.ZP, draw.Src)

				if ch != ' ' {
					src = &image.Uniform{fg}
					draw.DrawMask(im, rect, src, image.ZP, c.font, image.Point{int(ch%32) * 16, int(ch/32) * 16}, draw.Over)
				}
			}
		}
	}
}

// Width returns the width of the console in cells.
func (con *Console) Width() int {
	return con.w
}

// Height returns the height of the console in cells.
func (con *Console) Height() int {
	return con.h
}
