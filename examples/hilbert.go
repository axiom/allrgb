package examples

import (
	"image"
	"image/color"
)

type hilbert struct {
	next chan image.Point
}

func rot(n, x, y, rx, ry int) (int, int) {
	x_ := x
	y_ := y

	if ry == 0 {
		if rx == 1 {
			x_ = n - 1 - x
			y_ = n - 1 - y
		}

		x_, y_ = y_, x_
	}
	return x_, y_
}

func NewHilbertPlacer(rect image.Rectangle) *hilbert {
	h := hilbert{
		next: make(chan image.Point),
	}
	go func() {
		n := rect.Max.X * rect.Max.Y
		for d := 0; d <= n; d++ {
			var x, y, rx, ry, t int
			t = d
			for s := 1; s < n; s *= 2 {
				rx = 1 & (t / 2)
				ry = 1 & (t ^ rx)
				x, y = rot(s, x, y, rx, ry)
				x += s * rx
				y += s * ry
				t /= 4
			}
			h.next <- image.Point{X: x, Y: y}
		}
		close(h.next)
	}()
	return &h
}

func (h hilbert) Place(c color.Color) image.Point {
	return <-h.next
}
