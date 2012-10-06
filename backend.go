package rog

type Backend interface {
	Open(int, int, int, string)
	IsOpen() bool
	Close()
	Name(string)
	Render(*Console)
	Mouse() *MouseData
	Key() int
}
