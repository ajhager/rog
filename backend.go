package rog

import (
	gl "github.com/chsc/gogl/gl21"
	"github.com/go-gl/glfw"
	"image"
	"image/draw"
	_ "image/png"
	"runtime"
)

const (
    NOKEY      = -1
	Esc        = glfw.KeyEsc
	F1         = glfw.KeyF1
	F2         = glfw.KeyF2
	F3         = glfw.KeyF3
	F4         = glfw.KeyF4
	F5         = glfw.KeyF5
	F6         = glfw.KeyF6
	F7         = glfw.KeyF7
	F8         = glfw.KeyF8
	F9         = glfw.KeyF9
	F10        = glfw.KeyF10
	F11        = glfw.KeyF11
	F12        = glfw.KeyF12
	F13        = glfw.KeyF13
	F14        = glfw.KeyF14
	F15        = glfw.KeyF15
	F16        = glfw.KeyF16
	F17        = glfw.KeyF17
	F18        = glfw.KeyF18
	F19        = glfw.KeyF19
	F20        = glfw.KeyF20
	F21        = glfw.KeyF21
	F22        = glfw.KeyF22
	F23        = glfw.KeyF23
	F24        = glfw.KeyF24
	F25        = glfw.KeyF25
	Up         = glfw.KeyUp
	Down       = glfw.KeyDown
	Left       = glfw.KeyLeft
	Right      = glfw.KeyRight
	Lshift     = glfw.KeyLshift
	Rshift     = glfw.KeyRshift
	Lctrl      = glfw.KeyLctrl
	Rctrl      = glfw.KeyRctrl
	Lalt       = glfw.KeyLalt
	Ralt       = glfw.KeyRalt
	Tab        = glfw.KeyTab
	Enter      = glfw.KeyEnter
	Backspace  = glfw.KeyBackspace
	Insert     = glfw.KeyInsert
	Del        = glfw.KeyDel
	Pageup     = glfw.KeyPageup
	Pagedown   = glfw.KeyPagedown
	Home       = glfw.KeyHome
	End        = glfw.KeyEnd
	KP0        = glfw.KeyKP0
	KP1        = glfw.KeyKP1
	KP2        = glfw.KeyKP2
	KP3        = glfw.KeyKP3
	KP4        = glfw.KeyKP4
	KP5        = glfw.KeyKP5
	KP6        = glfw.KeyKP6
	KP7        = glfw.KeyKP7
	KP8        = glfw.KeyKP8
	KP9        = glfw.KeyKP9
	KPDivide   = glfw.KeyKPDivide
	KPMultiply = glfw.KeyKPMultiply
	KPSubtract = glfw.KeyKPSubtract
	KPAdd      = glfw.KeyKPAdd
	KPDecimal  = glfw.KeyKPDecimal
	KPEqual    = glfw.KeyKPEqual
	KPEnter    = glfw.KeyKPEnter
	KPNumlock  = glfw.KeyKPNumlock
	Capslock   = glfw.KeyCapslock
	Scrolllock = glfw.KeyScrolllock
	Pause      = glfw.KeyPause
	Lsuper     = glfw.KeyLsuper
	Rsuper     = glfw.KeyRsuper
	Menu       = glfw.KeyMenu
)

var (
	vs       = []float32{0, 0, 0, 0, 0, 0, 0, 0}
	cs       = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	ts       = []float32{0, 0, 0, 0, 0, 0, 0, 0}
	textures []gl.Uint
)

type glfwBackend struct {
	open                bool
	mouse               *MouseData
	font                *FontData
	key                 int
	s, t                float32
	width, height, zoom int
	verts               []float32
}

