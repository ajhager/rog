package term

import (
	"github.com/ajhager/rog"
	"github.com/nsf/termbox-go"
	"image"
)

func Backend() rog.Backend {
	return new(termboxBackend)
}

type termboxBackend struct {
	open  bool
	mouse *rog.MouseData
	key   string
}

func (w *termboxBackend) Open(width, height, zoom int) {
	if err := termbox.Init(); err != nil {
		panic(err)
	}

	w.mouse = new(rog.MouseData)
	println(termbox.ColorBlue | termbox.AttrBold)

	go w.pollKeys()

	w.open = true
}

func (w *termboxBackend) IsOpen() bool {
	return w.open
}

func (w *termboxBackend) Close() {
	w.open = false
	termbox.Close()
}

func (w *termboxBackend) Name(title string) {
}

func cLevel(i uint8) int {
	if i < 64 {
		return 1
	} else if i < 128 {
		return 2
	} else if i < 192 {
		return 3
	}
	return 4
}

var cols map[int]termbox.Attribute = map[int]termbox.Attribute{
	111: termbox.ColorBlack,
	121: termbox.ColorGreen,
	131: termbox.ColorGreen,
	141: termbox.ColorGreen | termbox.AttrBold,
	112: termbox.ColorBlue,
	113: termbox.ColorBlue,
	114: termbox.ColorBlue | termbox.AttrBold,
	122: termbox.ColorCyan,
	123: termbox.ColorBlue,
	124: termbox.ColorBlue | termbox.AttrBold,
	132: termbox.ColorGreen,
	133: termbox.ColorCyan,
	134: termbox.ColorBlue | termbox.AttrBold,
	142: termbox.ColorGreen | termbox.AttrBold,
	143: termbox.ColorCyan | termbox.AttrBold,
	144: termbox.ColorCyan | termbox.AttrBold,
	211: termbox.ColorRed,
	221: termbox.ColorYellow,
	231: termbox.ColorGreen,
	241: termbox.ColorGreen | termbox.AttrBold,
	212: termbox.ColorMagenta,
	213: termbox.ColorMagenta,
	214: termbox.ColorBlue,
	222: termbox.ColorBlack | termbox.AttrBold,
	223: termbox.ColorBlue,
	224: termbox.ColorBlue | termbox.AttrBold,
	232: termbox.ColorGreen,
	233: termbox.ColorCyan,
	234: termbox.ColorBlue | termbox.AttrBold,
	242: termbox.ColorGreen | termbox.AttrBold,
	243: termbox.ColorCyan | termbox.AttrBold,
	244: termbox.ColorCyan | termbox.AttrBold,
	311: termbox.ColorRed,
	321: termbox.ColorYellow,
	331: termbox.ColorYellow | termbox.AttrBold,
	341: termbox.ColorGreen | termbox.AttrBold,
	312: termbox.ColorMagenta,
	313: termbox.ColorMagenta,
	314: termbox.ColorMagenta,
	322: termbox.ColorRed,
	323: termbox.ColorMagenta | termbox.AttrBold,
	324: termbox.ColorMagenta | termbox.AttrBold,
	332: termbox.ColorYellow | termbox.AttrBold,
	333: termbox.ColorWhite,
	334: termbox.ColorBlue | termbox.AttrBold,
	342: termbox.ColorGreen | termbox.AttrBold,
	343: termbox.ColorGreen | termbox.AttrBold,
	344: termbox.ColorCyan | termbox.AttrBold,
	411: termbox.ColorRed,
	421: termbox.ColorRed | termbox.AttrBold,
	431: termbox.ColorYellow | termbox.AttrBold,
	441: termbox.ColorYellow | termbox.AttrBold,
	412: termbox.ColorBlue | termbox.AttrBold,
	413: termbox.ColorBlue | termbox.AttrBold,
	414: termbox.ColorBlue | termbox.AttrBold,
	422: termbox.ColorRed | termbox.AttrBold,
	423: termbox.ColorMagenta | termbox.AttrBold,
	424: termbox.ColorMagenta | termbox.AttrBold,
	432: termbox.ColorYellow | termbox.AttrBold,
	433: termbox.ColorRed | termbox.AttrBold,
	434: termbox.ColorMagenta | termbox.AttrBold,
	442: termbox.ColorYellow | termbox.AttrBold,
	443: termbox.ColorYellow | termbox.AttrBold,
	444: termbox.ColorWhite | termbox.AttrBold}

func rgbToAnsi(c rog.RGB) termbox.Attribute {
	n := cLevel(c.R)*100 + cLevel(c.G)*10 + cLevel(c.B)
	return cols[n]
}

func (w *termboxBackend) Render(console *rog.Console) {
	if w.IsOpen() {
		w.key = ""

		termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				fg, bg, ch := console.Get(x, y)
				termbox.SetCell(x, y, ch, rgbToAnsi(fg), rgbToAnsi(bg))
			}
		}
		termbox.Flush()
	}
}

func (w *termboxBackend) Screen() image.Image {
	return &image.Uniform{rog.RGB{255, 255, 255}}
}

func (w *termboxBackend) Mouse() *rog.MouseData {
	return w.mouse
}

func (w *termboxBackend) pollKeys() {
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if ev.Ch == 0 {
				switch ev.Key {
				case termbox.KeyEsc:
					w.key = rog.Escape
				}
			} else {
				switch ev.Ch {
				case 'h':
					w.key = rog.H
				case 'j':
					w.key = rog.J
				case 'k':
					w.key = rog.K
				case 'l':
					w.key = rog.L
				}
			}
		}
	}
}

func (w *termboxBackend) Key() string {
	return w.key
}
