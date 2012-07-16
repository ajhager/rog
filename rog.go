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
	win wde.Window
	Dt  float64
	Fps int64
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
		window := &Window{console, dw, 0, 0}

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

		oldTime := time.Now()
		newTime := time.Now()
		elapsed := float64(0)
		frames := int64(0)
		mr := image.Rectangle{image.Point{0, 0}, image.Point{16, 16}}
		for {
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
