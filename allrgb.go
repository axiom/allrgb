package allrgb

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
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

func newHilbert(rect image.Rectangle) *hilbert {
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

////////////////////////////////////////////////////////

type trivialPlacer struct {
	image.Rectangle
	count int
}

func newTrivialPlacer(rect image.Rectangle) *trivialPlacer {
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

/////////////////////////////////////////////////////////

type randomPlacer struct {
	image.Rectangle
	positions []int
	index     int
}

func newRandomPlacer(rect image.Rectangle) randomPlacer {
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

/////////////////////////////////////////////////////////

func sampleColorProducer() chan color.Color {
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
