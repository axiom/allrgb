package examples

import (
	sadcolor "code.google.com/p/sadbox/color"
	"fmt"
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
	placeQueue    chan queueItem
	previousPoint image.Point
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	f.frontier = make(map[image.Point]bool)
	f.taken = make(map[image.Point]sadcolor.HSL)
	f.placeQueue = make(chan queueItem)
	f.Max = rect.Max

	// Start at some place.
	start := image.Point{rect.Dx() / 2, rect.Dy() / 2}
	f.previousPoint = start
	f.extend([]image.Point{start})

	for i := 0; i < 1; i++ {
		go func() {
			f.place()
		}()
	}

	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	pointResult := make(chan image.Point)
	f.placeQueue <- queueItem{result: pointResult, color: c}
	return <-pointResult
}

// Place colors from the queue.
func (f *frontier) place() {
	fmt.Println("Started worker")
	for queueItem := range f.placeQueue {
		best := f.previousPoint
		shortest := 10e100

		// We want to find the place with the least color differance in the frontier.
		for p, _ := range f.frontier {
			var colors []sadcolor.HSL
			for _, neighbour := range f.takenNeighbours(p) {
				colors = append(colors, f.taken[neighbour])
			}

			if distance := colorDistance(queueItem.color, colors) + f.distancePrevious(p); distance < shortest {
				shortest = distance
				best = p
			}
		}

		f.take(best, queueItem.color)

		queueItem.result <- best
		close(queueItem.result)
	}
}

// Get all the possible neighbourns of the given point.
func (f frontier) neighbours(p image.Point) []image.Point {
	// There can be at most 8 neighbours
	neighbours := make([]image.Point, 0, 8)

	for dx := -1; dx <= 1; dx++ {
		if p.X+dx >= 0 && p.X+dx < f.Dx() {
			for dy := -1; dy <= 1; dy++ {
				if p.Y+dy >= 0 && p.Y+dy < f.Dy() {
					neighbours = append(neighbours, image.Point{X: p.X + dx, Y: p.Y + dy})
				}
			}
		}
	}

	return neighbours
}

// Get only the available (unpainted) neighbours of the given point.
func (f frontier) availableNeighbours(p image.Point) []image.Point {
	neighbours := f.neighbours(p)
	available := make([]image.Point, 0, len(neighbours))

	for _, neighbour := range neighbours {
		if _, ok := f.taken[neighbour]; !ok {
			available = append(available, neighbour)
		}
	}

	return available
}

// Get only the taken (painted) neighbours of the given point.
func (f frontier) takenNeighbours(p image.Point) []image.Point {
	neighbours := f.neighbours(p)
	taken := make([]image.Point, 0, len(neighbours))

	for _, neighbour := range neighbours {
		if _, ok := f.taken[neighbour]; ok {
			taken = append(taken, neighbour)
		}
	}

	return taken
}

// Extend the frontier with more points.
func (f *frontier) extend(ps []image.Point) {
	for _, p := range ps {
		f.frontier[p] = true
	}
}

// Take some points from the frontier.
func (f *frontier) take(p image.Point, c color.Color) {
	// We need to make sure we don't place a color ontop of another, if se we try again.
	if _, ok := f.taken[p]; ok {
		panic("I lost the race")
	}

	f.taken[p] = sadcolor.HSLModel.Convert(c).(sadcolor.HSL)
	f.extend(f.availableNeighbours(p))
	delete(f.frontier, p)
	f.previousPoint = p
}

func (f frontier) distancePrevious(p image.Point) float64 {
	dx := f.previousPoint.X - p.X
	dy := f.previousPoint.Y - p.Y
	return 3 * math.Sqrt(math.Pow(float64(dx), 2)+math.Pow(float64(dy), 2))
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
