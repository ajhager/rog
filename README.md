rog A roguelike game library written in go
===
[Documentation](http://hagerbot.com/rog/docs.html "Documentation")

![Rog Screenshot](http://hagerbot.com/img/screenshot_rog_fov.png)

Setup
-----
rog currently depends on [github.com/go-gl/glfw](http://github.com/go-gl/glfw). You can skip this step if you already have that up and running.
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
