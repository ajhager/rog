package rog

func New(width, height int, path string, cellW, cellH int) *Console {
	return &Console{
		Buffer: NewBuffer(width, height),
		font:   NewFont(path, cellW, cellH),
	}
}
