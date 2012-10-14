package rog

import (
	"image"
	"image/draw"
	"io"
	"os"
)

type FontData struct {
	Image                 image.Image
	Width, Height         int
	CellWidth, CellHeight int
	mapping               map[rune]int
}

func (fd *FontData) Map(ch rune) (int, bool) {
	if fd.mapping == nil {
		return int(ch), true
	}
	position, ok := fd.mapping[ch]
	return position, ok
}

func ReadFont(r io.Reader, cellWidth, cellHeight int, maps string) *FontData {
	m, _, err := image.Decode(r)
	if err != nil {
		panic(err)
	}

	b := m.Bounds()
	newm := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newm, newm.Bounds(), m, b.Min, draw.Src)

	var mapping map[rune]int
	if len(maps) > 0 {
		mapping = make(map[rune]int)
		i := 0
		for _, v := range maps {
			mapping[v] = i
			i++
		}
	}

	width := m.Bounds().Max.X / cellWidth
	height := m.Bounds().Max.Y / cellHeight

	return &FontData{newm, width, height, cellWidth, cellHeight, mapping}
}

func Font(path string, cellWidth, cellHeight int, maps string) *FontData {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	return ReadFont(file, cellWidth, cellHeight, maps)
}
