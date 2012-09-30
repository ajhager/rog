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
         if rog.Key() == rog.Escape {
             rog.Close()
         }
         rog.Flush()
     }
 }
*/
package rog

import (
	"fmt"
	"image/png"
	"os"
)

var (
	workspace Workspace
	console   *Console
	timing    *stats
)

// IsOpen returns whether the rog window is open or not.
func IsOpen() bool {
	return workspace.IsOpen()
}

// Open creates a window and a root console with size width by height cells.
func Open(width, height int, title string, ws Workspace) {
	timing = new(stats)
	console = NewConsole(width, height)

	workspace = ws
	workspace.Open(width, height)
	workspace.Name(title)
}

// Close shuts down the windowing system.
// No rog functions should be called after this.
func Close() {
	workspace.Close()
}

// Screenshot will save the window buffer as an image to name.png.
func Screenshot(name string) (err error) {
	file, err := os.Create(fmt.Sprintf("%v.%v", name, "png"))
	if err != nil {
		return
	}
	defer file.Close()

	err = png.Encode(file, workspace.Screen())
	return
}

// SetTitle changes the title of the window.
func SetTitle(title string) {
	workspace.Name(title)
}

// Flush renders the root console to the window.
func Flush() {
	workspace.Render(console)
	timing.Update()
}

// Mouse returns a struct representing the state of the mouse.
func Mouse() *MouseData {
	return workspace.Mouse()
}

// Key returns the last key typed this frame.
func Key() string {
	return workspace.Key()
}

// Dt returns length of the last frame in seconds.
func Dt() float64 {
	return timing.Dt
}

// Fps returns the number of rendered frames per second.
func Fps() int64 {
	return timing.Fps
}

func Blit(con *Console, x, y int) {
    for i := 0; i < con.Width(); i++ {
        for j := 0; j < con.Height(); j++ {
            fg, bg, ch := con.Get(i, j)
            console.Set(x+i, y+j, fg, bg, string(byte(ch)))
        }
    }
}

// Set draws on the root console.
func Set(x, y int, fg, bg Blender, data string, rest ...interface{}) {
	console.Set(x, y, fg, bg, data, rest...)
}

// Set draws on the root console with wrapping bounds of x, y, w, h.
func SetR(x, y, w, h int, fg, bg Blender, data string, rest ...interface{}) {
	console.SetR(x, y, w, h, fg, bg, data, rest...)
}

// Get returns the fg, bg colors and rune of the cell on the root console.
func Get(x, y int) (RGB, RGB, rune) {
	return console.Get(x, y)
}

// Fill draws a rect on the root console.
func Fill(x, y, w, h int, fg, bg Blender, ch rune) {
	console.Fill(x, y, w, h, fg, bg, ch)
}

// Clear draws a rect over the entire root console.
func Clear(fg, bg Blender, ch rune) {
	console.Clear(fg, bg, ch)
}

// Width returns the width of the root console in cells.
func Width() int {
	return console.Width()
}

// Height returns the height of the root console in cells.
func Height() int {
	return console.Height()
}
