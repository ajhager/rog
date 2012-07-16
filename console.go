package rog

import (
	"fmt"
	"image/color"
)

func colorEq(x, y color.Color) bool {
	xr, xg, xb, xa := x.RGBA()
	yr, yg, yb, ya := y.RGBA()
	return xr == yr && xg == yg && xb == yb && xa == ya
}

type Console struct {
	bg, fg [][]color.Color
	ch     [][]rune
	dirt   [][]bool
	w, h   int
	dBg, dFg color.Color
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

	con := &Console{bg, fg, ch, dirt, width, height, color.Black, color.White}
	
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

func (con *Console) SetDefaults(fg, bg color.Color) {
	if fg != nil { con.dFg = fg }
	if bg != nil { con.dBg = bg }
}

func (con *Console) Clear() {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.Set(x, y, ' ', con.dFg, con.dBg)
		}
	}
}

func (con *Console) Fill(ch rune, fg, bg color.Color) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.Set(x, y, ch, fg, bg)
		}
	}
}

func (con *Console) Set(x, y int, ch rune, fg, bg color.Color) {
	if ch > 0 && con.ch[y][x] != ch {
		con.ch[y][x] = ch
		con.dirt[y][x] = true
	}

	if fg != nil && !colorEq(fg, con.fg[y][x]) {
		con.fg[y][x] = fg
		con.dirt[y][x] = true
	}

	if bg != nil && !colorEq(bg, con.bg[y][x]) {
		con.bg[y][x] = bg
		con.dirt[y][x] = true
	}
}

func (con *Console) Put(x, y int, ch rune) {
	con.Set(x, y, ch, con.dFg, con.dBg)
}

func (con *Console) Print(s string, rest...interface{}) {
	runes := []rune(fmt.Sprintf(s, rest...))
	for x := 0; x < len(runes); x++ {
		con.Set(x, 5, runes[x], nil, nil)
	}
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
