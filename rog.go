package rog

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"runtime"
	"sync"
	"time"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
	wg sync.WaitGroup
)

type driver func(*Window)
type drawer func(draw.Image)

type Window struct {
	*Console
	win   wde.Window
	Dt    float64
	Fps   int64
	Mouse *Mouse
}

func (this *Window) Close() {
	this.win.Close()
}

func (this *Window) Draw(drawer drawer) {
	drawer(this.win.Screen())
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

func Open(width, height int, title string, driver driver) {
	wg.Add(1)
	go func() {
		dw, err := wde.NewWindow(width*16, height*16)
		ww, wh := dw.Size()
		fmt.Printf("w:%v h:%v\n", ww, wh)
		if err != nil {
			fmt.Println(err)
			return
		}
		dw.SetTitle(title)
		dw.Show()

		console := NewConsole(width, height)
		window := &Window{console, dw, 0, 0, new(Mouse)}

		f := font()
		buf := bytes.NewBuffer(f)
		mask, _, err := image.Decode(buf)
		if err != nil {
			panic(err)
		}

		events := dw.EventChan()
		oldTime := time.Now()
		newTime := time.Now()
		elapsed := float64(0)
		frames := int64(0)
		mr := image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
		for {
			window.Mouse.DPos.X = 0
			window.Mouse.DPos.Y = 0
			window.Mouse.DCell.X = 0
			window.Mouse.DCell.Y = 0
			window.Mouse.Left.Released = false
			window.Mouse.Right.Released = false
			window.Mouse.Middle.Released = false
			select {
			case ei := <-events:
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
				case wde.ResizeEvent:
					console.Dirty()
				case wde.CloseEvent:
					dw.Close()
					wg.Done()
				}
			default:
			}

			screen := dw.Screen()

			// Update state of the console.
			driver(window)

			// Render the console to the screen
			for y := 0; y < window.h; y++ {
				for x := 0; x < window.w; x++ {
					if window.dirt[y][x] {
						window.dirt[y][x] = false
						r := mr.Add(image.Point{x * 16, y * 16})
						bg := window.bg[y][x]
						src := &image.Uniform{bg}
						draw.Draw(screen, r, src, image.ZP, draw.Src)

						ch := window.ch[y][x]
						fg := window.fg[y][x]
						if ch != ' ' {
							src = &image.Uniform{fg}
							draw.DrawMask(screen, r, src, image.ZP, mask, image.Point{int(ch%32) * 16, int(ch/32) * 16}, draw.Over)
						}
					}
				}
			}
			dw.FlushImage()

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

func Stop() {
	wde.Stop()
}
