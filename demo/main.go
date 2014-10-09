package main

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/ajhager/engi"
	"github.com/ajhager/rogood"
)

var (
	player  *Player
	monster *Entity
	actionQ *rog.Queue
	grid    rog.Grid2d
)

type Entity struct {
	id   string
	X, Y int
}

func NewEntity(x, y int) *Entity {
	id := make([]byte, 16)
	rand.Read(id)
	return &Entity{hex.EncodeToString(id), x, y}
}

func (e *Entity) Id() string {
	return e.id
}

func (e *Entity) Act() float64 {
	println("MONSTER")
	return 2
}

type Player struct {
	*Entity
	Thinking bool
}

func (player *Player) Act() float64 {
	println("PLAYER")
	player.Thinking = true
	return 1
}

type Console struct {
	*rog.Console
}

func OpenConsole(title string, width, height int, zoom int, path string, cellW, cellH int) {
	console := &Console{rog.New(width, height, path, cellW, cellH)}
	engi.Open(title, width*cellW*zoom, height*cellH*zoom, false, console)
}

func (con *Console) Update(dt float32) {
	con.Clear(magenta, black, ' ')

	width, height := grid.Bounds()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			if grid.Get(x, y) {
				con.Set(x, y, dgrey, black, ".")
			} else {
				con.Set(x, y, dblue, black, "#")
			}
		}
	}

	con.Set(player.X, player.Y, bold, black, "@")
	con.Set(monster.X, monster.Y, dred, black, "M")

	if !player.Thinking && actionQ.Len() > 0 {
		actor := actionQ.Pop()
		time := actor.Act()
		if time >= 0 {
			actionQ.Push(actor, time)
		}
	}
}

func (c *Console) Type(char rune) {
	if player.Thinking {
		switch char {
		case 'h':
			if grid.Get(player.X-1, player.Y) {
				player.X -= 1
			}
		case 'j':
			if grid.Get(player.X, player.Y+1) {
				player.Y += 1
			}
		case 'k':
			if grid.Get(player.X, player.Y-1) {
				player.Y -= 1
			}
		case 'l':
			if grid.Get(player.X+1, player.Y) {
				player.X += 1
			}
		}
		player.Thinking = false
	}
}

func main() {
	player = &Player{NewEntity(5, 12), false}
	monster = NewEntity(20, 9)

	actionQ = rog.NewQueue()
	actionQ.Push(player, 1)
	actionQ.Push(monster, 2)

	grid = rog.NewSparseGrid(32, 20)
	rog.DigArena(grid, 0, 0, 32, 20)
	rog.DigArena(grid, 8, 8, 8, 8)

	OpenConsole("rog", 32, 20, 2, "data/font2.png", 16, 16)
}

const (
	black    = 0x212121
	dgrey    = 0x42403c
	dred     = 0x903849
	red      = 0xb95f70
	dgreen   = 0x51673b
	green    = 0x7a915a
	dyellow  = 0x8a643e
	yellow   = 0xaa845d
	dblue    = 0x415860
	blue     = 0x6b818b
	dmagenta = 0x55445e
	magenta  = 0x816a89
	dcyan    = 0x50755a
	cyan     = 0x7c9a84
	grey     = 0x76716a
	white    = 0xafa79d
	bold     = 0xeeeeee
)
