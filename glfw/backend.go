package glfw

import (
	"github.com/ajhager/rog"
    "github.com/jteeuwen/glfw"
    "github.com/banthar/gl"
    "bytes"
    "image"
    "image/color"
    "image/draw"
)

func Backend() rog.Backend {
	return new(glfwBackend)
}

type glfwBackend struct {
	open         bool
	mouse        *rog.MouseData
	key          string
	font   image.Image
    zoom int
}

func (w *glfwBackend) Open(width, height, zoom int) {
    if err := glfw.Init(); err != nil {
        panic(err)
    }

    w.zoom = zoom

    glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)
    err := glfw.OpenWindow(width*16*zoom, height*16*zoom, 8, 8, 8, 8, 0, 0, glfw.Windowed)
    if err != nil {
        panic(err)
    }

    glfw.SetWindowCloseCallback(func() int {w.Close(); return 0;})
    glfw.SetKeyCallback(func(key, state int) {w.setKey(key, state)})
    glfw.Enable(glfw.KeyRepeat)

	w.mouse = new(rog.MouseData)
    glfw.SetMousePosCallback(func(x, y int) {w.mouseMove(x, y)})
    glfw.SetMouseButtonCallback(func(but, state int) {w.mousePress(but, state)})

	font, _, err := image.Decode(bytes.NewBuffer(rog.FontData()))
	if err != nil {
		panic(err)
	}
    w.font = font

    glInit(width*16*zoom, height*16*zoom, w.font)

	w.open = true
}

func (w *glfwBackend) IsOpen() bool {
	return w.open && glfw.WindowParam(glfw.Opened) == 1
}

func (w *glfwBackend) Close() {
	w.open = false
    glfw.CloseWindow()
    glfw.Terminate()
}

func (w *glfwBackend) Name(title string) {
    glfw.SetWindowTitle(title)
}

func (w *glfwBackend) Render(console *rog.Console) {
	if w.IsOpen() {
	    w.mouse.Left.Released = false
	    w.mouse.Right.Released = false
	    w.mouse.Middle.Released = false
        w.key = ""

        gl.Clear(gl.COLOR_BUFFER_BIT)

        for y := 0; y < console.Height(); y++ {
            for x := 0; x < console.Width(); x++ {
				fg, bg, ch := console.Get(x, y)

                setColor(bg)
                w.letter(x, y, 0)
            
                if (ch != 0 && ch != 32)  {
                    setColor(fg)
                    w.letter(x, y, ch)
                }
            }
        }

        glfw.SwapBuffers()
	}
}

func (w *glfwBackend) Screen() image.Image {
	return &image.Uniform{rog.RGB{255, 255, 255}}
}

func (w *glfwBackend) Mouse() *rog.MouseData {
	return w.mouse
}

func (w *glfwBackend) mouseMove(x, y int) {
    w.mouse.Pos.X = x
    w.mouse.Pos.Y = y
    w.mouse.Cell.X = x / (16 * w.zoom)
    w.mouse.Cell.Y = y / (16 * w.zoom)
}

func (w *glfwBackend) mousePress(button, state int) {
    switch state {
    case glfw.KeyPress:
        switch button {
	    case glfw.MouseLeft:
			w.mouse.Left.Pressed = true
		case glfw.MouseRight:
			w.mouse.Right.Pressed = true
		case glfw.MouseMiddle:
			w.mouse.Middle.Pressed = true
        }
    case glfw.KeyRelease:
		switch button {
		case glfw.MouseLeft:
			w.mouse.Left.Pressed = false
			w.mouse.Left.Released = true
		case glfw.MouseRight:
			w.mouse.Right.Pressed = false
			w.mouse.Right.Released = true
		case glfw.MouseMiddle:
			w.mouse.Middle.Pressed = false
			w.mouse.Middle.Released = true
		}
    }
}

func (w *glfwBackend) Key() string {
	return w.key
}

func (w *glfwBackend) setKey(key, state int) {
    if state == glfw.KeyPress {
        rogKey, exists := glfwToRogKey[key]
        if exists {
            w.key = rogKey
        }
    }
    // KeyPageup
    // KeyPagedown
    // KeyKP0
    // KeyKP1
    // KeyKP2
    // KeyKP3
    // KeyKP4
    // KeyKP5
    // KeyKP6
    // KeyKP7
    // KeyKP8
    // KeyKP9
    // KeyKPDidivde
    // KeyKPMultiply
    // KeyKPSubtract
    // KeyKPAdd
    // KeyKPDecimal
    // KeyKPEqual
    // KeyKPEnter
    // KeyKPNumlock
    // KeyScrolllock
    // KeyPause
    // KeyMenu
    // KeyLast = KeyMenu
}

