// Submitted by Dustin Lacewell <dlacewell@gmail.com>
package main

import (
	"github.com/iand/perlin"
	"hagerbot.com/rog"
	"math"
)

const (
	width  = 40
	height = 40
)

var (
	seed    = 1.0
	alpha   = 1.08
	beta    = 0.0
	octaves = 6

	alpha_lerp = pingpong(0.008, 0.0, 0.6, 1.5)
	beta_lerp  = pingpong(0.005, 0.0, 0.0, 1.0)

	r_lerp = pingpong(0.01, 0.0, 0.0, 255.0)
	g_lerp = pingpong(0.02, 128.0, 255.0, 0.0)
	b_lerp = pingpong(0.04, 0.0, 255.0, 0.0)

	base_x = 100.0
	base_y = 100.0
)

func pingpong(delta, start, min, max float64) func() float64 {
	var (
		lerp_t   = math.Max(min, math.Min(max, start))
		lerp_d   = delta
		lerp_min = min
		lerp_max = max
	)

	lerp := func(t, a, b float64) float64 {
		return a + t*(b-a)
	}

	return func() float64 {
		lerp_t += lerp_d
		if lerp_t >= 1.0 || lerp_t <= 0.0 {
			lerp_d *= -1.0
		}
		lerp_t = math.Min(lerp_t, 1.0)
		lerp_t = math.Max(lerp_t, 0.0)
		return lerp(lerp_t, lerp_min, lerp_max)
	}
}

func render() {
	base_x += 1.0
	base_y += 1.0

	alpha = alpha_lerp()
	beta = beta_lerp()

	r := uint8(r_lerp())
	g := uint8(g_lerp())
	b := uint8(b_lerp())

	for x := 0.0; x < width; x++ {
		for y := 0.0; y < height; y++ {
			px := base_x + x
			py := base_y + y

			perlin := perlin.Noise2D(px, py, int64(seed), alpha, beta, octaves)
			perlin = (perlin * .5) + .5
			perlin = math.Max(0.0, perlin)
			perlin = math.Min(1.0, perlin)

			color := rog.RGB{r, g, b}
			result := rog.Black.Alpha(color, perlin)

			rog.Set(int(x), int(y)+1, nil, result, "")
		}
	}

}

func main() {
	rog.Open(width, height+2, 1, false, "Perlin-noise Test", nil)
	for rog.Running() {
		render()
		rog.Set(0, height+1, nil, nil, "%v", rog.Fps())
		if rog.Key() == rog.Esc {
			rog.Close()
		}
		rog.Flush()
	}
}
