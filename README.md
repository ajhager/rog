rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit, scaling console with built in unicode font
* Cross platform with pluggable backends
* Field of view, lighting, and pathfinding algorithms
* Procedural color and palette creation

```go
package main

import (
    "github.com/ajhager/rog"
    "github.com/ajhager/rog/wde"
)

func main() {
    rog.Open(20, 11, 2, "rog", wde.Backend())
    for rog.IsOpen() {
        rog.Set(5, 5, nil, nil, "Hello, 世界!")
        if rog.Key() == rog.Escape {
            rog.Close()
        }
        rog.Flush()
    }
}
```

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.
* The glfw backend is the most performant and stable at the moment, but wde has less dependencies.

Thanks
------
* libTCOD [http://doryen.eptalys.net/libtcod/]
* Yamamushi [http://www.theasciiproject.com/]
* John Asmuth [http://github.com/skelterjohn/go.wde]
* jteeuwen [http://github.com/jteeuwen/glfw]
* nsf [http://github.com/nsf/termbox-go]

Plans
-----
* Website, documentation, tutorial, and more examples.
* Window resizing and fullscreen
* Better keyboard handling
* User supplied, non-square, and possilby ttf fonts
* Image blitting
* Random color Blenders
* Palette support
* Merge in lighting
* World creation
* Noise generators
* More fov algorithms
* Curses like API
* Common GUI widgets
