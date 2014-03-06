package examples

import (
	"image"
	"image/color"
)

type trivialPlacer struct {
	image.Rectangle
	count int
}

func NewTrivialPlacer(rect image.Rectangle) *trivialPlacer {
	tp := trivialPlacer{}
	tp.Max = rect.Max
	return &tp
}

func (tp *trivialPlacer) Place(_ color.Color) image.Point {
	p := image.Point{
		X: tp.count / tp.Dy(),
		Y: tp.count % tp.Dy(),
	}
	tp.count++
	return p
}

func SampleColorProducer() chan color.Color {
	nextColor := make(chan color.Color)
	go func() {
		for r := 0; r < 32; r++ {
			for g := 0; g < 32; g++ {
				for b := 0; b < 32; b++ {
					nextColor <- color.RGBA{
						R: uint8(r << 3 & 0xf8),
						G: uint8(g << 3 & 0xf8),
						B: uint8(b << 3 & 0xf8),
						A: 255}
				}
			}
		}
		close(nextColor)
	}()
	return nextColor
}
