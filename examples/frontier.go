package examples

import (
	"fmt"
	"image"
	"image/color"
)

type frontier struct {
	image.Rectangle
	frontier map[image.Point]bool
	taken    map[image.Point]bool
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	f.frontier = make(map[image.Point]bool)
	f.taken = make(map[image.Point]bool)
	f.Max = rect.Max

	// Start at some place.
	f.take([]image.Point{image.Point{X: 0, Y: 0}})

	fmt.Printf("%#v\n", f)

	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	p := image.Point{}
	for p, _ = range f.frontier {
		break
	}

	f.take([]image.Point{p})

	return p
}

// Get all the possible neighbourns of the given point.
func (f frontier) neighbours(p image.Point) []image.Point {
	// There can be at most 8 neighbours
	neighbours := make([]image.Point, 0, 8)

	if p.Y > 0 {
		neighbours = append(neighbours, image.Point{X: p.X, Y: p.Y - 1}) // North
		if p.X < f.Max.X {
			neighbours = append(neighbours, image.Point{X: p.X + 1, Y: p.Y - 1}) // North-east
		}
		if p.X > 0 {
			neighbours = append(neighbours, image.Point{X: p.X - 1, Y: p.Y - 1}) // North-west
		}
	}

	if p.Y < f.Max.Y {
		neighbours = append(neighbours, image.Point{X: p.X, Y: p.Y + 1}) // South
		if p.X < f.Max.X {
			neighbours = append(neighbours, image.Point{X: p.X + 1, Y: p.Y + 1}) // South-east
		}
		if p.X > 0 {
			neighbours = append(neighbours, image.Point{X: p.X - 1, Y: p.Y + 1}) // South-west
		}
	}

	if p.X < f.Max.X {
		neighbours = append(neighbours, image.Point{X: p.X + 1, Y: p.Y}) // East
	}

	if p.X > 0 {
		neighbours = append(neighbours, image.Point{X: p.X - 1, Y: p.Y}) // West
	}

	return neighbours
}

// Get only the available (unpainted) neighbours of the given point.
func (f frontier) availableNeighbours(p image.Point) []image.Point {
	neighbours := f.neighbours(p)
	available := make([]image.Point, 0, len(neighbours))

	for _, neighbour := range neighbours {
		if taken, ok := f.taken[neighbour]; !ok || ok && !taken {
			available = append(available, neighbour)
		}
	}

	return available
}

// Extend the frontier with more points.
func (f *frontier) extend(ps []image.Point) {
	for _, p := range ps {
		f.frontier[p] = true
	}
}

// Take some points from the frontier.
func (f *frontier) take(ps []image.Point) {
	for _, p := range ps {
		f.taken[p] = true
		delete(f.frontier, p)
		f.extend(f.availableNeighbours(p))
	}
}
