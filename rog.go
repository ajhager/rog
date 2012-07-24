package rog

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
    "image/png"
    "os"
	"sync"
	"time"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
)

var (
	wg sync.WaitGroup
)

type driver func(*Window)

type Window struct {
	*Console
	win   wde.Window
	Dt    float64
	Fps   int64
	Mouse *Mouse
    Key string
}

func (this *Window) Close() {
	this.win.Close()
	wg.Done()
}

func (this* Window) Screenshot(name string) (err error) {
    file, err := os.Create(fmt.Sprintf("%v.%v", name, "png"))
    if err != nil {
        return
    }
    defer file.Close()

    err = png.Encode(file, this.win.Screen())
    return
}

func (this *Window) SetTitle(title string) {
	this.win.SetTitle(title)
}

type MouseButton struct {
	Pressed, Released bool
}

type Mouse struct {
	Pos, DPos, Cell, DCell image.Point
	Left, Right, Middle    MouseButton
}

func handleEvents(window *Window) (done bool) {
	window.Mouse.DPos.X = 0
	window.Mouse.DPos.Y = 0
	window.Mouse.DCell.X = 0
	window.Mouse.DCell.Y = 0
	window.Mouse.Left.Released = false
	window.Mouse.Right.Released = false
	window.Mouse.Middle.Released = false
    window.Key = ""
	select {
	case ei := <-window.win.EventChan():
		switch e := ei.(type) {
		case wde.MouseMovedEvent:
			window.Mouse.Pos.X = e.Where.X
			window.Mouse.Pos.Y = e.Where.Y
			window.Mouse.DPos.X = e.From.X
			window.Mouse.DPos.Y = e.From.Y
			window.Mouse.Cell.X = e.Where.X / 16
			window.Mouse.Cell.Y = e.Where.Y / 16
			window.Mouse.DCell.X = e.From.X / 16
			window.Mouse.DCell.Y = e.From.Y / 16
		case wde.MouseDownEvent:
			switch e.Which {
			case wde.LeftButton:
				window.Mouse.Left.Pressed = true
			case wde.RightButton:
				window.Mouse.Right.Pressed = true
			case wde.MiddleButton:
				window.Mouse.Right.Pressed = true
			}
		case wde.MouseUpEvent:
			switch e.Which {
			case wde.LeftButton:
				window.Mouse.Left.Pressed = false
				window.Mouse.Left.Released = true
			case wde.RightButton:
				window.Mouse.Right.Pressed = false
				window.Mouse.Right.Released = true
			case wde.MiddleButton:
				window.Mouse.Right.Pressed = false
				window.Mouse.Right.Released = true
			}
        case wde.KeyTypedEvent:
            window.Key = e.Key
		case wde.ResizeEvent:
		case wde.CloseEvent:
            window.Close()
            done = true
		}
	default:
	}
    return
}

func Open(width, height int, title string, driver driver) {
	wg.Add(1)
	go func() {
		dw, err := wde.NewWindow(width*16, height*16)
		if err != nil {
			fmt.Println(err)
			return
		}
		dw.SetTitle(title)
		dw.Show()

		console := NewConsole(width, height)
        backbuf := NewConsole(width, height)
		window := &Window{console, dw, 0, 0, new(Mouse), ""}

		f := font()
		buf := bytes.NewBuffer(f)
		mask, _, err := image.Decode(buf)
		if err != nil {
			panic(err)
		}

		oldTime := time.Now()
		newTime := time.Now()
		elapsed := float64(0)
		frames := int64(0)
		mr := image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
		for {
            // Handle events.
            if handleEvents(window) {
                break
            }

			// Update state of the console.
			driver(window)

			// Render the console to the screen
			for y := 0; y < window.h; y++ {
				for x := 0; x < window.w; x++ {
					bg := window.bg[y][x]
					fg := window.fg[y][x]
					ch := window.ch[y][x]
					if bg != backbuf.bg[y][x] ||
                       fg != backbuf.fg[y][x] ||
                       ch != backbuf.ch[y][x] {
                        backbuf.bg[y][x] = bg
                        backbuf.fg[y][x] = fg
                        backbuf.ch[y][x] = ch
						r := mr.Add(image.Point{x * 16, y * 16})
						src := &image.Uniform{bg}
						draw.Draw(window.win.Screen(), r, src, image.ZP, draw.Src)

						if ch != ' ' {
							src = &image.Uniform{fg}
							draw.DrawMask(window.win.Screen(), r, src, image.ZP, mask, image.Point{int(ch%32) * 16, int(ch/32) * 16}, draw.Over)
						}
					}
				}
			}
			window.win.FlushImage()

			oldTime = newTime
			newTime = time.Now()
			window.Dt = newTime.Sub(oldTime).Seconds()
			elapsed += window.Dt
			frames += 1
			if elapsed >= 1 {
				window.Fps = frames
				frames = 0
				elapsed -= elapsed
			}
		}
	}()
}

func Start() {
	go func() {
		wg.Wait()
		wde.Stop()
	}()
	wde.Run()
}
