package rog

type Backend interface {
	Open(int, int, int)
	IsOpen() bool
	Close()
	Name(string)
	Render(*Console)
	Mouse() *MouseData
	Key() string
}
