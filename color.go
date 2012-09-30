package rog

import (
	"math"
)

var (
	Black RGB = RGB{0, 0, 0}
	White RGB = RGB{255, 255, 255}
)

func colorToFloats(c RGB) (r, g, b float64) {
	r = float64(c.R) / 255.0
	g = float64(c.G) / 255.0
	b = float64(c.B) / 255.0
	return
}

func clampF(low, high, value float64) float64 {
	return math.Min(high, math.Max(low, value))
}

func multiply(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(topR * botR * 255),
		uint8(topG * botG * 255),
		uint8(topB * botB * 255),
	}
}

func dodgeF(top, bot float64) (out uint8) {
	if bot != 1 {
		out = uint8(255 * clampF(0, 1, top/(1-bot)))
	} else {
		out = uint8(255)
	}
	return
}

func dodge(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		dodgeF(topR, botR),
		dodgeF(topG, botG),
		dodgeF(topB, botB),
	}
}

func screen(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * (1 - ((1 - topR) * (1 - botR)))),
		uint8(255 * (1 - ((1 - topG) * (1 - botG)))),
		uint8(255 * (1 - ((1 - topB) * (1 - botB)))),
	}
}

func overlayF(top, bot float64) (out uint8) {
	if bot < 0.5 {
		out = uint8(2 * top * bot * 255)
	} else {
		out = uint8(255 * (1 - 2*(1-top)*(1-bot)))
	}
	return
}

func overlay(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		overlayF(topR, botR),
		overlayF(topG, botG),
		overlayF(topB, botB),
	}
}

func lighten(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * math.Max(topR, botR)),
		uint8(255 * math.Max(topG, botG)),
		uint8(255 * math.Max(topB, botB)),
	}
}

func darken(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * math.Min(topR, botR)),
		uint8(255 * math.Min(topG, botG)),
		uint8(255 * math.Min(topB, botB)),
	}
}

func burn(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR+topR-1)),
		uint8(255 * clampF(0, 1, botG+topG-1)),
		uint8(255 * clampF(0, 1, botB+topB-1)),
	}
}

func scale(bot RGB, s float64) RGB {
	botr, botg, botb := colorToFloats(bot)
	return RGB{
		uint8(255 * botr * s),
		uint8(255 * botg * s),
		uint8(255 * botb * s),
	}
}

func add(top, bot RGB) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR+topR)),
		uint8(255 * clampF(0, 1, botG+topG)),
		uint8(255 * clampF(0, 1, botB+topB)),
	}
}

func addAlpha(top, bot RGB, a float64) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR*a+topR)),
		uint8(255 * clampF(0, 1, botG*a+topG)),
		uint8(255 * clampF(0, 1, botB*a+topB)),
	}
}

func alpha(top, bot RGB, a float64) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * (botR + (topR-botR)*a)),
		uint8(255 * (botG + (topG-botG)*a)),
		uint8(255 * (botB + (topB-botB)*a)),
	}
}

// Blender interface
type Blender interface {
	Blend(RGB) RGB
}

// RGB represents a traditional 24-bit alpha-premultiplied color, having 8 bits for each of red, green, and blue.
type RGB struct {
	R, G, B uint8
}

func (c RGB) RGBA() (r, g, b, a uint32) {
	r = uint32(c.R)
	r |= r << 8
	g = uint32(c.G)
	g |= g << 8
	b = uint32(c.B)
	b |= b << 8
	a = 0xffff
	return
}

// HEX returns parses a uint32 into RGB components.
func Hex(n uint32) RGB {
	r := uint8((n >> 16) & 0xFF)
	g := uint8((n >> 8) & 0xFF)
	b := uint8(n & 0xFF)
	return RGB{r, g, b}
}

// Multiply = old * new
func (c RGB) Multiply(o RGB) RGB {
	return multiply(o, c)
}

// Dodge = new / (white - old)
func (c RGB) Dodge(o RGB) RGB {
	return dodge(o, c)
}

// Screen = white - (white - old) * (white - new)
func (c RGB) Screen(o RGB) RGB {
	return screen(o, c)
}

// Overlay = new.x <= 0.5 ? 2*new*old : white - 2*(white-new)*(white-old)
func (c RGB) Overlay(o RGB) RGB {
	return overlay(o, c)
}

// Darken = MIN(old, new)
func (c RGB) Darken(o RGB) RGB {
	return darken(o, c)
}

// Lighten = MIN(old, new)
func (c RGB) Lighten(o RGB) RGB {
	return lighten(o, c)
}

// Burn = old + new - white
func (c RGB) Burn(o RGB) RGB {
	return burn(o, c)
}

// Scale = old * s
func (c RGB) Scale(s float64) RGB {
	return scale(c, s)
}

// Add = old + new
func (c RGB) Add(o RGB) RGB {
	return add(o, c)
}

// AddAlpha = old + alpha*new
func (c RGB) AddAlpha(o RGB, a float64) RGB {
	return addAlpha(o, c, a)
}

// Alpha = (1-alpha)*old + alpha*(new-old)
func (c RGB) Alpha(o RGB, a float64) RGB {
	return alpha(o, c, a)
}

// RGB Blender interface
func (c RGB) Blend(o RGB) RGB {
	return c
}

// BlendFunc
type BlendFunc func(RGB) RGB

// BlendFunc Blender interface
func (bf BlendFunc) Blend(o RGB) RGB {
	return bf(o)
}

func Multiply(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return multiply(top, bot)
	}
}

func Dodge(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return dodge(top, bot)
	}
}

func Screen(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return screen(top, bot)
	}
}

func Overlay(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return overlay(top, bot)
	}
}

func Lighten(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return lighten(top, bot)
	}
}

func Darken(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return darken(top, bot)
	}
}

func Burn(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return burn(top, bot)
	}
}

func Scale(s float64) BlendFunc {
	return func(bot RGB) RGB {
		return scale(bot, s)
	}
}

func Add(top RGB) BlendFunc {
	return func(bot RGB) RGB {
		return add(top, bot)
	}
}

func AddAlpha(top RGB, a float64) BlendFunc {
	return func(bot RGB) RGB {
		return addAlpha(top, bot, a)
	}
}

func Alpha(top RGB, a float64) BlendFunc {
	return func(bot RGB) RGB {
		return addAlpha(top, bot, a)
	}
}
