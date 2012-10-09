rog
===
![Rog Screenshot](http://github.com/ajhager/rog/raw/master/data/screenshot.png)

* 24bit, scaling console with custom font support
* Cross platform rendering and input
* Field of view, lighting, and pathfinding algorithms
* Procedural color and palette manipulation

[Documentation](http://go.pkgdoc.org/github.com/ajhager/rog "Documentation")

```go
package main

import (
    "github.com/ajhager/rog"
)

func main() {
    rog.Open(20, 11, 2, "rog", nil)
    for rog.Running() {
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
* You will need glfw static libs and development headers installed for now.
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Plans
-----
* Website, tutorial, and more demos
* Window resizing and fullscreen
* Better keyboard handling
* Audio generation and output
* Image (subcell) blitting
* Custom drawing callback
* Noise generators
* Merge in lighting
* World creation
* More fov algorithms
