rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit color, unicode console with built in font.
* Cross platform windowing via github.com/skelterjohn/go.wde.
* Field of view, lighting, and pathfinding algorithms.

```go
package main

import (
    "github.com/ajhager/rog"
)

func main() {
    rog.Open(48, 32, "rog")
    for rog.IsOpen() {
        rog.Set(20, 15, nil, nil, "Hello, 世界")
        if rog.Key == rog.Escape {
            rog.Close()
        }
        rog.Flush()
    }
}
```

Issues
------
* Input, drawing performance, and stability issues with go.wde.

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Thanks
------
* libTCOD [http://doryen.eptalys.net/libtcod/]
* John Asmuth [http://github.com/skelterjohn/go.wde]
* Yamamushi [http://www.theasciiproject.com/]

Plans
-----
* Input and Render interfaces to prepare for..
* Multiple back ends that can be selected at compile time
* Blitter interface for blitting consoles, images, etc. on to consoles
* User supplied and non-square fonts
* World creation
* Noise generators
* Random color Blenders
* Palette support
* Lighting
* More fov algorithms
