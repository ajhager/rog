package rog

import (
	"image"
)

type MouseButton struct {
	Pressed, Released bool
}

type MouseData struct {
	Pos, Cell image.Point
	Left, Right, Middle    MouseButton
}
