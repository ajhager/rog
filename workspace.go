package rog

import (
	"image"
)

type Workspace interface {
	Open(int, int)
	IsOpen() bool
	Close()
	Name(string)
	Render(*Console)
	Screen() image.Image
	Mouse() *MouseData
	Key() string
}
