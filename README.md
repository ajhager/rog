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

Thanks
------
* libTCOD [http://doryen.eptalys.net/libtcod/]
* Yamamushi [http://www.theasciiproject.com/]
* John Asmuth [http://github.com/skelterjohn/go.wde]
* jteeuwen [https://github.com/jteeuwen/glfw]

Plans
-----
* Better keyboard input
* Image blitting
* Termbox backend
* User supplied and non-square fonts
* World creation
* Noise generators
* Random color Blenders
* Palette support
* Merge in lighting
* More fov algorithms
* Common GUI widgets
* Curses like API
* Window resizing and fullscreen
* Website, documentation, tutorial, and more examples.
