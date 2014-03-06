package examples

import (
	"fmt"
	"image"
	"image/color"
)

type frontier struct {
	image.Rectangle
	frontier map[image.Point]bool
	taken    map[image.Point]color.Color
}

func NewFrontier(rect image.Rectangle) *frontier {
	f := frontier{}
	f.frontier = make(map[image.Point]bool)
	f.taken = make(map[image.Point]color.Color)
	f.Max = rect.Max

	// Start at some place.
	f.extend([]image.Point{{rect.Dx() / 2, rect.Dy() / 2}})

	return &f
}

func (f *frontier) Place(c color.Color) image.Point {
	var best image.Point
	shortest := 0xfffffffffffffff

	// We want to find the place with the least color differance in the frontier.
	for p, _ := range f.frontier {
		var colors []color.Color
		for _, neighbour := range f.takenNeighbours(p) {
			colors = append(colors, f.taken[neighbour])
		}

		if distance := colorDistance(c, colors); distance < shortest {
			shortest = distance
			best = p
		}
	}

	if shortest == 0xfffffffffffffff {
		fmt.Println("COULD NOT FIND GOOD MATCH")
	}

	if _, ok := f.taken[best]; ok {
		panic("OMG TRIED TO PLACE ON TAKEN POINT")
	}

	f.take(best, c)
	return best
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
	f.taken[p] = c
	f.extend(f.availableNeighbours(p))
	delete(f.frontier, p)
}

// Get a distance value for the differance of the given color to the slice of colors.
func colorDistance(c color.Color, colors []color.Color) int {
	var diff, r, g, b int
	var rr, gg, bb uint32

	rr, gg, bb, _ = c.RGBA()
	r = int(rr)
	g = int(gg)
	b = int(bb)

	for _, color := range colors {
		rr, gg, bb, _ := color.RGBA()

		diff += abs(r-int(rr)) / 3
		diff += abs(g-int(gg)) / 3
		diff += abs(b-int(bb)) / 3
	}

	return diff / (1 + len(colors))
}

func abs(v int) int {
	if v < 0 {
		return -int(v)
	} else {
		return int(v)
	}
}
