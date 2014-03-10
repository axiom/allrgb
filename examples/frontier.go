package examples

import (
	//"github.com/axiom/allrgb"
	"image"
	"image/color"
	"math"
	"sync"
)

type candidate struct {
	point image.Point
	cost  float64
}

type frontier struct {
	image.Rectangle
	frontier map[image.Point]bool
	canvas   map[image.Point]color.Color
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	f.frontier = make(map[image.Point]bool, 32*32*32)
	f.canvas = make(map[image.Point]color.Color, 32*32*32)
	f.Max = rect.Max

	// Start at some place.
	f.extend([]image.Point{{f.Max.X / 2, f.Max.Y / 2}})

	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	// Create a chan that will produce all the current frontier points
	frontier := make(chan image.Point)
	go func() {
		for fp, _ := range f.frontier {
			frontier <- fp
		}
		close(frontier)
	}()

	// Workers will send possible candidates from the frontier points together
	// with their cost on this chan.
	candidates := make(chan candidate)
	var wg sync.WaitGroup
	for i := 0; i < 7; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for p := range frontier {
				//colorCost := f.colorDistance(p, c)
				neighbourhoodCost := f.neighbourhoodCost(p)
				//distanceCost := f.distanceCost(p)
				//cityDistanceCost := f.cityDistance(p, image.Point{f.Rectangle.Max.X / 2, f.Rectangle.Max.Y / 2})
				cost := neighbourhoodCost * f.roflCost(p) * f.lolCost(p)
				candidates <- candidate{p, cost}
			}
		}()
	}

	// Wait for the workers to produce all candidates in a seperate goroutine
	// so that we can start working on the candidates. We will still wait for
	// the candidates chan to be closed before this function returns.
	go func() {
		wg.Wait()
		close(candidates)
	}()

	// Start with the first candidate to make sure we don't pick an already taken position.
	start := <-candidates
	best := start.point
	cheapest := start.cost

	// Find the best candidate from the frontier points...
	for candidate := range candidates {
		if candidate.cost < cheapest {
			cheapest = candidate.cost
			best = candidate.point
		}
	}

	f.take(best, c)
	return best
}

// Get all the possible neighbourns of the given point.
func (f *frontier) neighbours(p image.Point, moore int, filter func(image.Point) bool) []image.Point {
	// There can be at most 8 neighbours
	neighbours := make([]image.Point, 0, 8)

	maxX := f.Rectangle.Max.X
	maxY := f.Rectangle.Max.Y

	for dx := -moore; dx <= moore; dx++ {
		for dy := -moore; dy <= moore; dy++ {
			if dx == 0 && dy == 0 {
				continue
			}

			np := p
			np.X = (np.X + maxX + dx) % maxX
			np.Y = (np.Y + maxY + dy) % maxY

			if filter(np) {
				neighbours = append(neighbours, np)
			}
		}
	}

	return neighbours
}

// Get only the available (unpainted) neighbours of the given point.
func (f *frontier) availableNeighbours(p image.Point, moore int) []image.Point {
	return f.neighbours(p, moore, func(point image.Point) bool {
		_, ok := f.canvas[point]
		return !ok
	})
}

// Get only the taken (painted) neighbours of the given point.
func (f *frontier) takenNeighbours(p image.Point, moore int) []image.Point {
	return f.neighbours(p, moore, func(point image.Point) bool {
		_, ok := f.canvas[point]
		return ok
	})
}

// Extend the frontier with more points.
func (f *frontier) extend(points []image.Point) {
	for _, p := range points {
		f.frontier[p] = true
	}
}

// Take some points from the frontier.
func (f *frontier) take(p image.Point, c color.Color) {
	// We need to make sure we don't place a color ontop of another, if se we try again.
	if _, ok := f.canvas[p]; ok {
		panic("I tried to take an occupied position")
	}

	f.canvas[p] = c
	f.extend(f.availableNeighbours(p, 1))
	delete(f.frontier, p)
}

/*
func (f *frontier) distancePrevious(p image.Point) float64 {

*/

