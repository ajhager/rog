rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit color, unicode console with built in font.
* Cross platform with pluggable backends.
* Field of view, lighting, and pathfinding algorithms.
* Procedural color functions.

```go
package main

import (
    "github.com/ajhager/rog"
    "github.com/ajhager/rog/wde"
)

func main() {
    rog.Open(48, 32, "rog", wde.Backend())
    for rog.IsOpen() {
        rog.Set(20, 15, nil, nil, "Hello, 世界")
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
* Blitter interface for blitting consoles, images, etc. on to consoles
* User supplied and non-square fonts
* World creation
* Noise generators
* Random color Blenders
* Palette support
* Merge in lighting
* More fov algorithms
* Website, documentation, tutorial, and more examples.
