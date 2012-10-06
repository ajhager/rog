package glfw

import (
	"bytes"
	"github.com/ajhager/rog"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"image"
	_ "image/png"
    "runtime"
)

func init() {
	rog.SetBackend(new(glfwBackend))
}

var (
    vs = []float32{0, 0, 0, 0, 0, 0, 0, 0}
    cs = []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
    ts = []float32{0, 0, 0, 0, 0, 0, 0, 0}
)

type glfwBackend struct {
	open  bool
	mouse *rog.MouseData
	key   int
	font  image.Image
    s, t float32
	width, height, zoom int
    verts []float32
}

func (w *glfwBackend) Open(width, height, zoom int) {
    fontChan := make(chan image.Image)
    go func() {
	    font, _, err := image.Decode(bytes.NewBuffer(rog.FontData()))
	    if err != nil {
		    panic(err)
	    }
	    fontChan <- font
    }()

	if err := glfw.Init(); err != nil {
		panic(err)
	}

	w.zoom = zoom
    w.width = width
    w.height = height

	glfw.OpenWindowHint(glfw.WindowNoResize, gl.TRUE)
	err := glfw.OpenWindow(width*16*zoom, height*16*zoom, 8, 8, 8, 8, 0, 0, glfw.Windowed)
	if err != nil {
		panic(err)
	}

	glfw.SetWindowCloseCallback(func() int { w.Close(); return 0 })
	glfw.SetKeyCallback(func(key, state int) { w.setKey(key, state) })
	glfw.Enable(glfw.KeyRepeat)

	w.mouse = new(rog.MouseData)
	glfw.SetMousePosCallback(func(x, y int) { w.mouseMove(x, y) })
	glfw.SetMouseButtonCallback(func(but, state int) { w.mousePress(but, state) })

	// font, _, err := image.Decode(bytes.NewBuffer(rog.FontData()))
	// if err != nil {
	// 	panic(err)
	// }
	// w.font = font
    // w.s = 16 / float32(font.Bounds().Max.X)
    // w.t = 16 / float32(font.Bounds().Max.Y)

	fc := float32(16 * zoom)

	for y := 0; y < height; y++ {
	    for x := 0; x < width; x++ {
	        cx := float32(x) * fc
	        cy := float32(y) * fc
	        w.verts = append(w.verts, cx, cy, cx, cy + fc, cx + fc, cy + fc, cx + fc, cy)
		}
	}

    runtime.LockOSThread()
	glInit(width*16*zoom, height*16*zoom)

    w.font = <-fontChan
    m := w.font.(*image.NRGBA)
    w.s = 16 / float32(m.Bounds().Max.X)
    w.t = 16 / float32(m.Bounds().Max.Y)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, m.Bounds().Max.X, m.Bounds().Max.Y, 0, gl.RGBA, gl.UNSIGNED_BYTE, m.Pix)

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
		w.key = -1

		gl.Clear(gl.COLOR_BUFFER_BIT)

		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				fg, bg, ch := console.Get(x, y)
				w.letter(x, y, 0, bg)

				if ch != 0 && ch != 32 {
					w.letter(x, y, ch, fg)
				}
			}
		}

		glfw.SwapBuffers()
	}
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

func (w *glfwBackend) Key() int {
	return w.key
}

func (w *glfwBackend) setKey(key, state int) {
	if state == glfw.KeyPress {
		rogKey, exists := glfwToRogKey[key]
		if exists {
			w.key = rogKey
		}

		if key < 256 {
			w.key = key
		}
	}
}

func glInit(width, height int) {
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

	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.COLOR_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.VertexPointer(2, 0, vs)
	gl.ColorPointer(3, 0, cs)
	gl.TexCoordPointer(2, 0, ts)
}

// Draw a letter at a certain coordinate
func (w *glfwBackend) letter(lx, ly int, c rune, cl rog.RGB) {
    start := 8 * (ly*w.width+lx)

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

	u := float32(c % 256) * w.s
	v := float32(c / 256) * w.t
    ts[0] = u
    ts[1] = v
    ts[2] = u
    ts[3] = v+w.t
    ts[4] = u+w.s
    ts[5] = v+w.t
    ts[6] = u+w.s
    ts[7] = v

	gl.DrawArrays(gl.POLYGON, 0, 4)
}

var glfwToRogKey = map[int]int{
	glfw.KeyBackspace:  rog.Backspace,
	glfw.KeyTab:        rog.Tab,
	glfw.KeyEsc:        rog.Escape,
	glfw.KeySpace:      rog.Space,
	glfw.KeyDel:        rog.Delete,
	glfw.KeyLsuper:     rog.LSuper,
	glfw.KeyRsuper:     rog.RSuper,
	glfw.KeyLshift:     rog.LShift,
	glfw.KeyRshift:     rog.RShift,
	glfw.KeyLctrl:      rog.LControl,
	glfw.KeyRctrl:      rog.RControl,
	glfw.KeyLalt:       rog.LAlt,
	glfw.KeyRalt:       rog.RAlt,
	glfw.KeyF1:         rog.F1,
	glfw.KeyF2:         rog.F2,
	glfw.KeyF3:         rog.F3,
	glfw.KeyF4:         rog.F4,
	glfw.KeyF5:         rog.F5,
	glfw.KeyF6:         rog.F6,
	glfw.KeyF7:         rog.F7,
	glfw.KeyF8:         rog.F8,
	glfw.KeyF9:         rog.F9,
	glfw.KeyF10:        rog.F10,
	glfw.KeyF11:        rog.F11,
	glfw.KeyF12:        rog.F12,
	glfw.KeyF13:        rog.F13,
	glfw.KeyF14:        rog.F14,
	glfw.KeyF15:        rog.F15,
	glfw.KeyF16:        rog.F16,
	glfw.KeyUp:         rog.Up,
	glfw.KeyDown:       rog.Down,
	glfw.KeyLeft:       rog.Left,
	glfw.KeyRight:      rog.Right,
	glfw.KeyEnter:      rog.Return,
	glfw.KeyInsert:     rog.Insert,
	glfw.KeyHome:       rog.Home,
	glfw.KeyEnd:        rog.End,
	glfw.KeyCapslock:   rog.Capslock,
	glfw.KeyKPDivide:  rog.KPDivide,
	glfw.KeyKPMultiply: rog.KPMultiply,
	glfw.KeyKPSubtract: rog.KPSubtract,
	glfw.KeyKPAdd:      rog.KPAdd,
	glfw.KeyKPDecimal:  rog.KPDecimal,
	glfw.KeyKPEqual:    rog.KPEqual,
	glfw.KeyKPEnter:    rog.KPEnter,
	glfw.KeyKPNumlock:  rog.KPNumlock,
	glfw.KeyKP0:        rog.KP0,
	glfw.KeyKP1:        rog.KP1,
	glfw.KeyKP2:        rog.KP2,
	glfw.KeyKP3:        rog.KP3,
	glfw.KeyKP4:        rog.KP4,
	glfw.KeyKP5:        rog.KP5,
	glfw.KeyKP6:        rog.KP6,
	glfw.KeyKP7:        rog.KP7,
	glfw.KeyKP8:        rog.KP8,
	glfw.KeyKP9:        rog.KP9}
