package examples

import (
	"image"
	"image/color"
	"math/rand"
)

type randomPlacer struct {
	image.Rectangle
	positions []int
	index     int
}

func NewRandomPlacer(rect image.Rectangle) randomPlacer {
	rp := randomPlacer{}
	rp.positions = make([]int, rect.Dx()*rect.Dy())
	rp.Max = rect.Max

	// Initiate positions with all possible positions
	for i := 0; i < len(rp.positions); i++ {
		rp.positions[i] = i
	}

	// Shuffle the available positions for instance randomness
	for i := 0; i < len(rp.positions); i++ {
		j := rand.Intn(len(rp.positions))
		k := rand.Intn(len(rp.positions))
		rp.positions[j], rp.positions[k] = rp.positions[k], rp.positions[j]
	}

	return rp
}

func (rp *randomPlacer) Place(_ color.Color) image.Point {
	position := rp.positions[rp.index]
	rp.index++
	return image.Point{
		X: position / rp.Dy(),
		Y: position % rp.Dy(),
	}
}

func RandomColorProducer() chan color.Color {
	nextColor := make(chan color.Color)
	go func() {
		colors := make([]color.Color, 32*32*32)
		i := 0
		for r := 0; r < 32; r++ {
			for g := 0; g < 32; g++ {
				for b := 0; b < 32; b++ {
					colors[i] = color.RGBA{
						R: uint8(r << 3 & 0xf8),
						G: uint8(g << 3 & 0xf8),
						B: uint8(b << 3 & 0xf8),
						A: 255,
					}
					i++
				}
			}
		}

		for _ = range colors {
			i, j := rand.Intn(len(colors)), rand.Intn(len(colors))
			colors[i], colors[j] = colors[j], colors[i]
		}

		for _, color := range colors {
			nextColor <- color
		}
		close(nextColor)
	}()
	return nextColor
}
