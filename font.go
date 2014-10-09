package rog

import "github.com/ajhager/engi"

type Font struct {
	native     *engi.Font
	path       string
	cellWidth  int
	cellHeight int
}

func NewFont(path string, cellWidth, cellHeight int) *Font {
	return &Font{nil, path, cellWidth, cellHeight}
}

func (font *Font) Remap(mapping string) {
	font.native.Remap(mapping)
}
