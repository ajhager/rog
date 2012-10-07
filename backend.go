package rog

type Backend interface {
	Open(int, int, int, *FontData)
	Close()
	Running() bool
	Name(string)
	Render(*Console)
	Mouse() *MouseData
	Key() int
}
