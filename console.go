package rog

import (
	"fmt"
	"image/color"
	"math"
)

func colorEq(x, y color.Color) bool {
	xr, xg, xb, xa := x.RGBA()
	yr, yg, yb, ya := y.RGBA()
	return xr == yr && xg == yg && xb == yb && xa == ya
}

func colorToFloats(c color.Color) (rr, gg, bb, aa float64) {
	const M = float64(1<<16 - 1)
	r, g, b, a := c.RGBA()
	rr = float64(r) / M
	gg = float64(g) / M
	bb = float64(b) / M
	aa = float64(a) / M
	return
}

func overlay(top, bot float64) (out uint8) {
	if bot < 0.5 {
		out = uint8(2 * top * bot * 255)
	} else {
		out = uint8(255 * (1 - 2*(1-top)*(1-bot)))
	}
	return
}

func dodge(top, bot float64) (out uint8) {
	if bot != 1 {
		out = uint8(255 * clamp(0, 1, top / (1 - bot)))
	} else {
		out = uint8(255)
	}
	return
}

func clamp(low, high, value float64) float64 {
	return math.Min(high, math.Max(low, value))
}

type ColorBlend func(color.Color, color.Color) color.Color

func Normal(top, bot color.Color) color.Color {
	return top
}

func Multiply(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(topR * botR * 255),
		uint8(topG * botG * 255),
		uint8(topB * botB * 255),
		uint8(topA * botA * 255),
	}
}

func Screen(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * (1 - ((1 - topR) * (1 - botR)))),
		uint8(255 * (1 - ((1 - topG) * (1 - botG)))),
		uint8(255 * (1 - ((1 - topB) * (1 - botB)))),
		uint8(255 * (1 - ((1 - topA) * (1 - botA)))),
	}
}

func Overlay(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		overlay(topR, botR),
		overlay(topG, botG),
		overlay(topB, botB),
		overlay(topA, botA),
	}
}

func Lighten(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * math.Max(topR, botR)),
		uint8(255 * math.Max(topG, botG)),
		uint8(255 * math.Max(topB, botB)),
		uint8(255 * math.Max(topA, botA)),
	}
}

func Darken(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * math.Min(topR, botR)),
		uint8(255 * math.Min(topG, botG)),
		uint8(255 * math.Min(topB, botB)),
		uint8(255 * math.Min(topA, botA)),
	}
}

func Burn(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * clamp(0, 1, botR + topR - 1)),
		uint8(255 * clamp(0, 1, botG + topG - 1)),
		uint8(255 * clamp(0, 1, botB + topB - 1)),
		uint8(255 * clamp(0, 1, botA + topA - 1)),
	}
}

func Dodge(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		dodge(topR, botR),
		dodge(topG, botG),
		dodge(topB, botB),
		dodge(topA, botA),
	}
}

func Add(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * clamp(0, 1, botR + topR)),
		uint8(255 * clamp(0, 1, botG + topG)),
		uint8(255 * clamp(0, 1, botB + topB)),
		uint8(255 * clamp(0, 1, botA + topA)),
	}
}

func AddAlpha(a float64) ColorBlend {
	return func(top, bot color.Color) color.Color {
		topR, topG, topB, topA := colorToFloats(top)
		botR, botG, botB, botA := colorToFloats(bot)
		return color.RGBA{
			uint8(255 * clamp(0, 1, botR * a + topR)),
			uint8(255 * clamp(0, 1, botG * a + topG)),
			uint8(255 * clamp(0, 1, botB * a + topB)),
			uint8(255 * clamp(0, 1, botA * a + topA)),
		}
	}
}

func Alpha(a float64) ColorBlend {
	return func(top, bot color.Color) color.Color {
		topR, topG, topB, topA := colorToFloats(top)
		botR, botG, botB, botA := colorToFloats(bot)
		return color.RGBA{
			uint8(255 * (botR + (topR - botR) * a)),
			uint8(255 * (botG + (topG - botG) * a)),
			uint8(255 * (botB + (topB - botB) * a)),
			uint8(255 * (botA + (topA - botA) * a)),
		}
	}
}

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
		con.ch[y][x] = ch
		con.dirt[y][x] = true
	}

	if fg != nil && !colorEq(fg, con.fg[y][x]) {
		con.fg[y][x] = fg
		con.dirt[y][x] = true
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
		con.Set(x, 5, runes[x], nil, nil, Normal)
	}
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
