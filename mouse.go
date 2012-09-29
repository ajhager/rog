package rog

import (
	"image"
)

type MouseButton struct {
	Pressed, Released bool
}

type MouseData struct {
	Pos, DPos, Cell, DCell image.Point
	Left, Right, Middle    MouseButton
}
