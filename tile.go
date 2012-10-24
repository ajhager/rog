package rog

type Renderable interface {
    GetX() int
    GetY() int
    SetX(int)
    SetY(int)

    GetFg() RGB
    GetBg() RGB
    GetGlyph() string
    Render(dx, dy int)
}

type Tileable interface {
    Renderable
    Viewable
}

type Tile struct {
    X, Y int
    Fg, Bg *RGB
    Glyph rune
    Roughness int
    Viewable bool
}

func (self *Tile) GetX() int { return self.X }
func (self *Tile) GetY() int { return self.Y }

func (self *Tile) SetX(v int) { self.X = v }
func (self *Tile) SetY(v int) { self.Y = v }

func (self *Tile) GetFg() RGB { return *self.Fg }
func (self *Tile) GetBg() RGB { return *self.Bg }
func (self *Tile) GetGlyph() string { return string(self.Glyph) }
func (self *Tile) GetRoughness() int { return self.Roughness }
func (self *Tile) IsViewable() bool { return self.Viewable }
func (self *Tile) Render(dx, dy int) {
    Set(
        self.GetX() + dx, self.GetY() + dy,
        self.GetFg(), self.GetBg(),
        string(self.GetGlyph()),
    )
}
