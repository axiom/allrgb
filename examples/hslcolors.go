package examples

import (
	sadcolor "code.google.com/p/sadbox/color"
	"github.com/axiom/allrgb"
	"image/color"
	"sort"
)

type component int

const (
	H component = iota
	S
	L
)

func HSLColorProducer() chan color.Color {
	return NewHSLColorProducer(H, S, L, false, true, true).Produce()
}

func NewHSLColorProducer(c1, c2, c3 component, r1, r2, r3 bool) allrgb.ColorProducer {
	nextColor := make(chan color.Color)

	go func() {
		colors := make([]sadcolor.HSL, 32*32*32)
		i := 0
		for color := range RGBColorProducer() {
			colors[i] = sadcolor.HSLModel.Convert(color).(sadcolor.HSL)
			i++
		}

		orderColors := hslcolors{
			colors:  colors,
			order:   [3]component{c1, c2, c3},
			reverse: [3]bool{r1, r2, r3},
		}
		sort.Sort(orderColors)

		for _, color := range orderColors.colors {
			nextColor <- color
		}
		close(nextColor)
	}()

	return allrgb.ColorProducerFunc(func() chan color.Color {
		return nextColor
	})
}

type hslcolors struct {
	colors  []sadcolor.HSL
	order   [3]component
	reverse [3]bool
}

func (cs hslcolors) Len() int      { return len(cs.colors) }
func (cs hslcolors) Swap(i, j int) { cs.colors[i], cs.colors[j] = cs.colors[j], cs.colors[i] }
func (cs hslcolors) Less(i, j int) bool {
	ci := []float64{cs.colors[i].H, cs.colors[i].S, cs.colors[i].L}
	cj := []float64{cs.colors[j].H, cs.colors[j].S, cs.colors[j].L}

	if ci[cs.order[0]] < cj[cs.order[0]] {
		return true != cs.reverse[0]
	}

	if ci[cs.order[0]] > cj[cs.order[0]] {
		return false != cs.reverse[0]
	}

	if ci[cs.order[1]] < cj[cs.order[1]] {
		return true != cs.reverse[1]
	}

	if ci[cs.order[1]] > cj[cs.order[1]] {
		return false != cs.reverse[1]
	}

	return ci[cs.order[2]] < cj[cs.order[2]] != cs.reverse[2]
}
