package wde

import (
	"bytes"
	"github.com/ajhager/rog"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
	"image"
	"image/color"
	"image/draw"
    _ "image/png"
)

func Backend() rog.Backend {
	return new(wdeBackend)
}

type nearestNeighborImage struct {
	image image.Image
	Zoom  int
}

func (nni nearestNeighborImage) ColorModel() color.Model {
	return nni.image.ColorModel()
}

func (nni nearestNeighborImage) Bounds() image.Rectangle {
	b := nni.image.Bounds()
	return image.Rect(b.Min.X, b.Min.Y, b.Max.X*nni.Zoom, b.Max.Y*nni.Zoom)
}

func (nni nearestNeighborImage) At(x, y int) color.Color {
	x = (x - x%nni.Zoom) / nni.Zoom
	y = (y - y%nni.Zoom) / nni.Zoom
	return nni.image.At(x, y)
}

type wdeBackend struct {
	open         bool
	window       wde.Window
	input        chan interface{}
	mouse        *rog.MouseData
	key          int
	bgbuf, fgbuf [][]color.Color
	chbuf        [][]rune
	font         image.Image
	zoom         int
}

func (w *wdeBackend) Open(width, height, zoom int) {
	w.window, _ = wde.NewWindow(width*16*zoom, height*16*zoom)
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
	w.zoom = zoom

	for y := 0; y < height; y++ {
		w.bgbuf[y] = make([]color.Color, width)
		w.fgbuf[y] = make([]color.Color, width)
		w.chbuf[y] = make([]rune, width)
	}

	font, _, err := image.Decode(bytes.NewBuffer(rog.FontData()))
	if err != nil {
		panic(err)
	}
	w.font = nearestNeighborImage{font, zoom}

	w.open = true
}

func (w *wdeBackend) IsOpen() bool {
	return w.open
}

func (w *wdeBackend) Close() {
	w.open = false
	w.window.Close()
	wde.Stop()
}

func (w *wdeBackend) Name(title string) {
	w.window.SetTitle(title)
}

func (w *wdeBackend) Render(console *rog.Console) {
	if w.IsOpen() {
		w.handleFrameEvents()

		im := w.window.Screen()
		maskRect := image.Rectangle{image.Point{0, 0}, image.Point{16 * w.zoom, 16 * w.zoom}}
		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				fg, bg, ch := console.Get(x, y)
				if bg != w.bgbuf[y][x] || fg != w.fgbuf[y][x] || ch != w.chbuf[y][x] {
					w.bgbuf[y][x] = bg
					w.fgbuf[y][x] = fg
					w.chbuf[y][x] = ch
					rect := maskRect.Add(image.Point{x * 16 * w.zoom, y * 16 * w.zoom})
					src := &image.Uniform{bg}
					draw.Draw(im, rect, src, image.ZP, draw.Src)

					if ch != ' ' {
						src = &image.Uniform{fg}
						draw.DrawMask(im, rect, src, image.ZP, w.font, image.Point{int(ch%256) * 16 * w.zoom, int(ch/256) * 16 * w.zoom}, draw.Over)
					}
				}
			}
		}

		w.window.FlushImage()
	}
}

func (w *wdeBackend) Mouse() *rog.MouseData {
	return w.mouse
}

func (w *wdeBackend) Key() int {
	return w.key
}

func (w *wdeBackend) handleRealtimeEvents() {
	for ei := range w.window.EventChan() {
		switch e := ei.(type) {
		case wde.MouseMovedEvent:
			w.mouse.Pos = e.Where
			w.mouse.Cell = e.Where.Div(16 * w.zoom)
		case wde.MouseDraggedEvent:
			w.mouse.Pos = e.Where
			w.mouse.Cell = e.Where.Div(16 * w.zoom)
		case wde.CloseEvent:
			w.Close()
		default:
			w.input <- ei
		}
	}
}

func (w *wdeBackend) handleFrameEvents() {
	w.mouse.Left.Released = false
	w.mouse.Right.Released = false
	w.mouse.Middle.Released = false
	w.key = -1
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
            if key, ok := wdeToRogKey[e.Key]; ok {
			    w.key = key
            }
		}
	default:
	}
}

