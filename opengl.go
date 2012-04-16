package rog

import (
    "github.com/banthar/gl"
)

func glInit(width, height int) {
    gl.Init()
    gl.Enable(gl.TEXTURE_2D)
    gl.Viewport(0, 0, width, height)
    gl.MatrixMode(gl.PROJECTION)
    gl.LoadIdentity()
    gl.Ortho(0, float64(width), float64(height), 0, -1, 1)
    gl.MatrixMode(gl.MODELVIEW)
    gl.LoadIdentity()
    gl.PixelStorei(gl.UNPACK_ALIGNMENT, 1)
    gl.Enable(gl.BLEND)
    gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)
    textures := make([]gl.Texture, 1)
    gl.GenTextures(textures)
    texture = textures[0]
    texture.Bind(gl.TEXTURE_2D)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)

    image := font()

    gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, image.Bounds().Max.X, image.Bounds().Max.Y, 0, gl.RGBA, gl.UNSIGNED_BYTE, image.Pix)
}

// Draw a letter at a certain coordinate
func letter(lx, ly int, c uint8) {
    fc := float32(cellSize)
    cx := float32(lx) * fc
    cy := float32(ly) * fc
    verts := []float32{cx, cy, cx, cy+fc, cx+fc, cy+fc, cx+fc, cy, cx, cy}
    y := float32(c / 16)
    x := float32(c % 16)
    t := float32(8) / float32(128)
    u := x * t
    v := y * t
    gl.EnableClientState(gl.VERTEX_ARRAY)
    gl.EnableClientState(gl.TEXTURE_COORD_ARRAY)
    gl.VertexPointer(2, 0, verts)
    gl.TexCoordPointer(2, 0, []float32{u, v, u, v+t, u+t, v+t, u+t, v, u, v})
    gl.DrawArrays(gl.POLYGON, 0, len(verts)/2-1)
    gl.DisableClientState(gl.VERTEX_ARRAY)
    gl.DisableClientState(gl.TEXTURE_COORD_ARRAY)
}

// Set the opengl drawing color
func setColor(c Color) {
    gl.Color3ub(c.R, c.G, c.B)
}

func setColorA(c Color, a uint8) {
    gl.Color4ub(c.R, c.G, c.B, a)
}