func glInit(width, height int, font image.Image) {
    gl.Init()
    gl.Enable(gl.TEXTURE_2D)
    gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
    gl.Viewport(0, 0, width, height)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    textures := make([]gl.Texture, 1)
    gl.GenTextures(textures)
    textures[0].Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

    b := font.Bounds()
    m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
    draw.Draw(m, m.Bounds(), font, b.Min, draw.Src)

    gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, m.Bounds().Max.X, m.Bounds().Max.Y, 0, gl.RGBA, gl.UNSIGNED_BYTE, m.Pix)
}

// Draw a letter at a certain coordinate
func (w *glfwBackend) letter(lx, ly int, c rune) {
    b := w.font.Bounds()
    fc := float32(16*w.zoom)
    cx := float32(lx) * fc
    cy := float32(ly) * fc
    verts := []float32{cx, cy, cx, cy+fc, cx+fc, cy+fc, cx+fc, cy, cx, cy}
    x := float32(c % 256)
    y := float32(c / 256)
    s := float32(16) / float32(b.Max.X)
    t := float32(16) / float32(b.Max.Y)
    u := x * s
    v := y * t
    gl.EnableClientState(gl.VERTEX_ARRAY)
    gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
    gl.VertexPointer(2, 0, verts)
    gl.TexCoordPointer(2, 0, []float32{u, v, u, v+t, u+s, v+t, u+s, v, u, v})
    gl.DrawArrays(gl.POLYGON, 0, len(verts)/2-1)
    gl.DisableClientState(gl.VERTEX_ARRAY)
    gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
}

// Set the opengl drawing color
func setColor(c color.Color) {
    r, g, b, _ := c.RGBA()
    gl.Color3ub(uint8(r), uint8(g), uint8(b))
}


var glfwToRogKey map[int]string = map[int]string {
    glfw.KeySpace: rog.Space,
    glfw.KeyEsc: rog.Escape,
    glfw.KeyF1: rog.F1,
    glfw.KeyF2: rog.F2,
    glfw.KeyF3: rog.F3,
    glfw.KeyF4: rog.F4,
    glfw.KeyF5: rog.F5,
    glfw.KeyF6: rog.F6,
    glfw.KeyF7: rog.F7,
    glfw.KeyF8: rog.F8,
    glfw.KeyF9: rog.F9,
    glfw.KeyF10: rog.F10,
    glfw.KeyF11: rog.F11,
    glfw.KeyF12: rog.F12,
    glfw.KeyF13: rog.F13,
    glfw.KeyF14: rog.F14,
    glfw.KeyF15: rog.F15,
    glfw.KeyF16: rog.F16,
    glfw.KeyUp: rog.Up,
    glfw.KeyDown: rog.Down,
    glfw.KeyLeft: rog.Left,
    glfw.KeyRight: rog.Right,
    glfw.KeyLshift: rog.LShift,
    glfw.KeyRshift: rog.RShift,
    glfw.KeyLctrl: rog.LControl,
    glfw.KeyRctrl: rog.RControl,
    glfw.KeyLalt: rog.LAlt,
    glfw.KeyRalt: rog.RAlt,
    glfw.KeyLsuper: rog.LSuper,
    glfw.KeyRsuper: rog.RSuper,
    glfw.KeyTab: rog.Tab,
    glfw.KeyEnter: rog.PadEnter,
    glfw.KeyBackspace: rog.Backspace,
    glfw.KeyInsert: rog.Insert,
    glfw.KeyDel: rog.Delete,
    glfw.KeyHome: rog.Home,
    glfw.KeyEnd: rog.End,
    glfw.KeyCapslock: rog.CapsLock,
    48: rog.N0,
    49: rog.N1,
    50: rog.N2,
    51: rog.N3,
    52: rog.N4,
    53: rog.N5,
    54: rog.N6,
    55: rog.N7,
    56: rog.N8,
    57: rog.N9,
    65: rog.A,
    66: rog.B,
    67: rog.C,
    68: rog.D,
    69: rog.E,
    70: rog.F,
    71: rog.G,
    72: rog.H,
    73: rog.I,
    74: rog.J,
    75: rog.K,
    76: rog.L,
    77: rog.M,
    78: rog.N,
    79: rog.O,
    80: rog.P,
    81: rog.Q,
    82: rog.R,
    83: rog.S,
    84: rog.T,
    85: rog.U,
    86: rog.V,
    87: rog.W,
    88: rog.X,
    89: rog.Y,
    90: rog.Z}
