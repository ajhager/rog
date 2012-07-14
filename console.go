package rog

import (
	"image/color"
)

type Console struct {
	bg, fg [][]color.Color
	ch     [][]rune
	w, h   int
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

	con := &Console{bg, fg, ch, width, height}
	con.Clear()
	return con
}

func (con *Console) Clear() {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.Set(x, y, 0, color.RGBA{0, 0, 0, 0}, color.RGBA{255, 255, 255, 255})
		}
	}
}

func (con *Console) Fill(ch rune, bg, fg color.Color) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.Set(x, y, ch, bg, fg)
		}
	}
}

func (con *Console) FillCh(ch rune) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.SetCh(x, y, ch)
		}
	}
}

func (con *Console) FillBg(bg color.Color) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.SetBg(x, y, bg)
		}
	}
}

func (con *Console) FillFg(fg color.Color) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.SetFg(x, y, fg)
		}
	}
}

func (con *Console) Set(x, y int, ch rune, bg, fg color.Color) {
	con.bg[y][x] = bg
	con.fg[y][x] = fg
	con.ch[y][x] = ch
}

func (con *Console) SetCh(x, y int, ch rune) {
	con.ch[y][x] = ch
}

func (con *Console) SetBg(x, y int, bg color.Color) {
	con.bg[y][x] = bg
}

func (con *Console) SetFg(x, y int, fg color.Color) {
	con.fg[y][x] = fg
}

func (con *Console) SetChFg(x, y int, ch rune, fg color.Color) {
	con.ch[y][x] = ch
	con.fg[y][x] = fg
}

func (con *Console) SetChBg(x, y int, ch rune, bg color.Color) {
	con.ch[y][x] = ch
	con.bg[y][x] = bg
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
