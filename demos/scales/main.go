package main

import (
	"github.com/ajhager/rog"
	_ "github.com/ajhager/rog/wde"
)

var rainbow = rog.Discrete(rog.Red, rog.Orange, rog.Yellow, rog.Green, rog.Blue, rog.Purple, rog.Magenta)

func main() {
	rog.Open(21, 11, 2, "rog")
    rog.Set(3, 5, rog.Burn(rog.Grey), rainbow, "Discrete Scale!")
	for rog.IsOpen() {
		if rog.Key() == rog.Escape {
			rog.Close()
		}
		rog.Flush()
	}
}