var wdeToRogKey map[string]int = map[string]int {
	wde.KeyBackspace: rog.Backspace,
	wde.KeyTab: rog.Tab,
	wde.KeySpace: rog.Space,
	wde.KeyDelete: rog.Delete,
	wde.KeyReturn: rog.Return,
	wde.KeyA: 65,
	wde.KeyB: 66,
	wde.KeyC: 67,
	wde.KeyD: 68,
	wde.KeyE: 69,
	wde.KeyF: 70,
	wde.KeyG: 71,
	wde.KeyH: 72,
	wde.KeyI: 73,
	wde.KeyJ: 74,
	wde.KeyK: 75,
	wde.KeyL: 76,
	wde.KeyM: 77,
	wde.KeyN: 78,
	wde.KeyO: 79,
	wde.KeyP: 80,
	wde.KeyQ: 81,
	wde.KeyR: 82,
	wde.KeyS: 83,
	wde.KeyT: 84,
	wde.KeyU: 85,
	wde.KeyV: 86,
	wde.KeyW: 87,
	wde.KeyX: 88,
	wde.KeyY: 89,
	wde.KeyZ: 90,
	wde.Key0: 48,
	wde.Key1: 49,
	wde.Key2: 50,
	wde.Key3: 51,
	wde.Key4: 52,
	wde.Key5: 53,
	wde.Key6: 54,
	wde.Key7: 55,
	wde.Key8: 56,
	wde.Key9: 57,
	wde.KeyLeftSuper: rog.LSuper,
	wde.KeyRightSuper: rog.RSuper,
	wde.KeyLeftShift: rog.LShift,
	wde.KeyRightShift: rog.RShift,
	wde.KeyLeftControl: rog.LControl,
	wde.KeyRightControl: rog.RControl,
	wde.KeyLeftAlt: rog.LAlt,
	wde.KeyRightAlt: rog.RAlt,
	wde.KeyF1: rog.F1,
	wde.KeyF2: rog.F2,
	wde.KeyF3: rog.F3,
	wde.KeyF4: rog.F4,
	wde.KeyF5: rog.F5,
	wde.KeyF6: rog.F6,
	wde.KeyF7: rog.F7,
	wde.KeyF8: rog.F8,
	wde.KeyF9: rog.F9,
	wde.KeyF10: rog.F10,
	wde.KeyF11: rog.F11,
	wde.KeyF12: rog.F12,
	wde.KeyF13: rog.F13,
	wde.KeyF14: rog.F14,
	wde.KeyF15: rog.F15,
	wde.KeyF16: rog.F16,
	wde.KeyUpArrow: rog.Up,
	wde.KeyDownArrow: rog.Down,
	wde.KeyLeftArrow: rog.Left,
	wde.KeyRightArrow: rog.Right,
	wde.KeyInsert: rog.Insert,
	wde.KeyHome: rog.Home,
	wde.KeyEnd: rog.End,
	wde.KeyCapsLock: rog.Capslock,
	wde.KeyPadSlash: rog.KPDivide,
	wde.KeyPadStar: rog.KPMultiply,
	wde.KeyPadMinus: rog.KPSubtract,
	wde.KeyPadPlus: rog.KPAdd,
	wde.KeyPadDot: rog.KPDecimal,
	wde.KeyPadEqual: rog.KPEqual,
	wde.KeyPadEnter: rog.KPEnter,
	wde.KeyNumlock: rog.KPNumlock,
	wde.KeyBackTick: 96,
	wde.KeyMinus: 45,
	wde.KeyEqual: 61,
	wde.KeyLeftBracket: 91,
	wde.KeyRightBracket: 93,
	wde.KeyBackslash: 92,
	wde.KeySemicolon: 59,
	wde.KeyQuote: 39,
	wde.KeyComma: 44,
	wde.KeyPeriod: 46,
	wde.KeySlash: 47,
	wde.KeyEscape: rog.Escape}
