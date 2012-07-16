package rog

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
	"runtime"
	"sync"
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
	win wde.Window
	Dt  float64
	*Console
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

func Open(width, height int, title string, driver driver) {
	wg.Add(1)
	go func() {
		dw, err := wde.NewWindow(width*16+4, height*16+4)
		if err != nil {
			fmt.Println(err)
			return
		}
		dw.SetTitle(title)
		dw.Show()

		console := NewConsole(width, height)
		window := &Window{dw, 0, console}

		f := font()
		buf := bytes.NewBuffer(f)
		mask, _, err := image.Decode(buf)
		if err != nil {
			panic(err)
		}

		events := dw.EventChan()
		done := make(chan bool)

		go func() {
		loop:
			for ei := range events {
				runtime.Gosched()
				switch e := ei.(type) {
				case wde.KeyTypedEvent:
					fmt.Println("KeyDownEvent", e.Glyph)
				case wde.CloseEvent:
					dw.Close()
					break loop
				}
			}
			done <- true
		}()

		for {
			screen := dw.Screen()
			draw.Draw(screen, screen.Bounds(), &image.Uniform{color.Black}, image.ZP, draw.Src)

			// Update state of the console.
			driver(window)

			mr := image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
			// Render the console to the screen
			for y := 0; y < window.h; y++ {
				for x := 0; x < window.w; x++ {
					r := mr.Add(image.Point{x * 16, y * 16})
					bg := window.bg[y][x]
					src := &image.Uniform{bg}
					draw.Draw(screen, r, src, image.ZP, draw.Src)

					ch := window.ch[y][x]
					fg := window.fg[y][x]
					if ch != 0 && ch != 32 {
						src = &image.Uniform{fg}
						draw.DrawMask(screen, r, src, image.ZP, mask, image.Point{int(ch%32) * 16, int(ch/32) * 16}, draw.Over)
					}
				}
			}
			dw.FlushImage()

			// Check for console close.	
			select {
			case <-done:
				wg.Done()
				return
			default:
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
