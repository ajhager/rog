package rog

type Renderable interface {
    X() int
    Y() int
    SetX(int)
    SetY(int)

    Fg() RGB
    Bg() RGB
    Glyph() rune
    Render()
}

type Tileable interface {
    Renderable
    Viewable
}

type Tile struct {
    x, y int

    fg, bg RGB
    glyph rune
    blocks bool
}

func (self *Tile) X() int { return self.x }
func (self *Tile) Y() int { return self.y }

func (self *Tile) SetX(v int) { self.x = v }
func (self *Tile) SetY(v int) { self.y = v }

func (self *Tile) Fg() RGB { return self.fg }
func (self *Tile) Bg() RGB { return self.bg }
func (self *Tile) Glyph() rune { return self.glyph }
func (self *Tile) Blocks() bool { return self.blocks }
func (self *Tile) Render(dx, dy int) {
    Set(
        self.X() + dx, self.Y() + dy,
        self.Fg(), self.Bg(),
        string(self.Glyph()),
    )
}

type Floor struct { Tile }
func NewFloor(x, y int) *Floor {
    return &Floor{
        Tile{
            x, y,
            Hex(0x885040), Hex(0x885040),
            '.', false,

        },
    }
}

type Wall struct { Tile }
func NewWall(x, y int) *Wall {
    return &Wall{
        Tile{
            x, y,
            Hex(0xffbb99), Hex(0xffbb99),
            '#', true,
        },
    }
}