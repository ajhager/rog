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
            if j == y + h {
                break
            }
		    con.put(i, j, fg, bg, runes[k])
            i += 1
	    }
    } else {
		con.put(i, j, fg, bg, -1)
    }
}

func (con *Console) Set(x, y int, fg, bg interface{}, data string, rest ...interface{}) {
    con.set(x, y, 0, 0, con.w, con.h, fg, bg, data, rest...)
}

func (con *Console) SetR(x, y, w, h int, fg, bg interface{}, data string, rest ...interface{}) {
    con.set(x, y, x, y, w, h, fg, bg, data, rest...)
}

func (con *Console) Fill(x, y, w, h int, ch rune, fg, bg interface{}) {
	for i := x; i < x+w; i++ {
		for j := y; j < y+h; j++ {
			con.Set(i, j, fg, bg, string(ch))
		}
	}
}

func (con *Console) Width() int {
	return con.w
}

func (con *Console) Height() int {
	return con.h
}
