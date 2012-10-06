package main

import (
	"github.com/ajhager/rog"
	_ "github.com/ajhager/rog/glfw"
)

func main() {
	rog.Open(20, 11, 2, "rog", "../../data/font.png")
	for rog.IsOpen() {
		rog.Set(5, 5, nil, nil, "Hello, 世界!")
		if rog.Key() == rog.Escape {
			rog.Close()
		}
		rog.Flush()
	}
}
