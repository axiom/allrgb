package examples

import (
	"image"
	"image/color"
)

type frontier struct {
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	p := image.Point{}
	return p
}
