package rog

func Viewable(grid Grid2d, x, y int) []struct{ X, Y int } {
	visible := make([]struct{ X, Y int }, 0)
	visible = append(visible, struct{ X, Y int }{20, 20})
	return visible
}
