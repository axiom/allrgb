package examples

import (
	"github.com/axiom/allrgb"
	"image"
	"image/color"
	. "launchpad.net/gocheck"
	"sort"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})

func TestHSLColorProducer(t *testing.T) {
	const expected = 32768
	if unique := countUniqueColors(allrgb.ColorProducerFunc(HSLColorProducer)); unique != expected {
		t.Errorf("HSLColorProducer produced %v unique colors, want %v", unique, expected)
	}
}
func TestRGBColorProducer(t *testing.T) {
	const expected = 32768
	if unique := countUniqueColors(allrgb.ColorProducerFunc(RGBColorProducer)); unique != expected {
		t.Errorf("RGBColorProducer produced %v unique colors, want %v", unique, expected)
	}
}

func (s *MySuite) TestFrontierNeighbours(c *C) {
	rect := image.Rectangle{
		Max: image.Point{8, 4},
	}
	frontier := NewFrontier(rect)
	testcases := []struct {
		point      image.Point
		neighbours []image.Point
	}{
		{
			image.Point{0, 0},
			[]image.Point{
				{1, 0},
				{1, 1},
				{0, 1},
			},
		},
		{
			image.Point{1, 0},
			[]image.Point{
				{0, 0},
				{0, 1},
				{1, 1},
				{2, 1},
				{2, 0},
			},
		},
		{
			image.Point{1, 1},
			[]image.Point{
				{0, 0},
				{1, 0},
				{2, 0},
				{2, 1},
				{2, 2},
				{1, 2},
				{0, 2},
				{0, 1},
			},
		},
		{
			image.Point{7, 3},
			[]image.Point{
				{6, 2},
				{7, 2},
				{6, 3},
			},
		},
	}

	for _, tc := range testcases {
		found := frontier.neighbours(tc.point)
		sort.Sort(neighbours(found))
		sort.Sort(neighbours(tc.neighbours))

		c.Check(found, DeepEquals, tc.neighbours)
	}
}

func countUniqueColors(producer allrgb.ColorProducer) int {
	unique := make(map[color.Color]bool)
	for c := range producer.Produce() {
		unique[c] = true
	}

	return len(unique)
}

type neighbours []image.Point

func (ns neighbours) Len() int      { return len(ns) }
func (ns neighbours) Swap(i, j int) { ns[i], ns[j] = ns[j], ns[i] }
func (ns neighbours) Less(i, j int) bool {
	if ns[i].X < ns[j].X {
		return true
	}

	if ns[i].X == ns[j].X {
		return ns[i].Y < ns[j].Y
	}

	return false
}
