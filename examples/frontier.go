package examples

import (
	sadcolor "code.google.com/p/sadbox/color"
	"github.com/axiom/allrgb"
	"image"
	"image/color"
	"math"
	"sync"
)

type queueItem struct {
	result chan image.Point
	color  color.Color
}

type frontier struct {
	sync.Mutex
	image.Rectangle
	frontier      map[image.Point]bool
	taken         map[image.Point]sadcolor.HSL
	occupied      []bool
	placeQueue    chan queueItem
	previousPoint image.Point
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	f.frontier = make(map[image.Point]bool)
	f.taken = make(map[image.Point]sadcolor.HSL)
	f.occupied = make([]bool, rect.Dx()*rect.Dy())
	f.placeQueue = make(chan queueItem)
	f.Max = rect.Max

	// Start at some place.
	start := image.Point{rect.Dx() / 2, rect.Dy() / 2}
	f.previousPoint = start
	f.extend([]image.Point{start})

	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	best := f.previousPoint
	shortest := 10e128

	// We want to find the place with the least color differance in the frontier.
	for p, _ := range f.frontier {
		var colors []sadcolor.HSL

		// Get the colors of all the taken neighbours so we can use those for the distance calculations.
		for _, neighbour := range f.takenNeighbours(p) {
			colors = append(colors, f.taken[neighbour])
		}

		if distance := colorDistance(c, colors) + f.distancePrevious(p); distance < shortest {
			shortest = distance
			best = p
		}
	}

	f.take(best, c)
	return best
}

// Get all the possible neighbourns of the given point.
func (f frontier) neighbours(p image.Point) []image.Point {
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
func (f frontier) availableNeighbours(p image.Point) []image.Point {
	neighbours := make([]image.Point, 0, 8)
	for _, neighbour := range f.neighbours(p) {
		if !f.occupied[allrgb.PointToOffset(neighbour, f.Rectangle)] {
			neighbours = append(neighbours, neighbour)
		}
	}
	return neighbours
}

// Get only the taken (painted) neighbours of the given point.
func (f frontier) takenNeighbours(p image.Point) []image.Point {
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

	f.previousPoint = p

	f.taken[p] = sadcolor.HSLModel.Convert(c).(sadcolor.HSL)
	f.occupied[offset] = true

	f.extend(f.availableNeighbours(p))
	delete(f.frontier, p)
}

func (f frontier) distancePrevious(p image.Point) float64 {
	dx := f.previousPoint.X - p.X
	dy := f.previousPoint.Y - p.Y
	return 3 * math.Sqrt(math.Pow(float64(dx), 2)+math.Pow(float64(dy), 2))
}

func (f *frontier) offset(p image.Point) int {
	return allrgb.PointToOffset(p, f.Rectangle)
}

// Get a distance value for the differance of the given color to the slice of colors.
func colorDistance(color color.Color, colors []sadcolor.HSL) float64 {
	c := sadcolor.HSLModel.Convert(color).(sadcolor.HSL)
	diff := 0.0
	for _, cc := range colors {
		diff += 3*math.Pow(c.H-cc.H, 2) + 2*math.Pow(c.L-cc.L, 2) + 1*math.Pow(c.S-cc.S, 2)
	}

	return 100 * math.Sqrt(diff) / float64(len(colors))
}
