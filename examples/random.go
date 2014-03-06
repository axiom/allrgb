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
