package rog

import (
    "fmt"
    "math"
)

func addo(a, b byte) byte {
    c := int(a) + int(b)
    if (c > 255) {
        return 255
    }
    return byte(c)
}

func subo(a, b byte) byte {
    c := int(a) - int(b)
    if (c < 0) {
        return 0
    }
    return byte(c)
}

type Color struct {
    R, G, B uint8
}

func (c Color) String() string {
    return fmt.Sprintf("Color{%d, %d, %d}", c.R, c.G, c.B)
}

func (c Color) Equal(o Color) bool {
    return c.R == o.R && c.G == o.G && c.B == o.B
}

func (c Color) Add(o Color) Color {
    return Color{addo(c.R, o.R), addo(c.G, o.G), addo(c.B, o.B)}
}

func (c Color) Sub(o Color) Color {
    return Color{subo(c.R, o.R), subo(c.G, o.G), subo(c.B, o.B)}
}

func FromHSV(h, s, v float64) (r Color) {
    h60 := h / 60.0
    hi := int(math.Floor(h60)) % 6
    f := h60 - math.Floor(h60)
    p := v * (1.0 - s)
    q := v * (1.0 - s * f)
    t := v * (1.0 - s * (1.0 - f))
    switch hi {
    case 0:
        r.R = byte(v * 255.0 + 0.5)
        r.G = byte(t * 255.0 + 0.5)
        r.B = byte(p * 255.0 + 0.5)
    case 1:
        r.R = byte(q * 255.0 + 0.5)
        r.G = byte(v * 255.0 + 0.5)
        r.B = byte(p * 255.0 + 0.5)
    case 2:
        r.R = byte(p * 255.0 + 0.5)
        r.G = byte(v * 255.0 + 0.5)
        r.B = byte(t * 255.0 + 0.5)
    case 3:
        r.R = byte(p * 255.0 + 0.5)
        r.G = byte(q * 255.0 + 0.5)
        r.B = byte(v * 255.0 + 0.5)
    case 4:
        r.R = byte(t * 255.0 + 0.5)
        r.G = byte(p * 255.0 + 0.5)
        r.B = byte(v * 255.0 + 0.5)
    case 5:
        r.R = byte(v * 255.0 + 0.5)
        r.G = byte(p * 255.0 + 0.5)
        r.B = byte(q * 255.0 + 0.5)
    }
    return
}

func HgrMono(v float64) Color {
    return FromHSV(0, 0, v)
}

func HgrYellow(v float64) Color {
    return FromHSV(45, .6, v)
}

func HgrOrange(v float64) Color {
    return FromHSV(25, .6, v)
}

func HgrRed(v float64) Color {
    return FromHSV(5, .6, v)
}

func HgrMagenta(v float64) Color {
    return FromHSV(345, .6, v)
}

func HgrPurple(v float64) Color {
    return FromHSV(295, .45, v)
}

func HgrBlue(v float64) Color {
    return FromHSV(245, .3, v)
}

func HgrCyan(v float64) Color {
    return FromHSV(195, .3, v)
}

func HgrCeladon(v float64) Color {
    return FromHSV(145, .3, v)
}

func HgrGreen(v float64) Color {
    return FromHSV(95, .3, v)
}

var (
)
