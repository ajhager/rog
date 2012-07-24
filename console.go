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

func (con *Console) Fill(x0, y0, x1, y1 int, ch rune, fg, bg interface{}) {
	for x := x0; x < x1; x++ {
		for y := y0; y < y1; y++ {
			con.Set(x, y, string(ch), fg, bg)
		}
	}
}

func (con *Console) put(x, y int, ch rune, fg, bg interface{}) {
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

func (con *Console) Set(x, y int, data string, fg, bg interface{}) {
    runes := []rune(data)
    if len(runes) > 0 {
	    for xx := 0; xx < len(runes); xx++ {
		    con.put(xx+x, y, runes[xx], fg, bg)
	    }
    } else {
		con.put(x, y, -1, fg, bg)
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
