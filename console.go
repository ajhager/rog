package rog

import (
	"fmt"
	"image/color"
)

type Console struct {
	bg, fg [][]color.Color
	ch     [][]rune
	w, h   int
	dBg, dFg color.Color
}

func NewConsole(width, height int) *Console {
	bg := make([][]color.Color, height)
	fg := make([][]color.Color, height)
	ch := make([][]rune, height)

	for y := 0; y < height; y++ {
		bg[y] = make([]color.Color, width)
		fg[y] = make([]color.Color, width)
		ch[y] = make([]rune, width)
	}

	con := &Console{bg, fg, ch, width, height, color.Black, color.White}
	con.Clear()
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
	if ch > 0 { con.ch[y][x] = ch }
	if fg != nil { con.fg[y][x] = fg }
	if bg != nil { con.bg[y][x] = bg }
}

func (con *Console) Put(x, y int, ch rune) {
	con.ch[y][x] = ch
	con.fg[y][x] = con.dFg
	con.bg[y][x] = con.dBg
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
