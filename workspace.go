package rog

import (
	"image/draw"
)

type Workspace interface {
	Open(int, int)
	IsOpen() bool
	Close()
	Name(string)
	Render(*Console)
	Screen() draw.Image
	Mouse() *MouseData
	Key() string
}
