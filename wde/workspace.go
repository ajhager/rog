package wde

import (
	"bytes"
	"github.com/ajhager/rog"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/color"
	"image/draw"
)

func Backend() *wdeWorkspace {
	return new(wdeWorkspace)
}

type wdeWorkspace struct {
	open         bool
	window       wde.Window
	input        chan interface{}
	mouse        *rog.MouseData
	key          string
	bgbuf, fgbuf [][]color.Color
	chbuf        [][]rune
	font         image.Image
}

func (w *wdeWorkspace) Open(width, height int) {
	w.window, _ = wde.NewWindow(width*16, height*16)
	w.window.Show()
	go func() {
		wde.Run()
	}()
	w.mouse = new(rog.MouseData)
	w.input = make(chan interface{}, 16)
	go w.handleRealtimeEvents()

	w.bgbuf = make([][]color.Color, height)
	w.fgbuf = make([][]color.Color, height)
	w.chbuf = make([][]rune, height)

	for y := 0; y < height; y++ {
		w.bgbuf[y] = make([]color.Color, width)
		w.fgbuf[y] = make([]color.Color, width)
		w.chbuf[y] = make([]rune, width)
	}

	font, _, err := image.Decode(bytes.NewBuffer(rog.FontData()))
	if err != nil {
		panic(err)
	}
	w.font = font

	w.open = true
}

func (w *wdeWorkspace) IsOpen() bool {
	return w.open
}

func (w *wdeWorkspace) Close() {
	w.open = false
	w.window.Close()
	wde.Stop()
}

func (w *wdeWorkspace) Name(title string) {
	w.window.SetTitle(title)
}

func (w *wdeWorkspace) Render(console *rog.Console) {
	if w.IsOpen() {
		w.handleFrameEvents()

		im := w.window.Screen()
		maskRect := image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				fg, bg, ch := console.Get(x, y)
				if bg != w.bgbuf[y][x] || fg != w.fgbuf[y][x] || ch != w.chbuf[y][x] {
					w.bgbuf[y][x] = bg
					w.fgbuf[y][x] = fg
					w.chbuf[y][x] = ch
					rect := maskRect.Add(image.Point{x * 16, y * 16})
					src := &image.Uniform{bg}
					draw.Draw(im, rect, src, image.ZP, draw.Src)

					if ch != ' ' {
						src = &image.Uniform{fg}
						draw.DrawMask(im, rect, src, image.ZP, w.font, image.Point{int(ch%256) * 16, int(ch/256) * 16}, draw.Over)
					}
				}
			}
		}

		w.window.FlushImage()
	}
}

func (w *wdeWorkspace) Screen() image.Image {
	return w.window.Screen()
}

func (w *wdeWorkspace) Mouse() *rog.MouseData {
	return w.mouse
}

func (w *wdeWorkspace) Key() string {
	return w.key
}

func (w *wdeWorkspace) handleRealtimeEvents() {
	for ei := range w.window.EventChan() {
		switch e := ei.(type) {
		case wde.MouseMovedEvent:
			w.mouse.Pos = e.Where
			w.mouse.DPos = e.From
			w.mouse.Cell = e.Where.Div(16)
			w.mouse.DCell = e.From.Div(16)
		case wde.MouseDraggedEvent:
			w.mouse.Pos = e.Where
			w.mouse.DPos = e.From
			w.mouse.Cell = e.Where.Div(16)
			w.mouse.DCell = e.From.Div(16)
		case wde.CloseEvent:
			w.Close()
		default:
			w.input <- ei
		}
	}
}

func (w *wdeWorkspace) handleFrameEvents() {
	w.mouse.DPos = image.ZP
	w.mouse.DCell = image.ZP
	w.mouse.Left.Released = false
	w.mouse.Right.Released = false
	w.mouse.Middle.Released = false
	w.key = ""
	select {
	case ei := <-w.input:
		switch e := ei.(type) {
		case wde.MouseDownEvent:
			switch e.Which {
			case wde.LeftButton:
				w.mouse.Left.Pressed = true
			case wde.RightButton:
				w.mouse.Right.Pressed = true
			case wde.MiddleButton:
				w.mouse.Middle.Pressed = true
			}
		case wde.MouseUpEvent:
			switch e.Which {
			case wde.LeftButton:
				w.mouse.Left.Pressed = false
				w.mouse.Left.Released = true
			case wde.RightButton:
				w.mouse.Right.Pressed = false
				w.mouse.Right.Released = true
			case wde.MiddleButton:
				w.mouse.Middle.Pressed = false
				w.mouse.Middle.Released = true
			}
		case wde.KeyTypedEvent:
			w.key = e.Key
		}
	default:
	}
}
