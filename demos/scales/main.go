package main

import (
	"hagerbot.com/rog"
)

var drainbow = rog.Discrete(rog.Red, rog.Orange, rog.Yellow, rog.Green, rog.Blue, rog.Purple, rog.Magenta)
var lrainbow = rog.Linear(rog.Red, rog.Orange, rog.Yellow, rog.Green, rog.Blue, rog.Purple, rog.Magenta)

func main() {
	rog.Open(20, 10, 2, false, "scales", nil)
	rog.Set(0, 3, rog.Burn(rog.Grey), drainbow, "   Discrete Scale   ")
	rog.Set(0, 6, rog.Burn(rog.Grey), lrainbow, "    Linear Scale    ")
	for rog.Running() {
		rog.Set(0, 9, nil, nil, "%v", rog.Fps())
		if rog.Key() == rog.Esc {
			rog.Close()
		}
		rog.Flush()
	}
}
