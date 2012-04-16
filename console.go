package rog

type Console struct {
    bg, fg [][]Color
    ch [][]byte
    w, h int
}

func NewConsole(width, height int) *Console {
    bg := make([][]Color, height)
    fg := make([][]Color, height)
    ch := make([][]uint8, height)

    for y := 0; y < height; y++ {
        bg[y] = make([]Color, width)
        fg[y] = make([]Color, width)
        ch[y] = make([]uint8, width)
    }

    con := &Console{bg, fg, ch, width, height}
    con.Clear()
    return con
}

func (con *Console) Clear() {
    for x := 0; x < con.w; x++ {
        for y:= 0; y < con.h; y++ {
            con.Set(x, y, 0, Color{0, 0, 0}, Color{255, 255, 255})
        }
    }
}

func (con *Console) Fill(ch byte, bg, fg Color) {
    for x := 0; x < con.w; x++ {
        for y:= 0; y < con.h; y++ {
            con.Set(x, y, ch, bg, fg)
        }
    }
}

func (con *Console) FillCh(ch byte) {
    for x := 0; x < con.w; x++ {
        for y:= 0; y < con.h; y++ {
            con.SetCh(x, y, ch)
        }
    }
}

func (con *Console) FillBg(bg Color) {
    for x := 0; x < con.w; x++ {
        for y:= 0; y < con.h; y++ {
            con.SetBg(x, y, bg)
        }
    }
}

func (con *Console) FillFg(fg Color) {
    for x := 0; x < con.w; x++ {
        for y:= 0; y < con.h; y++ {
            con.SetFg(x, y, fg)
        }
    }
}

func (con *Console) Set(x, y int, ch byte, bg, fg Color) {
    con.bg[y][x] = bg
    con.fg[y][x] = fg
    con.ch[y][x] = ch
}

func (con *Console) SetCh(x, y int, ch byte) {
    con.ch[y][x] = ch
}

func (con *Console) SetBg(x, y int, bg Color) {
    con.bg[y][x] = bg
}

func (con *Console) SetFg(x, y int, fg Color) {
    con.fg[y][x] = fg
}

func (con *Console) SetChFg(x, y int, ch byte, fg Color) {
    con.ch[y][x] = ch
    con.fg[y][x] = fg
}

func (con *Console) SetChBg(x, y int, ch byte, bg Color) {
    con.ch[y][x] = ch
    con.bg[y][x] = bg
}

func (con *Console) Width() int {
    return con.w
}

func (con *Console) Height() int {
    return con.h
}
