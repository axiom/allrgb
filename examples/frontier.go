package examples

import (
	sadcolor "code.google.com/p/sadbox/color"
	"github.com/axiom/allrgb"
	"image"
	"image/color"
	"math"
	"math/rand"
)

type frontier struct {
	image.Rectangle
	frontier       map[image.Point]bool
	taken          map[image.Point]sadcolor.HSL
	occupied       []bool
	previousPoints [10]image.Point
	direction      float64
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	f.frontier = make(map[image.Point]bool)
	f.taken = make(map[image.Point]sadcolor.HSL)
	f.occupied = make([]bool, rect.Dx()*rect.Dy())
	f.Max = rect.Max

	rand.Intn(10)

	// Start at some place.
	f.extend([]image.Point{
		// Center
		{f.Max.X / 2, f.Max.Y / 2},

		// Corners
		/*
			{0, 0},
			{f.Max.X - 1, 0},
			f.Max.Sub(image.Point{1, 1}),
			{0, f.Max.Y - 1},
		*/

		// Random starting points.
		/*
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
			{rand.Intn(f.Max.X), rand.Intn(f.Max.Y)},
		*/
	})

	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	best := f.previousPoints[0]
	shortest := 10e128

	// We want to find the place with the least color differance in the frontier.
	for p, _ := range f.frontier {
		var colors []sadcolor.HSL

		// Get the colors of all the taken neighbours so we can use those for the distance calculations.
		neighbours := f.takenNeighbours(p)
		for _, neighbour := range neighbours {
			colors = append(colors, f.taken[neighbour])
		}

		distance := 0.0
		distance += 0 * colorDistance(c, colors)
		distance += -1 * neighboursCount(p, neighbours)
		distance += -1 * f.distancePrevious(p)
		distance += 0 * f.distanceDirection(p)
		if distance < shortest {
			shortest = distance
			best = p
		}
	}

	f.take(best, c)
	return best
}

// Get all the possible neighbourns of the given point.
func (f *frontier) neighbours(p image.Point) []image.Point {
	// There can be at most 8 neighbours
	neighbours := make([]image.Point, 0, 8)

	if p.X > 0 {
		neighbours = append(neighbours, p.Add(image.Point{-1, 0}))
		if p.Y > 0 {
			neighbours = append(neighbours, p.Add(image.Point{-1, -1}))
		}
		if p.Y+1 < f.Rectangle.Max.Y {
			neighbours = append(neighbours, p.Add(image.Point{-1, 1}))
		}
	}

	if p.X+1 < f.Rectangle.Max.X {
		neighbours = append(neighbours, p.Add(image.Point{1, 0}))
		if p.Y > 0 {
			neighbours = append(neighbours, p.Add(image.Point{1, -1}))
		}
		if p.Y+1 < f.Rectangle.Size().Y {
			neighbours = append(neighbours, p.Add(image.Point{1, 1}))
		}
	}

	if p.Y > 0 {
		neighbours = append(neighbours, p.Add(image.Point{0, -1}))
	}
	if p.Y+1 < f.Rectangle.Size().Y {
		neighbours = append(neighbours, p.Add(image.Point{0, 1}))
	}

	return neighbours
}

// Get only the available (unpainted) neighbours of the given point.
func (f *frontier) availableNeighbours(p image.Point) []image.Point {
	neighbours := make([]image.Point, 0, 8)
	for _, neighbour := range f.neighbours(p) {
		if !f.occupied[allrgb.PointToOffset(neighbour, f.Rectangle)] {
			neighbours = append(neighbours, neighbour)
		}
	}
	return neighbours
}

// Get only the taken (painted) neighbours of the given point.
func (f *frontier) takenNeighbours(p image.Point) []image.Point {
	neighbours := make([]image.Point, 0, 8)
	for _, neighbour := range f.neighbours(p) {
		if f.occupied[allrgb.PointToOffset(neighbour, f.Rectangle)] {
			neighbours = append(neighbours, neighbour)
		}
	}
	return neighbours
}

// Extend the frontier with more points.
func (f *frontier) extend(ps []image.Point) {
	for _, p := range ps {
		f.frontier[p] = true
	}
}

// Take some points from the frontier.
func (f *frontier) take(p image.Point, c color.Color) {
	offset := f.offset(p)

	// We need to make sure we don't place a color ontop of another, if se we try again.
	if f.occupied[offset] {
		panic("I lost the race")
	}

	f.direction = math.Atan2(float64(f.previousPoints[0].X-p.X), float64(f.previousPoints[0].Y-p.Y))
	for i := len(f.previousPoints) - 1; i > 0; i-- {
		f.previousPoints[i] = f.previousPoints[i-1]
	}
	f.previousPoints[0] = p

	f.taken[p] = sadcolor.HSLModel.Convert(c).(sadcolor.HSL)
	f.occupied[offset] = true

	f.extend(f.availableNeighbours(p))
	delete(f.frontier, p)
}

func (f *frontier) offset(p image.Point) int {
	return allrgb.PointToOffset(p, f.Rectangle)
}

func (f *frontier) distancePrevious(p image.Point) float64 {
	diff := 0.0
	for _, pp := range f.previousPoints {
		diff += math.Pow(float64(pp.X-p.X), 2) + math.Pow(float64(pp.Y-p.Y), 2)

	}
	return math.Sqrt(diff / float64(len(f.previousPoints)))
}

func (f *frontier) distanceDirection(p image.Point) float64 {
	direction := math.Atan2(float64(f.previousPoints[0].X-p.X), float64(f.previousPoints[0].Y-p.Y))
	diff := math.Abs(direction - f.direction)
	return math.Pow(diff, 2)
}

func neighboursCount(p image.Point, neighbours []image.Point) float64 {
	return math.Pow(float64(len(neighbours)), 2)
}

// Get a distance value for the differance of the given color to the slice of colors.
func colorDistance(color color.Color, colors []sadcolor.HSL) float64 {
	if len(colors) == 0 {
		return 0
	}

	c := sadcolor.HSLModel.Convert(color).(sadcolor.HSL)
	diff := 0.0
	for _, cc := range colors {
		diff += 2*math.Pow(c.H-cc.H, 2) + math.Pow(c.L-cc.L, 2) + math.Pow(c.S-cc.S, 2)
	}

	return math.Sqrt(diff / float64(len(colors)))
}
