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

func conv(r, g, b, a uint32) (rr, gg, bb, aa float64) {
	const M = float64(1<<16 - 1)
	rr = float64(r) / M
	gg = float64(g) / M
	bb = float64(b) / M
	aa = float64(a) / M
	return
}

func colorMultiply(x, y color.Color) color.Color {
	xr, xg, xb, xa := conv(x.RGBA())
	yr, yg, yb, ya := conv(y.RGBA())
	// println(uint8(float64(xg) * float64(yg) / float64(255)))
	// println(uint8(xg * yg / uint32(255)))
	// println(float64(xr) / float64(255))
	// println(float64(xr) / float64(math.MaxUint32))
	// fmt.Printf("%v * %v / 255 = %v (%v)\n", uint8(xr), uint8(yr), uint8(xr)*uint8(yr)/uint8(255), uint8(xr*yr/255)):w
	return color.RGBA{
		uint8(xr * yr * 255),
		uint8(xg * yg * 255),
		uint8(xb * yb * 255),
		uint8(xa * ya * 255),
	}
}

func colorScreen(x, y color.Color) color.Color {
	xr, xg, xb, xa := conv(x.RGBA())
	yr, yg, yb, ya := conv(y.RGBA())
	return color.RGBA{
		uint8(255 * (1 - ((1 - xr) * (1 - yr)))),
		uint8(255 * (1 - ((1 - xg) * (1 - yg)))),
		uint8(255 * (1 - ((1 - xb) * (1 - yb)))),
		uint8(255 * (1 - ((1 - xa) * (1 - ya)))),
	}
}

func ol(top, bot float64) (out uint8) {
	if bot < 0.5 {
		out = uint8(2 * top * bot * 255)
	} else {
		out = uint8(255 * (1 - 2*(1-top)*(1-bot)))
	}
	return
}

func colorOverlay(x, y color.Color) color.Color {
	xr, xg, xb, xa := conv(x.RGBA())
	yr, yg, yb, ya := conv(y.RGBA())
	return color.RGBA{ol(xr, yr), ol(xg, yg), ol(xb, yb), ol(xa, ya)}
}

type BgFlag uint8

const (
	None BgFlag = iota
	Set
	Multiply
	Screen
	Overlay

	Lighten
	Darken
	Dodge
	CBurn
	Add
	AddA
	Burn
	Alpha
	Default
)

type Console struct {
	bg, fg   [][]color.Color
	ch       [][]rune
	dirt     [][]bool
	w, h     int
	dBg, dFg color.Color
	bgFlag   BgFlag
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

	con := &Console{bg, fg, ch, dirt, width, height, color.Black, color.White, None}

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

func (con *Console) SetDefaults(fg, bg color.Color, flag BgFlag) {
	if fg != nil {
		con.dFg = fg
	}
	if bg != nil {
		con.dBg = bg
	}
	con.bgFlag = flag
}

func (con *Console) Clear() {
	con.Fill(' ', con.dFg, con.dBg)
}

func (con *Console) Fill(ch rune, fg, bg color.Color) {
	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.Set(x, y, ch, fg, bg, Set)
		}
	}
}

func (con *Console) Set(x, y int, ch rune, fg, bg color.Color, flag BgFlag) {
	if ch > 0 && con.ch[y][x] != ch {
		con.ch[y][x] = ch
		con.dirt[y][x] = true
	}

	if fg != nil && !colorEq(fg, con.fg[y][x]) {
		con.fg[y][x] = fg
		con.dirt[y][x] = true
	}

	if bg != nil && !colorEq(bg, con.bg[y][x]) {
		con.dirt[y][x] = true
		switch flag {
		case None:
		case Set:
			con.bg[y][x] = bg
		case Multiply:
			con.bg[y][x] = colorMultiply(bg, con.bg[y][x])
		case Screen:
			con.bg[y][x] = colorScreen(bg, con.bg[y][x])
		case Overlay:
			con.bg[y][x] = colorOverlay(bg, con.bg[y][x])
		}
	}
}

func (con *Console) Put(x, y int, ch rune) {
	con.Set(x, y, ch, con.dFg, con.dBg, con.bgFlag)
}

func (con *Console) Print(s string, rest ...interface{}) {
	runes := []rune(fmt.Sprintf(s, rest...))
	for x := 0; x < len(runes); x++ {
		con.Set(x, 5, runes[x], nil, nil, Set)
	}
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
