package main

import (
	"fmt"
	"github.com/axiom/allrgb"
	"github.com/axiom/allrgb/examples"
	"image"
)

func main() {
	rect := image.Rectangle{Max: image.Point{X: 256, Y: 128}}

	configurations := map[string]image.Image{
		"trivial": allrgb.ColorDetermined(
			rect,
			allrgb.ColorProducerFunc(examples.SampleColorProducer),
			examples.NewTrivialPlacer(rect)),
		"hilbert": allrgb.ColorDetermined(
			rect,
			allrgb.ColorProducerFunc(examples.SampleColorProducer),
			examples.NewHilbertPlacer(rect)),
	}

	for name, img := range configurations {
		if err := allrgb.SaveImage(name, img); err != nil {
			fmt.Printf("Could not do %v: %v\n", name, err)
		}
	}
}
