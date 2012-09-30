package rog

import (
	"image"
)

type Backend interface {
	Open(int, int, int)
	IsOpen() bool
	Close()
	Name(string)
	Render(*Console)
	Screen() image.Image
	Mouse() *MouseData
	Key() string
}
