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

func (s *MySuite) TestHSLColorProducer(c *C) {
	c.Check(countUniqueColors(allrgb.ColorProducerFunc(HSLColorProducer)), Equals, 32768)
}

func (s *MySuite) TestRGBColorProducer(c *C) {
	c.Check(countUniqueColors(allrgb.ColorProducerFunc(RGBColorProducer)), Equals, 32768)
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
				{0, 1},
				{0, 3},
				{1, 0},
				{1, 1},
				{1, 3},
				{7, 0},
				{7, 1},
				{7, 3},
			},
		},
		{
			image.Point{1, 0},
			[]image.Point{
				{0, 0},
				{0, 1},
				{0, 3},
				{1, 1},
				{1, 3},
				{2, 0},
				{2, 1},
				{2, 3},
			},
		},
		{
			image.Point{1, 1},
			[]image.Point{
				{0, 0},
				{0, 1},
				{0, 2},
				{1, 0},
				{1, 2},
				{2, 0},
				{2, 1},
				{2, 2},
			},
		},
		{
			image.Point{7, 3},
			[]image.Point{
				{0, 0},
				{0, 2},
				{0, 3},
				{6, 0},
				{6, 2},
				{6, 3},
				{7, 0},
				{7, 2},
			},
		},
	}

	for _, tc := range testcases {
		found := frontier.neighbours(tc.point)
		expected := tc.neighbours
		sort.Sort(neighbours(found))
		sort.Sort(neighbours(expected))

		c.Check(found, DeepEquals, expected)
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

	if ns[i].X > ns[j].X {
		return false
	}

	return ns[i].Y < ns[j].Y
}
