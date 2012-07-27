package rog

import (
	"image/color"
	"math"
)

func colorToFloats(c color.Color) (rr, gg, bb float64) {
	const M = float64(1<<16 - 1)
	r, g, b, _ := c.RGBA()
	rr = float64(r) / M
	gg = float64(g) / M
	bb = float64(b) / M
	return
}

func clampF(low, high, value float64) float64 {
	return math.Min(high, math.Max(low, value))
}

func multiply(top, bot color.Color) RGB {
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

func dodge(top, bot color.Color) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		dodgeF(topR, botR),
		dodgeF(topG, botG),
		dodgeF(topB, botB),
	}
}

func screen(top, bot color.Color) RGB {
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

func overlay(top, bot color.Color) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		overlayF(topR, botR),
		overlayF(topG, botG),
		overlayF(topB, botB),
	}
}

func lighten(top, bot color.Color) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * math.Max(topR, botR)),
		uint8(255 * math.Max(topG, botG)),
		uint8(255 * math.Max(topB, botB)),
	}
}

func darken(top, bot color.Color) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * math.Min(topR, botR)),
		uint8(255 * math.Min(topG, botG)),
		uint8(255 * math.Min(topB, botB)),
	}
}

func burn(top, bot color.Color) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR+topR-1)),
		uint8(255 * clampF(0, 1, botG+topG-1)),
		uint8(255 * clampF(0, 1, botB+topB-1)),
	}
}

func scale(bot color.Color, s float64) RGB {
	botr, botg, botb := colorToFloats(bot)
	return RGB{
		uint8(255 * botr * s),
		uint8(255 * botg * s),
		uint8(255 * botb * s),
	}
}

func add(top, bot color.Color) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR+topR)),
		uint8(255 * clampF(0, 1, botG+topG)),
		uint8(255 * clampF(0, 1, botB+topB)),
	}
}

func addAlpha(top, bot color.Color, a float64) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR*a+topR)),
		uint8(255 * clampF(0, 1, botG*a+topG)),
		uint8(255 * clampF(0, 1, botB*a+topB)),
	}
}

func alpha(top, bot color.Color, a float64) RGB {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * (botR + (topR-botR)*a)),
		uint8(255 * (botG + (topG-botG)*a)),
		uint8(255 * (botB + (topB-botB)*a)),
	}
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
func (c RGB) Multiply(o color.Color) RGB {
	return multiply(o, c)
}

// Dodge = new / (white - old)
func (c RGB) Dodge(o color.Color) RGB {
	return dodge(o, c)
}

// Screen = white - (white - old) * (white - new)
func (c RGB) Screen(o color.Color) RGB {
	return screen(o, c)
}

// Overlay = new.x <= 0.5 ? 2*new*old : white - 2*(white-new)*(white-old)
func (c RGB) Overlay(o color.Color) RGB {
	return overlay(o, c)
}

// Darken = MIN(old, new)
func (c RGB) Darken(o color.Color) RGB {
	return darken(o, c)
}

// Lighten = MIN(old, new)
func (c RGB) Lighten(o color.Color) RGB {
	return lighten(o, c)
}

// Burn = old + new - white
func (c RGB) Burn(o color.Color) RGB {
	return burn(o, c)
}

// Scale = old * s
func (c RGB) Scale(s float64) RGB {
	return scale(c, s)
}

// Add = old + new
func (c RGB) Add(o color.Color) RGB {
	return add(o, c)
}

// AddAlpha = old + alpha*new
func (c RGB) AddAlpha(o color.Color, a float64) RGB {
	return addAlpha(o, c, a)
}

// Alpha = (1-alpha)*old + alpha*(new-old)
func (c RGB) Alpha(o color.Color, a float64) RGB {
	return alpha(o, c, a)
}

// Blender
type Blender func(color.Color) color.Color

func Multiply(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return multiply(top, bot)
	}
}

func Dodge(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return dodge(top, bot)
	}
}

func Screen(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return screen(top, bot)
	}
}

func Overlay(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return overlay(top, bot)
	}
}

func Lighten(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return lighten(top, bot)
	}
}

func Darken(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return darken(top, bot)
	}
}

func Burn(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return burn(top, bot)
	}
}

func Scale(s float64) Blender {
	return func(bot color.Color) color.Color {
		return scale(bot, s)
	}
}

func Add(top color.Color) Blender {
	return func(bot color.Color) color.Color {
		return add(top, bot)
	}
}

func AddAlpha(top color.Color, a float64) Blender {
	return func(bot color.Color) color.Color {
		return addAlpha(top, bot, a)
	}
}

func Alpha(top color.Color, a float64) Blender {
	return func(bot color.Color) color.Color {
		return addAlpha(top, bot, a)
	}
}
