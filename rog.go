package rog

import (
    "time"
    "github.com/jteeuwen/glfw"
    "github.com/banthar/gl"
)

var (
    root *Console
    texture gl.Texture
    running bool
    winColor Color = Color{0, 0, 0}
    fadeColor Color = Color{0, 0, 0}
    fadeAlpha uint8 = 0
    cellSize int
    newTime, oldTime, frames int64
    frameTime, Dt float64
    showFPS bool = false
)

func Clear() {
    setWinColor(Color{0, 0, 0})
    root.Clear()
}

func Fill(ch uint8, bg, fg Color) {
    root.Fill(ch, bg, fg)
}

func FillCh(ch uint8) {
    root.FillCh(ch)
}

func FillBg(bg Color) {
    root.FillBg(bg)
    setWinColor(bg)
}

func FillFg(fg Color) {
    root.FillFg(fg)
}

func Set(x, y int, ch uint8, bg, fg Color) {
    root.Set(x, y, ch, bg, fg)
}

func SetCh(x, y int, ch uint8) {
    root.SetCh(x, y, ch)
}

func SetBg(x, y int, bg Color) {
    root.SetBg(x, y, bg)
}

func SetFg(x, y int, fg Color) {
    root.SetFg(x, y, fg)
}

func SetChFg(x, y int, ch uint8, fg Color) {
    root.SetChFg(x, y, ch, fg)
}

func SetChBg(x, y int, ch uint8, bg Color) {
    root.SetChBg(x, y, ch, bg)
}

func SetFade(fade Color, alpha uint8) {
    fadeColor = fade
    fadeAlpha = alpha
}

func Key(key int) bool {
    return glfw.Key(key) == glfw.KeyPress
}

func setWinColor(win Color) {
    r := gl.GLclampf(float32(win.R) / 255.0)
    g := gl.GLclampf(float32(win.G) / 255.0)
    b := gl.GLclampf(float32(win.B) / 255.0)
    a := gl.GLclampf(1.0)
    gl.ClearColor(r, g, b, a)
    winColor = win
}

func Open(width, height, zoom int, title string, fullscreen bool) {
    // Set up the root console
    root = NewConsole(width, height)

    // Initialize the windowing library
    if err := glfw.Init(); err != nil {
        panic(err)
    }
    running = true

    state := glfw.Windowed
    if fullscreen {
        state = glfw.Fullscreen
        glfw.SetSwapInterval(1)
    }
    cellSize = fontSize * zoom

    glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)
    err := glfw.OpenWindow(width*cellSize, height*cellSize, 8, 8, 8, 8, 0, 0, state)
    if err != nil {
        panic(err)
    }

    glfw.SetWindowCloseCallback(Close)

    //glfw.SetWindowTitle(title)

    // Initialize opengl
    glInit(width*cellSize, height*cellSize)
    setWinColor(winColor)

    // Initialize timing
    newTime, oldTime = time.Now().UnixNano(), time.Now().UnixNano()
}

func IsOpen() bool {
    return running && glfw.WindowParam(glfw.Opened) == 1
}

func Close() int {
    running = false
    glfw.CloseWindow()
    glfw.Terminate()
    return 0
}

func ShowFPS(on bool) {
    showFPS = on
}

func Flush() {
    gl.Clear(gl.COLOR_BUFFER_BIT)

    for y := 0; y < root.h; y++ {
        for x := 0; x < root.w; x++ {
            bg := root.bg[y][x]
            if !bg.Equal(winColor) {
                setColor(root.bg[y][x])
                letter(x, y, 219)
            }
            
            ch := root.ch[y][x]
            fg := root.fg[y][x]
            if (ch != 0 && ch != 32) && !fg.Equal(bg) {
                setColor(root.fg[y][x])
                letter(x, y, ch)
            }
        }
    }
    
    if (fadeAlpha != 0) {
//        setColorA(fadeColor.R, fadeColor.G, fadeColor.B, fadeAlpha)
//        letter(0, 0, root.w *cellSize, root.h*cellSize, 219)
    }

    glfw.SwapBuffers()

    frames += 1
    newTime = time.Now().UnixNano()
    Dt = float64(newTime - oldTime) / 1e9
    frameTime += Dt
    oldTime = newTime
    
    if frameTime >= 1 {
        if (showFPS) {
            println(frames)
        }
        frameTime -= 1
        frames = 0
    }
}
