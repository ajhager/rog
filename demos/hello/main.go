package main

import (
	"github.com/ajhager/rog"
)

func main() {
	rog.Open(20, 11, 1, "rog", nil)
	for rog.Running() {
		rog.Set(5, 5, nil, nil, "Hello, 世界!")
		if rog.Key() == rog.Escape {
			rog.Close()
		}
		rog.Flush()
	}
}
