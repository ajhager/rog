package rog

import (
	"image/color"
	"math"
)

func colorEq(x, y color.Color) bool {
	xr, xg, xb, xa := x.RGBA()
	yr, yg, yb, ya := y.RGBA()
	return xr == yr && xg == yg && xb == yb && xa == ya
}

func colorToFloats(c color.Color) (rr, gg, bb, aa float64) {
	const M = float64(1<<16 - 1)
	r, g, b, a := c.RGBA()
	rr = float64(r) / M
	gg = float64(g) / M
	bb = float64(b) / M
	aa = float64(a) / M
	return
}

func overlay(top, bot float64) (out uint8) {
	if bot < 0.5 {
		out = uint8(2 * top * bot * 255)
	} else {
		out = uint8(255 * (1 - 2*(1-top)*(1-bot)))
	}
	return
}

func dodge(top, bot float64) (out uint8) {
	if bot != 1 {
		out = uint8(255 * clamp(0, 1, top / (1 - bot)))
	} else {
		out = uint8(255)
	}
	return
}

func clamp(low, high, value float64) float64 {
	return math.Min(high, math.Max(low, value))
}

type ColorBlend func(color.Color, color.Color) color.Color

func Normal(top, bot color.Color) color.Color {
	return top
}

func Multiply(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(topR * botR * 255),
		uint8(topG * botG * 255),
		uint8(topB * botB * 255),
		uint8(topA * botA * 255),
	}
}

func Screen(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * (1 - ((1 - topR) * (1 - botR)))),
		uint8(255 * (1 - ((1 - topG) * (1 - botG)))),
		uint8(255 * (1 - ((1 - topB) * (1 - botB)))),
		uint8(255 * (1 - ((1 - topA) * (1 - botA)))),
	}
}

func Overlay(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		overlay(topR, botR),
		overlay(topG, botG),
		overlay(topB, botB),
		overlay(topA, botA),
	}
}

func Lighten(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * math.Max(topR, botR)),
		uint8(255 * math.Max(topG, botG)),
		uint8(255 * math.Max(topB, botB)),
		uint8(255 * math.Max(topA, botA)),
	}
}

func Darken(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * math.Min(topR, botR)),
		uint8(255 * math.Min(topG, botG)),
		uint8(255 * math.Min(topB, botB)),
		uint8(255 * math.Min(topA, botA)),
	}
}

func Burn(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * clamp(0, 1, botR + topR - 1)),
		uint8(255 * clamp(0, 1, botG + topG - 1)),
		uint8(255 * clamp(0, 1, botB + topB - 1)),
		uint8(255 * clamp(0, 1, botA + topA - 1)),
	}
}

func Dodge(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		dodge(topR, botR),
		dodge(topG, botG),
		dodge(topB, botB),
		dodge(topA, botA),
	}
}

func Add(top, bot color.Color) color.Color {
	topR, topG, topB, topA := colorToFloats(top)
	botR, botG, botB, botA := colorToFloats(bot)
	return color.RGBA{
		uint8(255 * clamp(0, 1, botR + topR)),
		uint8(255 * clamp(0, 1, botG + topG)),
		uint8(255 * clamp(0, 1, botB + topB)),
		uint8(255 * clamp(0, 1, botA + topA)),
	}
}

func AddAlpha(a float64) ColorBlend {
	return func(top, bot color.Color) color.Color {
		topR, topG, topB, topA := colorToFloats(top)
		botR, botG, botB, botA := colorToFloats(bot)
		return color.RGBA{
			uint8(255 * clamp(0, 1, botR * a + topR)),
			uint8(255 * clamp(0, 1, botG * a + topG)),
			uint8(255 * clamp(0, 1, botB * a + topB)),
			uint8(255 * clamp(0, 1, botA * a + topA)),
		}
	}
}

func Alpha(a float64) ColorBlend {
	return func(top, bot color.Color) color.Color {
		topR, topG, topB, topA := colorToFloats(top)
		botR, botG, botB, botA := colorToFloats(bot)
		return color.RGBA{
			uint8(255 * (botR + (topR - botR) * a)),
			uint8(255 * (botG + (topG - botG) * a)),
			uint8(255 * (botB + (topB - botB) * a)),
			uint8(255 * (botA + (topA - botA) * a)),
		}
	}
}
