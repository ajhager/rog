package rog

import (
	"fmt"
	"image/color"
)

type Console struct {
	bg, fg   [][]color.Color
	ch       [][]rune
	w, h     int
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

	for x := 0; x < con.w; x++ {
		for y := 0; y < con.h; y++ {
			con.ch[y][x] = ' '
			con.fg[y][x] = color.White
			con.bg[y][x] = color.Black
		}
	}

	return con
}

func (con *Console) Fill(x0, y0, x1, y1 int, ch rune, fg, bg color.Color, bl ...ColorBlend) {
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			con.Set(x, y, string(ch), fg, bg, bl...)
		}
	}
}

func (con *Console) put(x, y int, ch rune, fg, bg color.Color, blend ColorBlend) {
	if ch > 0 {
		con.ch[y][x] = ch
	}

	if fg != nil {
		con.fg[y][x] = fg
	}

	if bg != nil {
		con.bg[y][x] = blend(bg, con.bg[y][x])
    }

}

func (con *Console) Set(x, y int, data string, fg, bg color.Color, bl ...ColorBlend) {
    blend := Normal
    if len(bl) > 0 {
        blend = bl[0]
    }

    runes := []rune(data)
    num := len(runes)
    if num > 0 {
	    for xx := 0; xx < num; xx++ {
		    con.put(xx+x, y, runes[xx], fg, bg, blend)
	    }
    } else {
		con.put(x, y, -1, fg, bg, blend)
    }
}

func (con *Console) P(s string, rest ...interface{}) string {
	return fmt.Sprintf(s, rest...)
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
