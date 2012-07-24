package rog

import (
	"image/color"
	"math"
)

// Utility
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

func multiply(top, bot color.Color) color.Color {
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
		out = uint8(255 * clampF(0, 1, top / (1 - bot)))
	} else {
		out = uint8(255)
	}
	return
}

func dodge(top, bot color.Color) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		dodgeF(topR, botR),
		dodgeF(topG, botG),
		dodgeF(topB, botB),
	}
}

func screen(top, bot color.Color) color.Color {
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

func overlay(top, bot color.Color) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		overlayF(topR, botR),
		overlayF(topG, botG),
		overlayF(topB, botB),
	}
}

func lighten(top, bot color.Color) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * math.Max(topR, botR)),
		uint8(255 * math.Max(topG, botG)),
		uint8(255 * math.Max(topB, botB)),
	}
}

func darken(top, bot color.Color) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * math.Min(topR, botR)),
		uint8(255 * math.Min(topG, botG)),
		uint8(255 * math.Min(topB, botB)),
	}
}

func burn(top, bot color.Color) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR + topR - 1)),
		uint8(255 * clampF(0, 1, botG + topG - 1)),
		uint8(255 * clampF(0, 1, botB + topB - 1)),
	}
}

func scale(bot color.Color, s float64) color.Color {
	botr, botg, botb := colorToFloats(bot)
    return RGB{
        uint8(255 * botr * s),
        uint8(255 * botg * s),
        uint8(255 * botb * s),
    }
}

func add(top, bot color.Color) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR + topR)),
		uint8(255 * clampF(0, 1, botG + topG)),
		uint8(255 * clampF(0, 1, botB + topB)),
	}
}

func addAlpha(top, bot color.Color, a float64) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * clampF(0, 1, botR * a + topR)),
		uint8(255 * clampF(0, 1, botG * a + topG)),
		uint8(255 * clampF(0, 1, botB * a + topB)),
	}
}

func alpha(top, bot color.Color, a float64) color.Color {
	topR, topG, topB := colorToFloats(top)
	botR, botG, botB := colorToFloats(bot)
	return RGB{
		uint8(255 * (botR + (topR - botR) * a)),
		uint8(255 * (botG + (topG - botG) * a)),
		uint8(255 * (botB + (topB - botB) * a)),
	}
}

// RGB
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

func HEX(n uint32) RGB {
    r := uint8((n >> 16) & 0xFF)
    g := uint8((n >> 8) & 0xFF)
    b := uint8(n & 0xFF)
    return RGB{r, g, b}
}

func (c RGB) Multiply(o color.Color) color.Color {
    return multiply(o, c)
}

func (c RGB) Dodge(o color.Color) color.Color {
    return dodge(o, c)
}

func (c RGB) Screen(o color.Color) color.Color {
    return screen(o, c)
}

func (c RGB) Overlay(o color.Color) color.Color {
    return overlay(o, c)
}

func (c RGB) Darken(o color.Color) color.Color {
    return darken(o, c)
}

func (c RGB) Burn(o color.Color) color.Color {
    return burn(o, c)
}

func (c RGB) Scale(s float64) color.Color {
    return scale(c, s)
}

func (c RGB) Add(o color.Color) color.Color {
    return add(o, c)
}

func (c RGB) AddAlpha(o color.Color, a float64) color.Color {
    return addAlpha(o, c, a)
}

func (c RGB) Alpha(o color.Color, a float64) color.Color {
    return alpha(o, c, a)
}

// Blenders
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
