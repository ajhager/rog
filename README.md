rog A roguelike game library written in go
===
![Rog Screenshot](http://hagerbot.com/img/screenshot_rog_fov.png)

* 24bit, scaling console with custom font support
* Cross platform rendering and input
* Field of view, lighting, and pathfinding algorithms
* Procedural color and palette manipulation

[Documentation](http://hagerbot.com/rog/docs.html "Documentation")

```go
package main

import (
    "hagerbot.com/rog"
)

func main() {
    rog.Open(20, 11, 2, false, "rog", nil)
    for rog.Running() {
        rog.Set(5, 5, nil, nil, "Hello, 世界!")
        if rog.Key() == rog.Esc {
            rog.Close()
        }
        rog.Flush()
    }
}
```

Setup
-----
* Ubuntu: apt-get install libglfw-dev
* OSX: brew install glfw
* Windows: download the glfw binaries, then drop the GL/ directory into C:\MingW\include and the files for your arch under libmingw into C:\MingW\lib. You will then need to install glfw.dll system wide or have it in the directory with your game.

Install
-------
go get hagerbot.com/rog

Notes
-----
* On Windows you can build your project with `go build -ldflags -Hwindowsgui` to inhibit the console window that pops up by default.

Plans
-----
* Tutorial for installing dependencies
* Generalize rendering to 2d tile engine and build rog semantics on top
* Blitting of image.Image interface to consoles
* User defined OpenGL callback
* Noise generators (Perlin, Fractal, Simplex, Wavelet, etc.)
* Clean up map system, merge in lighting, and add BSP level gen
* Screenshot method
* 1:2 custom font demo
* Host and link to the full 4096x4096 Unicode font texture
* Create tween module and unify it with color scales
