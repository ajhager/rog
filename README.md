rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit, scaling console with built in unicode font
* Cross platform with pluggable backends
* Field of view, lighting, and pathfinding algorithms
* Procedural color and palette creation

[Documentation](http://go.pkgdoc.org/github.com/ajhager/rog "Documentation")

```go
package main

import (
    "github.com/ajhager/rog"
    _ "github.com/ajhager/rog/wde"
)

func main() {
    rog.Open(20, 11, 2, "rog")
    for rog.IsOpen() {
        rog.Set(5, 5, nil, nil, "Hello, 世界!")
        if rog.Key() == rog.Escape {
            rog.Close()
        }
        rog.Flush()
    }
}
```

Backends
--------
* wde:   No dependencies, software rasterizer, somewhat unstable
* glfw:  GLFW dynamic libs needed, opengl rendering, fast and stable
* term:  Runs in a terminal, RGB->Ansi color, not feature complete yet
* html:  Coming soon?

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* libTCOD [http://doryen.eptalys.net/libtcod/]
* Yamamushi [http://www.theasciiproject.com/]
* John Asmuth [http://github.com/skelterjohn/go.wde]
* jteeuwen [http://github.com/jteeuwen/glfw]
* nsf [http://github.com/nsf/termbox-go]

Plans
-----
* Website, tutorial, and more demos
* Window resizing and fullscreen
* User supplied, non-square, and ttf fonts
* Image blitting
* Standard colors by name
* Random color Blenders
* Palette support
* Merge in lighting
* World creation
* Noise generators
* More fov algorithms
* Curses like API
* Common GUI widgets