/*
func (f *frontier) distanceDirection(offset int) float64 {
	pp := allrgb.OffsetToPoint(f.previousOffsets[0], f.Rectangle)
	p := allrgb.OffsetToPoint(offset, f.Rectangle)
	direction := math.Atan2(float64(pp.X-p.X), float64(pp.Y-p.Y))
	diff := math.Abs(math.Pi/3 - (direction - f.direction))
	return math.Pow(diff, 2)
}
*/

func (f *frontier) sineCost(p image.Point) float64 {
	x := image.Point{
		X: int(float64(f.Rectangle.Max.X)/2.0*math.Sin(float64(len(f.canvas))/2.0)) % f.Rectangle.Max.X,
		Y: int(float64(f.Rectangle.Max.Y)/2.0*math.Cos(float64(len(f.canvas))/1.0)) % f.Rectangle.Max.Y,
	}

	return f.cityDistance(x, p)
}

func (f *frontier) lolCost(p image.Point) float64 {
	points := [...]image.Point{
		p.Add(image.Point{X: -1, Y: 0}),
		p.Add(image.Point{X: 0, Y: -1}),
	}

	var count int
	for _, np := range points {
		np.X %= f.Rectangle.Max.X
		np.Y %= f.Rectangle.Max.Y
		if _, ok := f.canvas[np]; ok {
			count++
		}
	}

	return 1.0 / float64(count)
}

func (f *frontier) roflCost(p image.Point) float64 {
	var column, row int

	for y := 0; y < f.Rectangle.Max.Y; y++ {
		if _, ok := f.canvas[p.Add(image.Point{X: 0, Y: y})]; ok {
			column++
		}
	}

	for x := 0; x < f.Rectangle.Max.X; x++ {
		if _, ok := f.canvas[p.Add(image.Point{X: x, Y: 0})]; ok {
			row++
		}
	}

	return float64(1+column) / float64(1+row)

	return float64(f.Rectangle.Max.X+f.Rectangle.Max.Y) / float64(1+column+row)
}

func (f *frontier) distanceCost(p image.Point) float64 {
	return f.distance(p, image.Point{f.Rectangle.Max.X / 2, f.Rectangle.Max.Y / 2})
}

func (f *frontier) distance(a, b image.Point) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return math.Sqrt(
		math.Pow(math.Min(dx, float64(f.Rectangle.Dx())-dx), 2) +
			math.Pow(math.Min(dy, float64(f.Rectangle.Dy())-dy), 2))
}

func (f *frontier) cityDistance(a, b image.Point) float64 {
	dx := math.Abs(float64(a.X - b.X))
	dy := math.Abs(float64(a.Y - b.Y))
	return math.Min(dx, float64(f.Rectangle.Dx())-dx) + math.Min(dy, float64(f.Rectangle.Dy())-dy)
}

func (f *frontier) neighbourhoodCost(p image.Point) float64 {
	neighbourhood := len(f.takenNeighbours(p, 2))

	switch {
	case neighbourhood <= 10:
		return float64(10-neighbourhood) / 10.0
	case neighbourhood > 20:
		return -99999.9
	default:
		return float64(neighbourhood-11) / 14.0
	}

	if neighbourhood == 24 {
		return -999999.0
	} else if neighbourhood > 12 {
		return float64(neighbourhood)
	} else {
		return -(12 - float64(neighbourhood))
	}
	return (float64(neighbourhood)) / 24.0
}

// Get a distance value for the differance of the given color to the slice of colors.
func (f *frontier) colorDistance(p image.Point, c color.Color) float64 {
	neighbours := f.takenNeighbours(p, 2)
	colors := make([]color.Color, len(neighbours))
	for i, np := range neighbours {
		colors[i] = f.canvas[np]
	}

	if len(colors) == 0 {
		return 0
	}

	diff := 0.0
	r, g, b, _ := c.RGBA()

	for _, cc := range colors {
		rr, gg, bb, _ := cc.RGBA()
		diff += math.Pow(float64(r)-float64(rr), 2)
		diff += math.Pow(float64(g)-float64(gg), 2)
		diff += math.Pow(float64(b)-float64(bb), 2)
	}

	return 441 / (math.Sqrt(diff) / float64(len(colors)))
}
