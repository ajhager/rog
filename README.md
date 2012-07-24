rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit color, unicode console with built in font.
* Cross platform windowing via github.com/skelterjohn/go.wde.
* Field of view, lighting, and pathfinding algorithms.

```go
import (
    "github.com/ajhager/rog"
)

func main() {
    rog.Open(48, 32, "rog")
    for rog.IsOpen() {
        rog.Set(0, 0, nil, nil, "Hello, 世界")
        if rog.Key == "escape" {
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
* John Asmuth [http://github.com/skelterjohn/go.wde]
* Yamamushi [http://www.theasciiproject.com/]

Plans
-----
* Pathfinding.
* World generation.
* User supplied font sets and tilemaps.
* Noise generators.
* Image scale and blitting.
* Console to console blitting.
* Fold lighting into the library.
* Test suite.
