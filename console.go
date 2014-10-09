package rog

import "github.com/ajhager/engi"

type Console struct {
	*Buffer
	font  *Font
	Batch *engi.Batch
}

func (con *Console) Preload() {
	engi.Files.Add("font", con.font.path)
}

func (con *Console) Setup() {
	con.Batch = engi.NewBatch(float32(con.Width()*con.font.cellWidth), float32(con.Height()*con.font.cellHeight))
	con.font.native = engi.NewGridFont(engi.Files.Image("font"), con.font.cellWidth, con.font.cellHeight)
}

func (con *Console) Render() {
	con.Batch.Begin()
	con.Draw(con.Batch, con.font.native, con.font.cellWidth, con.font.cellHeight)
	con.Batch.End()
}

func (c *Console) Close()                                 {}
func (c *Console) Update(dt float32)                      {}
func (c *Console) Resize(w, h int)                        {}
func (c *Console) Mouse(x, y float32, action engi.Action) {}
func (c *Console) Scroll(amount float32)                  {}
func (c *Console) Key(key engi.Key, modifier engi.Modifier, action engi.Action) {
	if key == engi.Escape {
		engi.Exit()
	}
}
func (c *Console) Type(char rune) {}
