package glfw

import (
	"bytes"
	"github.com/ajhager/rog"
	"github.com/banthar/gl"
	"github.com/jteeuwen/glfw"
	"image"
	"image/color"
	"image/draw"
	_ "image/png"
)

func init() {
	rog.SetBackend(new(glfwBackend))
}

type glfwBackend struct {
	open  bool
	mouse *rog.MouseData
	key   int
	font  image.Image
	zoom  int
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

	glfw.SetWindowCloseCallback(func() int { w.Close(); return 0 })
	glfw.SetKeyCallback(func(key, state int) { w.setKey(key, state) })
	glfw.Enable(glfw.KeyRepeat)

	w.mouse = new(rog.MouseData)
	glfw.SetMousePosCallback(func(x, y int) { w.mouseMove(x, y) })
	glfw.SetMouseButtonCallback(func(but, state int) { w.mousePress(but, state) })

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
		w.key = -1

		gl.Clear(gl.COLOR_BUFFER_BIT)

		for y := 0; y < console.Height(); y++ {
			for x := 0; x < console.Width(); x++ {
				fg, bg, ch := console.Get(x, y)

				setColor(bg)
				w.letter(x, y, 0)

				if ch != 0 && ch != 32 {
					setColor(fg)
					w.letter(x, y, ch)
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
	fc := float32(16 * w.zoom)
	cx := float32(lx) * fc
	cy := float32(ly) * fc
	verts := []float32{cx, cy, cx, cy + fc, cx + fc, cy + fc, cx + fc, cy, cx, cy}
	x := float32(c % 256)
	y := float32(c / 256)
	s := float32(16) / float32(b.Max.X)
	t := float32(16) / float32(b.Max.Y)
	u := x * s
	v := y * t
	gl.EnableClientState(gl.VERTEX_ARRAY)
	gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
	gl.VertexPointer(2, 0, verts)
	gl.TexCoordPointer(2, 0, []float32{u, v, u, v + t, u + s, v + t, u + s, v, u, v})
	gl.DrawArrays(gl.POLYGON, 0, len(verts)/2-1)
	gl.DisableClientState(gl.VERTEX_ARRAY)
	gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
}

// Set the opengl drawing color
func setColor(c color.Color) {
	r, g, b, _ := c.RGBA()
	gl.Color3ub(uint8(r), uint8(g), uint8(b))
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
