// Copyright 2012 Joseph Hager. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package rog provides algorithms and data structures for creating roguelike games.

 package main

 import (
     "github.com/ajhager/rog"
 )

 func main() {
     rog.Open(48, 32, "rog")
     for rog.IsOpen() {
         rog.Set(20, 15, nil, nil, "Hello, 世界")
         if rog.Key == "escape" {
             rog.Close()
         }
         rog.Flush()
     }
 }
*/
package rog

import (
	"fmt"
	"image"
    "image/png"
    "os"
	"time"
	"github.com/skelterjohn/go.wde"
	_ "github.com/skelterjohn/go.wde/init"
)

var (
    open = false
    window wde.Window
    console *Console
    stats *timing
    Mouse *mouse
    Key string
)

// IsOpen returns whether the rog window is open or not.
func IsOpen() bool {
    return open
}

// Open creates a window and a root console with size width by height cells.
func Open(width, height int, title string) (err error) {
	window, err = wde.NewWindow(width*16, height*16)
	if err != nil {
        return
	}
	window.SetTitle(title)
	window.Show()
    
	console = NewConsole(width, height)
    stats = new(timing)
    Mouse = new(mouse)
    
	go func() {
	    wde.Run()
	}()

    open = true
    return
}

// Close shuts down the windowing system.
// No rog functions should be called after this.
func Close() {
    open = false
	window.Close()
    wde.Stop()
}

// Screenshot will save the window buffer as an image to name.png.
func Screenshot(name string) (err error) {
    file, err := os.Create(fmt.Sprintf("%v.%v", name, "png"))
    if err != nil {
        return
    }
    defer file.Close()

    err = png.Encode(file, window.Screen())
    return
}

// SetTitle changes the title of the window.
func SetTitle(title string) {
	window.SetTitle(title)
}

// Flush renders the root console to the window.
func Flush() {
    handleEvents()
    if open {
        console.Render(window.Screen())
        window.FlushImage()
    }
    stats.Update(time.Now())
}

// Dt returns length of the last frame in seconds.
func Dt() float64 {
    return stats.Dt
}

// Fps returns the number of rendered frames per second.
func Fps() int64 {
    return stats.Fps
}

// Set draws on the root console.
func Set(x, y int, fg, bg interface{}, data string, rest ...interface{}) {
    console.Set(x, y, fg, bg, data, rest...)
}

// Set draws on the root console with wrapping bounds of x, y, w, h.
func SetR(x, y, w, h int, fg, bg interface{}, data string, rest ...interface{}) {
    console.SetR(x, y, w, h, fg, bg, data, rest...)
}

// Fill draws a rect on the root console.
func Fill(x, y, w, h int, fg, bg interface{}, ch rune) {
    console.Fill(x, y, w, h, fg, bg, ch)
}

// Width returns the width of the root console in cells.
func Width() int {
    return console.Width()
}

// Height returns the height of the root console in cells.
func Height() int {
    return console.Height()
}

func handleEvents() {
	Mouse.DPos.X = 0
	Mouse.DPos.Y = 0
	Mouse.DCell.X = 0
	Mouse.DCell.Y = 0
	Mouse.Left.Released = false
	Mouse.Right.Released = false
	Mouse.Middle.Released = false
    Key = ""
	select {
	case ei := <-window.EventChan():
		switch e := ei.(type) {
		case wde.MouseMovedEvent:
			Mouse.Pos.X = e.Where.X
			Mouse.Pos.Y = e.Where.Y
			Mouse.DPos.X = e.From.X
			Mouse.DPos.Y = e.From.Y
			Mouse.Cell.X = e.Where.X / 16
			Mouse.Cell.Y = e.Where.Y / 16
			Mouse.DCell.X = e.From.X / 16
			Mouse.DCell.Y = e.From.Y / 16
		case wde.MouseDownEvent:
			switch e.Which {
			case wde.LeftButton:
				Mouse.Left.Pressed = true
			case wde.RightButton:
				Mouse.Right.Pressed = true
			case wde.MiddleButton:
				Mouse.Right.Pressed = true
			}
		case wde.MouseUpEvent:
			switch e.Which {
			case wde.LeftButton:
				Mouse.Left.Pressed = false
				Mouse.Left.Released = true
			case wde.RightButton:
				Mouse.Right.Pressed = false
				Mouse.Right.Released = true
			case wde.MiddleButton:
				Mouse.Right.Pressed = false
				Mouse.Right.Released = true
			}
        case wde.KeyTypedEvent:
            Key = e.Key
		case wde.ResizeEvent:
		case wde.CloseEvent:
            Close()
		}
	default:
	}
}

type mouseButton struct {
	Pressed, Released bool
}

type mouse struct {
	Pos, DPos, Cell, DCell image.Point
	Left, Right, Middle mouseButton
}

type timing struct {
    Then, Now time.Time
    Elapsed, Dt float64
    Frames, Fps int64
}

func (t *timing) Update(now time.Time) {
    t.Then = t.Now
    t.Now = now
    t.Dt = t.Now.Sub(t.Then).Seconds()
    t.Elapsed += t.Dt
    t.Frames += 1
	if t.Elapsed >= 1 {
	    t.Fps = t.Frames
        t.Frames = 0
        t.Elapsed -= t.Elapsed
    }
}
