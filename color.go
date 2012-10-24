package rog

import (
	"math"
	"math/rand"
)

func colorToFloats(c RGB) (r, g, b float64) {
	r = float64(c.R) / 255.0
	g = float64(c.G) / 255.0
	b = float64(c.B) / 255.0
	return
}

func colorToFloatsA(c RGB) (r, g, b float64) {
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
	a = clampF(0, 1, a)
	topR, topG, topB := float64(top.R), float64(top.G), float64(top.B)
	botR, botG, botB := float64(bot.R), float64(bot.G), float64(bot.B)
	return RGB{
		uint8(botR + (topR-botR)*a),
		uint8(botG + (topG-botG)*a),
		uint8(botB + (topB-botB)*a),
	}
}

// Blender interface
type Blender interface {
	Blend(RGB, int, int) RGB
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
func Hex(n uint32) *RGB {
	r := uint8((n >> 16) & 0xFF)
	g := uint8((n >> 8) & 0xFF)
	b := uint8(n & 0xFF)
	return &RGB{r, g, b}
}

// Rand returns a random RGB color
func Rand() RGB {
	return RGB{
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256)),
		uint8(rand.Int31n(256))}
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

// Scale the color a random amount.
func (c RGB) RandScale() RGB {
	return scale(c, rand.Float64())
}

// RGB Blender interface
func (c RGB) Blend(o RGB, i, t int) RGB {
	return c
}

// BlendFunc
type BlendFunc func(RGB) RGB

