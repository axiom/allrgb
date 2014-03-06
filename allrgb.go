package allrgb

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

//////////////////////////////////////////////////////////
// Color determined

// Placer determines where a color should go.
type Placer interface {
	Place(color.Color) image.Point
}

type PlacerFunc func(color.Color) image.Point

func (pf PlacerFunc) Place(c color.Color) image.Point {
	return pf(c)
}

// Used to produce a sequence of colors to be placed.
type ColorProducer interface {
	Produce() chan color.Color
}

type ColorProducerFunc func() chan color.Color

func (pf ColorProducerFunc) Produce() chan color.Color {
	return pf()
}

////////////////////////////////////////////////////////
// Position determined

// Colorer determines the color of a given point.
type Colorer interface {
	Color(image.Point) color.Color
}

type ColorerFunc func(image.Point) color.Color

func (cf ColorerFunc) Color(p image.Point) color.Color {
	return cf(p)
}

// Used to produce a sequence of postions to be painted.
type PlaceProducer interface {
	Produce() chan image.Point
}

type PlaceProducerFunc func() chan image.Point

func (ppf PlaceProducerFunc) Produce() chan image.Point {
	return ppf()
}

////////////////////////////////////////////////////////

func ColorDetermined(rect image.Rectangle, cp ColorProducer, p Placer) image.Image {
	img := image.NewRGBA(rect)
	for c := range cp.Produce() {
		p := p.Place(c)
		img.Set(p.X, p.Y, c)
	}
	return img
}

func PlaceDetermined(rect image.Rectangle, c Colorer, pp PlaceProducer) image.Image {
	img := image.NewRGBA(rect)
	for p := range pp.Produce() {
		c := c.Color(p)
		img.Set(p.X, p.Y, c)
	}
	return img
}

func SaveImage(name string, img image.Image) error {
	f, err := os.Create(fmt.Sprintf("%v.png", name))
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}
