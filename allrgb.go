package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
)

// Placer determines where a color should go.
type Placer interface {
	Place(color.Color) image.Point
}

type PlacerFunc func(color.Color) image.Point

func (pf PlacerFunc) Place(c color.Color) image.Point {
	return pf(c)
}

type Producer interface {
	Produce() chan color.Color
}

type ProducerFunc func() chan color.Color

func (pf ProducerFunc) Produce() chan color.Color {
	return pf()
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

func newHilbert(rect image.Rectangle) hilbert {
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
	return h
}

func (h hilbert) Place(c color.Color) image.Point {
	return <-h.next
}

////////////////////////////////////////////////////////

type trivialPlacer struct {
	image.Rectangle
	count int
}

func newTrivialPlacer(rect image.Rectangle) trivialPlacer {
	tp := trivialPlacer{}
	tp.Max = rect.Max
	return tp
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

func produceColors() chan color.Color {
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

func main() {
	bits := 15
	rect := image.Rectangle{Max: image.Point{X: 256, Y: 128}}
	img := image.NewRGBA(rect)
	fmt.Println(bits, rect)

	f, err := os.Create("allrgb.png")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	nextColor := ProducerFunc(produceColors).Produce()

	//h := newHilbert(rect)
	h := newTrivialPlacer(rect)
	//h := newRandomPlacer(rect)
	for c := range nextColor {
		p := h.Place(c)
		img.Set(p.X, p.Y, c)
	}

	if err = png.Encode(f, img); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
