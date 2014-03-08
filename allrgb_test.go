package allrgb

import (
	"image"
	"testing"
)

func TestOffsetConversion(t *testing.T) {
	rect := image.Rectangle{
		Max: image.Point{X: 8, Y: 4},
	}

	testcases := []struct {
		point  image.Point
		offset int
	}{
		// The four corners
		{image.Point{0, 0}, 0},
		{image.Point{7, 0}, 7},
		{image.Point{7, 3}, 31},
		{image.Point{0, 3}, 24},

		{image.Point{1, 0}, 1},
		{image.Point{0, 1}, 8},
		{image.Point{2, 2}, 18},
		{image.Point{7, 3}, 31},
	}

	for _, testcase := range testcases {
		if offset := PointToOffset(testcase.point, rect); testcase.offset != offset {
			t.Errorf("PointToOffset(%#v) = %v, want %v", testcase.point, offset, testcase.offset)
		}

		if point := OffsetToPoint(PointToOffset(testcase.point, rect), rect); point != testcase.point {
			t.Errorf(
				"OffsetToPoint(PointToOffset(%#v)) = %#v, want %#v",
				testcase.point,
				point,
				testcase.point,
			)
		}

		if offset := PointToOffset(OffsetToPoint(testcase.offset, rect), rect); offset != testcase.offset {
			t.Errorf(
				"PointToOffset(OffsetToPoint(%#v)) = %#v, want %#v",
				testcase.offset,
				offset,
				testcase.offset,
			)
		}
	}
}