func (w *glfwBackend) Open(width, height, zoom int, fs bool, font *FontData) {
	if err := glfw.Init(); err != nil {
		panic(err)
	}

	w.zoom = zoom
	w.width = width
	w.height = height

    var fwidth = width*font.CellWidth*zoom
    var fheight = height*font.CellHeight*zoom
    var twidth = fwidth
    var theight = fheight

    flag := glfw.Windowed
    if (fs) {
        flag = glfw.Fullscreen
        dm := glfw.DesktopMode()
        twidth = dm.W
        theight = dm.H
    }
    
	glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)
	err := glfw.OpenWindow(twidth, theight, 8, 8, 8, 8, 0, 0, flag)
	if err != nil {
		panic(err)
	}

    w.key = NOKEY
	glfw.SetWindowCloseCallback(func() int { w.Close(); return 0 })
	glfw.SetKeyCallback(func(key, state int) { w.setKey(key, state) })
	glfw.SetCharCallback(func(key, state int) { w.setKey(key, state) })
	glfw.Enable(glfw.KeyRepeat)

	w.mouse = new(MouseData)
    glfw.Enable(glfw.MouseCursor)
	glfw.SetMousePosCallback(func(x, y int) { w.mouseMove(x, y) })
	glfw.SetMouseButtonCallback(func(but, state int) { w.mousePress(but, state) })
    glfw.Enable(glfw.MouseCursor)

    xoff := float32(twidth - fwidth) / 2.0
    yoff := float32(theight - fheight) / 2.0

	fc := float32(font.CellWidth * zoom)
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			cx := xoff + float32(x) * fc
			cy := yoff + float32(y) * fc
			w.verts = append(w.verts, cx, cy, cx, cy+fc, cx+fc, cy+fc, cx+fc, cy)
		}
	}

	runtime.LockOSThread()
	glInit(twidth, theight)

	m := font.Image.(*image.RGBA)
	w.s = float32(font.Width) / float32(m.Bounds().Max.X)
	w.t = float32(font.Height) / float32(m.Bounds().Max.Y)
	textures = make([]gl.Uint, 2)
	gl.GenTextures(2, &textures[0])

	gl.BindTexture(gl.TEXTURE_2D, textures[0])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.Sizei(m.Bounds().Max.X), gl.Sizei(m.Bounds().Max.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&m.Pix[0]))

	m = image.NewRGBA(image.Rect(0, 0, font.Width, font.Height))
	draw.Draw(m, m.Bounds(), &image.Uniform{White}, image.ZP, draw.Src)
	gl.BindTexture(gl.TEXTURE_2D, textures[1])
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.Sizei(m.Bounds().Max.X), gl.Sizei(m.Bounds().Max.Y), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Pointer(&m.Pix[0]))

	w.font = font

	w.open = true
}

func (w *glfwBackend) Running() bool {
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

func (w *glfwBackend) Render(console *Console) {
	if w.Running() {
		w.mouse.Left.Released = false
		w.mouse.Right.Released = false
		w.mouse.Middle.Released = false
		w.key = NOKEY

		gl.Clear(gl.COLOR_BUFFER_BIT)

		gl.BindTexture(gl.TEXTURE_2D, textures[1])
		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				_, bg, _ := console.Get(x, y)
				w.letter(x, y, 0, bg)
			}
		}

		gl.BindTexture(gl.TEXTURE_2D, textures[0])
		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				fg, _, ch := console.Get(x, y)
				if position, ok := w.font.Map(ch); ok {
					w.letter(x, y, position, fg)
				}
			}
		}

		glfw.SwapBuffers()
	}
}

func (w *glfwBackend) Mouse() *MouseData {
	return w.mouse
}

func (w *glfwBackend) Cursor(on bool) {
    if on {
        glfw.Enable(glfw.MouseCursor)
    } else {
        glfw.Disable(glfw.MouseCursor)
    }
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

func (w *glfwBackend) Key() int {
	return w.key
}

func (w *glfwBackend) setKey(key, state int) {
	if state == glfw.KeyPress {
		w.key = key
	}
}

func glInit(width, height int) {
	gl.Init()
	gl.Enable(gl.TEXTURE_2D)
	gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, gl.Double(width), gl.Double(height), 0, -1, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.VertexPointer(2, gl.FLOAT, 0, gl.Pointer(&vs[0]))
	gl.ColorPointer(3, gl.UNSIGNED_BYTE, 0, gl.Pointer(&cs[0]))
	gl.TexCoordPointer(2, gl.FLOAT, 0, gl.Pointer(&ts[0]))
}

// Draw a letter at a certain coordinate
func (w *glfwBackend) letter(lx, ly int, c int, cl RGB) {
	start := 8 * (ly*w.width + lx)

	vs[0] = w.verts[start]
	vs[1] = w.verts[start+1]
	vs[2] = w.verts[start+2]
	vs[3] = w.verts[start+3]
	vs[4] = w.verts[start+4]
	vs[5] = w.verts[start+5]
	vs[6] = w.verts[start+6]
	vs[7] = w.verts[start+7]

	cs[0] = cl.R
	cs[1] = cl.G
	cs[2] = cl.B
	cs[3] = cl.R
	cs[4] = cl.G
	cs[5] = cl.B
	cs[6] = cl.R
	cs[7] = cl.G
	cs[8] = cl.B
	cs[9] = cl.R
	cs[10] = cl.G
	cs[11] = cl.B

	u := float32(c%w.font.Width) * w.s
	v := float32(c/w.font.Height) * w.t
	ts[0] = u
	ts[1] = v
	ts[2] = u
	ts[3] = v + w.t
	ts[4] = u + w.s
	ts[5] = v + w.t
	ts[6] = u + w.s
	ts[7] = v

	gl.DrawArrays(gl.POLYGON, 0, 4)
}
