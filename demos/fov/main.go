package main

import (
	"github.com/ajhager/rog"
	"runtime"
	"image"
)

func intensity(px, py, cx, cy, r int) float64 {
	r2 := float64(r * r)
	squaredDist := float64((px-cx)*(px-cx) + (py-cy)*(py-cy))
	coef1 := 1.0 / (1.0 + squaredDist/20)
	coef2 := coef1 - 1.0/(1.0+r2)
	return coef2 / (1.0 - 1.0/(1.0+r2))
}

func load_map(width, height int, tmap [][]rune) *rog.Map {
	new_map := rog.NewMap(width, height)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			switch tmap[y][x] {
				case '#':
					new_map.SetTile(NewWall(x, y), x, y)
				case ' ':
					new_map.SetTile(NewFloor(x, y), x, y)
			}
		}
	}
	return new_map
}

var (
	width  = 40
	height = 20

	frame_time = 1000.0 / 23.0

	lgrey  = *rog.Hex(0xc8c8c8)
	dgrey  = *rog.Hex(0x1e1e1e)

	stats runtime.MemStats

	tmap = [][]rune{
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
	}
)

func main() {
	rog.Open(width, height, 1, false, "Example", nil)

	tile_map := load_map(width, height, tmap)

	player := NewPlayer(23, 9)

	tile_map.Fov(player.GetX(), player.GetY(), 20, true, rog.FOVCircular)

	scenes := rog.NewSceneStack()

	first := GameScene {
		stack: scenes,
		tile_map: tile_map, 
		player: player, 
		path: nil,
	}

	scenes.Push(&first)

	frame_timer := rog.NewTimer(frame_time)

	has_scene := func() bool {
		if scenes.Top() != nil {
			return true
		}
		rog.Close()
		return false
	}

	for rog.Running() {
		if rog.CheckTimer(frame_timer) {
			if has_scene() {
				scenes.Top().HandleKeys()
			}
			if has_scene() {
				scenes.Top().Update()
			}
			if has_scene() {
				scenes.Top().Render()
			}
			rog.Flush()
		}
	}
}

type Floor struct { rog.Tile }
func NewFloor(x, y int) *Floor {
    return &Floor{
        rog.Tile {
            X: x, Y: y,
            Fg: rog.Hex(0x442020), 
            Bg: rog.Hex(0x885040),
            Glyph: '.', 
            Roughness: rog.PATH_MIN, 
            Viewable: true,
        },
    }
}

type Wall struct { rog.Tile }
func NewWall(x, y int) *Wall {
    return &Wall{
    	rog.Tile {
		    X: x, Y: y,
		    Fg: rog.Hex(0x885544), 
		    Bg: rog.Hex(0xffbb99),
		    Glyph: '#', 
		    Roughness: rog.PATH_MAX, 
		    Viewable: false,
		},
	}
}

type Player struct { rog.Tile }
func NewPlayer(x, y int) *Player {
	return &Player {
		rog.Tile {
			X: x, Y: y,
			Fg: &rog.White,
			Bg: nil,
			Glyph: '웃',
			Roughness: rog.PATH_MIN,
			Viewable: true,
		},
	}
}

type GameScene struct {
	stack *rog.SceneStack
	tile_map *rog.Map
    player *Player
	path []image.Point
}

func (self *GameScene) MovePlayer(x, y int) {
	tile, ok := self.tile_map.GetTile(x, y)
	if ok {
		in_bounds := x >= 0 && y > 0 && x < width && y < height-1
		is_floor := tile.GetRoughness() < rog.PATH_MAX
		if  in_bounds && is_floor {
			rog.Set(self.player.GetX(), self.player.GetY(), rog.White, nil, " ")
			self.player.SetX(x)
			self.player.SetY(y)
			self.tile_map.Fov(x, y, 20, true, rog.FOVCircular)
		}
	}
}

func (self *GameScene) HandleKeys() {
	x := self.player.GetX()
	y := self.player.GetY()

    switch rog.Key() {
        case rog.Esc:
            self.stack.Pop()
		case 'k', 'w', rog.Up:
			self.MovePlayer(x, y-1)
		case 'j', 's', rog.Down:
			self.MovePlayer(x, y+1)
		case 'h', 'a', rog.Left:
			self.MovePlayer(x-1, y)
		case 'l', 'd', rog.Right:
			self.MovePlayer(x+1, y)
    }

	if rog.Mouse().Left.Released {
		self.path = self.tile_map.Path(x, y, rog.Mouse().Cell.X, rog.Mouse().Cell.Y)
	}

}

func (self *GameScene) Update() {

}

func (self *GameScene) Render() {
	px := self.player.GetX()
	py := self.player.GetY()

	rog.Clear(nil, rog.Black, ' ')

	for point, _ := range self.tile_map.ViewMap() {
		tile, _ := self.tile_map.GetTile(point.X, point.Y)
		i := intensity(px, py, point.X, point.Y, 20)
		rog.Set(point.X, point.Y, tile.GetFg().Scale(i), tile.GetBg().Scale(i), string(tile.GetGlyph()))
	}

	rog.Set(px,  py, lgrey, nil, self.player.GetGlyph())

	for _, p := range self.path {
		if p.X != px || p.Y != py {
			rog.Set(p.X, p.Y, lgrey, nil, "*")
		}
	}

	runtime.ReadMemStats(&stats)
	rog.Fill(0, 0, rog.Width(), 1, lgrey, rog.Dodge(dgrey), ' ')
	rog.Set(0, 0, nil, nil, "%vFPS %vMB %vGC %vGR", rog.Fps(), stats.HeapAlloc/1000000, stats.NumGC, runtime.NumGoroutine())
	rog.Fill(0, height-1, rog.Width(), 1, lgrey, rog.Dodge(dgrey), ' ')
	rog.Set(0, height-1, nil, nil, "Pos: %v %v Cell: %v %v", rog.Mouse().Pos.X, rog.Mouse().Pos.Y, rog.Mouse().Cell.X, rog.Mouse().Cell.Y)
}
