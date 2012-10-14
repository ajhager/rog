package main

import (
	"hagerbot.com/rog"
)

func main() {
	rog.Open(20, 11, 1, false, "rog", nil)
	for rog.Running() {
		rog.Set(5, 5, nil, nil, "Hello, 世界!")
		if rog.Key() == rog.Esc {
			rog.Close()
		}
		rog.Flush()
	}
}
