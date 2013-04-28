rog A roguelike game library written in go
===

Note
----
rog is no longer maintained. Please check out [eng](http://github.com/ajhager/eng) for an up to date 2d game library for go.

Setup
-----
rog depends on [github.com/go-gl/glfw](http://github.com/go-gl/glfw). You can skip this step if you already have that up and running.
* Ubuntu: apt-get install libglfw-dev
* OSX: brew install glfw
* Windows: download the glfw binaries, then drop the GL directory into C:\MinGW\include and the files for your arch under libmingw into C:\MinGW\lib. You will then need to install glfw.dll system wide or have it in the directory with your game.

Install
-------
`go get hagerbot.com/rog`

Try it!
-------
```go
package main

import (
    "github.com/ajhager/rog"
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