// BlendFunc Blender interface
func (bf BlendFunc) Blend(o RGB, i, t int) RGB {
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

func RandScale() BlendFunc {
	return func(bot RGB) RGB {
		return bot.RandScale()
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

// Scales
type ScaleFunc func(RGB, int, int) RGB

func (sf ScaleFunc) Blend(c RGB, i, t int) RGB {
	return sf(c, i, t)
}

func Discrete(blenders ...Blender) ScaleFunc {
	return func(bot RGB, i, t int) RGB {
		return blenders[i%len(blenders)].Blend(bot, i, t)
	}
}

func Linear(blenders ...Blender) ScaleFunc {
	return func(bot RGB, i, t int) RGB {
		if i == 0 {
			return blenders[0].Blend(bot, i, t)
		}

		if i == (t - 1) {
			return blenders[len(blenders)-1].Blend(bot, i, t)
		}

		a := (float64(i) / float64(t-1)) * float64(len(blenders)-1)
		b := int(math.Floor(a))
		return alpha(blenders[b+1].Blend(bot, i, t), blenders[b].Blend(bot, i, t), a-float64(b))
	}
}

var (
	Black        = RGB{0, 0, 0}
	DarkestGrey  = RGB{31, 31, 31}
	DarkerGrey   = RGB{63, 63, 63}
	DarkGrey     = RGB{95, 95, 95}
	Grey         = RGB{127, 127, 127}
	LightGrey    = RGB{159, 159, 159}
	LighterGrey  = RGB{191, 191, 191}
	LightestGrey = RGB{223, 223, 223}
	White        = RGB{255, 255, 255}

	DarkestSepia  = RGB{31, 24, 15}
	DarkerSepia   = RGB{63, 50, 31}
	DarkSepia     = RGB{94, 75, 47}
	Sepia         = RGB{127, 101, 63}
	LightSepia    = RGB{158, 134, 100}
	LighterSepia  = RGB{191, 171, 143}
	LightestSepia = RGB{222, 211, 195}

	DesaturatedRed        = RGB{127, 63, 63}
	DesaturatedFlame      = RGB{127, 79, 63}
	DesaturatedOrange     = RGB{127, 95, 63}
	DesaturatedAmber      = RGB{127, 111, 63}
	DesaturatedYellow     = RGB{127, 127, 63}
	DesaturatedLime       = RGB{111, 127, 63}
	DesaturatedChartreuse = RGB{95, 127, 63}
	DesaturatedGreen      = RGB{63, 127, 63}
	DesaturatedSea        = RGB{63, 127, 95}
	DesaturatedTurquoise  = RGB{63, 127, 111}
	DesaturatedCyan       = RGB{63, 127, 127}
	DesaturatedSky        = RGB{63, 111, 127}
	DesaturatedAzure      = RGB{63, 95, 127}
	DesaturatedBlue       = RGB{63, 63, 127}
	DesaturatedHan        = RGB{79, 63, 127}
	DesaturatedViolet     = RGB{95, 63, 127}
	DesaturatedPurple     = RGB{111, 63, 127}
	DesaturatedFuchsia    = RGB{127, 63, 127}
	DesaturatedMagenta    = RGB{127, 63, 111}
	DesaturatedPink       = RGB{127, 63, 95}
	DesaturatedCrimson    = RGB{127, 63, 79}

	LightestRed        = RGB{255, 191, 191}
	LightestFlame      = RGB{255, 207, 191}
	LightestOrange     = RGB{255, 223, 191}
	LightestAmber      = RGB{255, 239, 191}
	LightestYellow     = RGB{255, 255, 191}
	LightestLime       = RGB{239, 255, 191}
	LightestChartreuse = RGB{223, 255, 191}
	LightestGreen      = RGB{191, 255, 191}
	LightestSea        = RGB{191, 255, 223}
	LightestTurquoise  = RGB{191, 255, 239}
	LightestCyan       = RGB{191, 255, 255}
	LightestSky        = RGB{191, 239, 255}
	LightestAzure      = RGB{191, 223, 255}
	LightestBlue       = RGB{191, 191, 255}
	LightestHan        = RGB{207, 191, 255}
	LightestViolet     = RGB{223, 191, 255}
	LightestPurple     = RGB{239, 191, 255}
	LightestFuchsia    = RGB{255, 191, 255}
	LightestMagenta    = RGB{255, 191, 239}
	LightestPink       = RGB{255, 191, 223}
	LightestCrimson    = RGB{255, 191, 207}

	LighterRed        = RGB{255, 127, 127}
	LighterFlame      = RGB{255, 159, 127}
	LighterOrange     = RGB{255, 191, 127}
	LighterAmber      = RGB{255, 223, 127}
	LighterYellow     = RGB{255, 255, 127}
	LighterLime       = RGB{223, 255, 127}
	LighterChartreuse = RGB{191, 255, 127}
	LighterGreen      = RGB{127, 255, 127}
	LighterSea        = RGB{127, 255, 191}
	LighterTurquoise  = RGB{127, 255, 223}
	LighterCyan       = RGB{127, 255, 255}
	LighterSky        = RGB{127, 223, 255}
	LighterAzure      = RGB{127, 191, 255}
	LighterBlue       = RGB{127, 127, 255}
	LighterHan        = RGB{159, 127, 255}
	LighterViolet     = RGB{191, 127, 255}
	LighterPurple     = RGB{223, 127, 255}
	LighterFuchsia    = RGB{255, 127, 255}
	LighterMagenta    = RGB{255, 127, 223}
	LighterPink       = RGB{255, 127, 191}
	LighterCrimson    = RGB{255, 127, 159}

	LightRed        = RGB{255, 63, 63}
	LightFlame      = RGB{255, 111, 63}
	LightOrange     = RGB{255, 159, 63}
	LightAmber      = RGB{255, 207, 63}
	LightYellow     = RGB{255, 255, 63}
	LightLime       = RGB{207, 255, 63}
	LightChartreuse = RGB{159, 255, 63}
	LightGreen      = RGB{63, 255, 63}
	LightSea        = RGB{63, 255, 159}
	LightTurquoise  = RGB{63, 255, 207}
	LightCyan       = RGB{63, 255, 255}
	LightSky        = RGB{63, 207, 255}
	LightAzure      = RGB{63, 159, 255}
	LightBlue       = RGB{63, 63, 255}
	LightHan        = RGB{111, 63, 255}
	LightViolet     = RGB{159, 63, 255}
	LightPurple     = RGB{207, 63, 255}
	LightFuchsia    = RGB{255, 63, 255}
	LightMagenta    = RGB{255, 63, 207}
	LightPink       = RGB{255, 63, 159}
	LightCrimson    = RGB{255, 63, 111}

	Red        = RGB{255, 0, 0}
	Flame      = RGB{255, 63, 0}
	Orange     = RGB{255, 127, 0}
	Amber      = RGB{255, 191, 0}
	Yellow     = RGB{255, 255, 0}
	Lime       = RGB{191, 255, 0}
	Chartreuse = RGB{127, 255, 0}
	Green      = RGB{0, 255, 0}
	Sea        = RGB{0, 255, 127}
	Turquoise  = RGB{0, 255, 191}
	Cyan       = RGB{0, 255, 255}
	Sky        = RGB{0, 191, 255}
	Azure      = RGB{0, 127, 255}
	Blue       = RGB{0, 0, 255}
	Han        = RGB{63, 0, 255}
	Violet     = RGB{127, 0, 255}
	Purple     = RGB{191, 0, 255}
	Fuchsia    = RGB{255, 0, 255}
	Magenta    = RGB{255, 0, 191}
	Pink       = RGB{255, 0, 127}
	Crimson    = RGB{255, 0, 63}

	DarkRed        = RGB{191, 0, 0}
	DarkFlame      = RGB{191, 47, 0}
	DarkOrange     = RGB{191, 95, 0}
	DarkAmber      = RGB{191, 143, 0}
	DarkYellow     = RGB{191, 191, 0}
	DarkLime       = RGB{143, 191, 0}
	DarkChartreuse = RGB{95, 191, 0}
	DarkGreen      = RGB{0, 191, 0}
	DarkSea        = RGB{0, 191, 95}
	DarkTurquoise  = RGB{0, 191, 143}
	DarkCyan       = RGB{0, 191, 191}
	DarkSky        = RGB{0, 143, 191}
	DarkAzure      = RGB{0, 95, 191}
	DarkBlue       = RGB{0, 0, 191}
	DarkHan        = RGB{47, 0, 191}
	DarkViolet     = RGB{95, 0, 191}
	DarkPurple     = RGB{143, 0, 191}
	DarkFuchsia    = RGB{191, 0, 191}
	DarkMagenta    = RGB{191, 0, 143}
	DarkPink       = RGB{191, 0, 95}
	DarkCrimson    = RGB{191, 0, 47}

	DarkerRed        = RGB{127, 0, 0}
	DarkerFlame      = RGB{127, 31, 0}
	DarkerOrange     = RGB{127, 63, 0}
	DarkerAmber      = RGB{127, 95, 0}
	DarkerYellow     = RGB{127, 127, 0}
	DarkerLime       = RGB{95, 127, 0}
	DarkerChartreuse = RGB{63, 127, 0}
	DarkerGreen      = RGB{0, 127, 0}
	DarkerSea        = RGB{0, 127, 63}
	DarkerTurquoise  = RGB{0, 127, 95}
	DarkerCyan       = RGB{0, 127, 127}
	DarkerSky        = RGB{0, 95, 127}
	DarkerAzure      = RGB{0, 63, 127}
	DarkerBlue       = RGB{0, 0, 127}
	DarkerHan        = RGB{31, 0, 127}
	DarkerViolet     = RGB{63, 0, 127}
	DarkerPurple     = RGB{95, 0, 127}
	DarkerFuchsia    = RGB{127, 0, 127}
	DarkerMagenta    = RGB{127, 0, 95}
	DarkerPink       = RGB{127, 0, 63}
	DarkerCrimson    = RGB{127, 0, 31}

	DarkestRed        = RGB{63, 0, 0}
	DarkestFlame      = RGB{63, 15, 0}
	DarkestOrange     = RGB{63, 31, 0}
	DarkestAmber      = RGB{63, 47, 0}
	DarkestYellow     = RGB{63, 63, 0}
	DarkestLime       = RGB{47, 63, 0}
	DarkestChartreuse = RGB{31, 63, 0}
	DarkestGreen      = RGB{0, 63, 0}
	DarkestSea        = RGB{0, 63, 31}
	DarkestTurquoise  = RGB{0, 63, 47}
	DarkestCyan       = RGB{0, 63, 63}
	DarkestSky        = RGB{0, 47, 63}
	DarkestAzure      = RGB{0, 31, 63}
	DarkestBlue       = RGB{0, 0, 63}
	DarkestHan        = RGB{15, 0, 63}
	DarkestViolet     = RGB{31, 0, 63}
	DarkestPurple     = RGB{47, 0, 63}
	DarkestFuchsia    = RGB{63, 0, 63}
	DarkestMagenta    = RGB{63, 0, 47}
	DarkestPink       = RGB{63, 0, 31}
	DarkestCrimson    = RGB{63, 0, 15}

	Brass  = RGB{191, 151, 96}
	Copper = RGB{197, 136, 124}
	Gold   = RGB{229, 191, 0}
	Silver = RGB{203, 203, 203}

	Celadon = RGB{172, 255, 175}
	Peach   = RGB{255, 159, 127}
)