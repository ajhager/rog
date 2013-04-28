package main

import (
	"github.com/ajhager/rog"
	"image"
)

type Level struct {
	data   [][]rune
	player image.Point
}

func (l *Level) MoveBlocked(x, y int) bool {
	return l.data[y][x] == '#'
}

var (
	level = &Level{[][]rune{
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("################    ####################"),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("    #           ####        #  #  #     "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("##############              #######     "),
		[]rune("#                                       "),
		[]rune("#            #                          "),
		[]rune("#            #                          "),
		[]rune("################ ## ## ## ## ## ## ## ##"),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
		[]rune("                                        "),
	}, image.Pt(23, 9)}

	path []image.Point
)

func main() {
	rog.Open(40, 20, 1, false, "pathfinding", nil)
	for rog.Running() {
		if rog.Mouse().Left.Released {
			limit := image.Rect(0, 0, 40, 20)
			target := image.Pt(rog.Mouse().Cell.X, rog.Mouse().Cell.Y)
			path = rog.Path(level, limit, level.player, target)
		}

		for y := 0; y < 20; y++ {
			for x := 0; x < 40; x++ {
				rog.Set(x, y, nil, nil, string(level.data[y][x]))
			}
		}

		rog.Set(level.player.X, level.player.Y, nil, nil, "@")

		for _, p := range path {
			if p.X != level.player.X || p.Y != level.player.Y {
				rog.Set(p.X, p.Y, rog.Red, nil, "*")
			}
		}

		if rog.Key() == rog.Esc {
			rog.Close()
		}
		rog.Flush()
	}
}
