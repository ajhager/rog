package rog

import (
	"fmt"
	"image/color"
)

type Console struct {
	bg, fg   [][]color.Color
	ch       [][]rune
	dirt     [][]bool
	w, h     int
	dBg, dFg color.Color
	blend    ColorBlend
}

func NewConsole(width, height int) *Console {
	bg := make([][]color.Color, height)
	fg := make([][]color.Color, height)
	ch := make([][]rune, height)
	dirt := make([][]bool, height)

	for y := 0; y < height; y++ {
		bg[y] = make([]color.Color, width)
		fg[y] = make([]color.Color, width)
		ch[y] = make([]rune, width)
		dirt[y] = make([]bool, width)
	}

	con := &Console{bg, fg, ch, dirt, width, height, color.Black, color.White, Normal}

	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.ch[y][x] = ' '
			con.fg[y][x] = con.dFg
			con.bg[y][x] = con.dBg
			con.dirt[y][x] = true
		}
	}

	return con
}

func (con *Console) SetDefaults(fg, bg color.Color, blend ColorBlend) {
	if fg != nil {
		con.dFg = fg
	}
	if bg != nil {
		con.dBg = bg
	}
	con.blend = blend
}

func (con *Console) Dirty() {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.dirt[y][x] = true
		}
	}
}

func (con *Console) Clear() {
	con.Fill(' ', con.dFg, con.dBg)
}

func (con *Console) Fill(ch rune, fg, bg color.Color) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.Set(x, y, ch, fg, bg, Normal)
		}
	}
}

func (con *Console) Set(x, y int, ch rune, fg, bg color.Color, blend ColorBlend) {
	if ch > 0 && con.ch[y][x] != ch {
		con.dirt[y][x] = true
		con.ch[y][x] = ch
	}

	if fg != nil && !colorEq(fg, con.fg[y][x]) {
		con.dirt[y][x] = true
		con.fg[y][x] = fg
	}

	if bg != nil && !colorEq(bg, con.bg[y][x]) {
		con.dirt[y][x] = true
		con.bg[y][x] = blend(bg, con.bg[y][x])
	}
}

func (con *Console) Put(x, y int, ch rune) {
	con.Set(x, y, ch, con.dFg, con.dBg, con.blend)
}

func (con *Console) Print(s string, rest ...interface{}) {
	runes := []rune(fmt.Sprintf(s, rest...))
	for x := 0; x < len(runes); x++ {
		con.Set(x, 0, runes[x], nil, nil, Normal)
	}
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
