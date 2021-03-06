package allrgb

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"time"
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

// Given a point and a rectangle return the offset representing that point.
func PointToOffset(p image.Point, rect image.Rectangle) int {
	return p.X + p.Y*rect.Dx()
}

// Given an offset and a rectangle return the point representing that offset.
func OffsetToPoint(offset int, rect image.Rectangle) image.Point {
	return image.Point{
		X: offset % rect.Dx(),
		Y: offset / rect.Dx(),
	}
}

////////////////////////////////////////////////////////

func ColorDetermined(rect image.Rectangle, cp ColorProducer, p Placer) image.Image {
	img := image.NewRGBA(rect)
	tick := timer(1000, rect.Dx()*rect.Dy())
	for c := range cp.Produce() {
		p := p.Place(c)
		img.Set(p.X, p.Y, c)
		tick()
	}
	return img
}

func PlaceDetermined(rect image.Rectangle, c Colorer, pp PlaceProducer) image.Image {
	img := image.NewRGBA(rect)
	tick := timer(1000, rect.Dx()*rect.Dy())
	for p := range pp.Produce() {
		c := c.Color(p)
		img.Set(p.X, p.Y, c)
		tick()
	}
	return img
}

func ColorDeterminedFrameSaver(rect image.Rectangle, cp ColorProducer, p Placer, rate int, name string) error {
	img := image.NewRGBA(rect)
	frame := 0
	tick := timer(rate, rect.Dx()*rect.Dy())
	for c := range cp.Produce() {
		p := p.Place(c)
		img.Set(p.X, p.Y, c)

		tick()
		if frame%rate == 0 {
			if err := SaveImage(fmt.Sprintf("%v-%05d.png", name, frame), img); err != nil {
				return err
			}
		}
		frame++
	}

	return SaveImage(fmt.Sprintf("%v-%05d.png", name, frame), img)
}

func SaveImage(name string, img image.Image) error {
	f, err := os.Create(name)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func timer(rate int, totalFrames int) func() {
	frame := 0
	previousFrame := 0
	timestamp := time.Now()

	return func() {
		if frame%rate == 0 || frame == totalFrames {
			now := time.Now()
			timediff := now.Sub(timestamp)
			framediff := frame - previousFrame
			fps := 0.0

			if timediff.Nanoseconds() > 0 {
				fps = float64(framediff) / timediff.Seconds()
			}

			eta := time.Duration(float64(totalFrames-frame)/fps) * time.Second

			fmt.Printf(
				"% 3.0f%% % 5.0f fps ETA %v\n",
				float64(100*frame)/float64(totalFrames),
				fps,
				eta,
			)

			previousFrame = frame
			timestamp = now
		}

		frame++
	}
}
