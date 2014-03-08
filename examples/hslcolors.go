package examples

import (
	sadcolor "code.google.com/p/sadbox/color"
	"image/color"
)

func HSLColorProducer() chan color.Color {
	nextColor := make(chan color.Color)
	go func() {
		for h := 0; h < 32; h++ {
			for l := 0; l < 32; l++ {
				for s := 0; s < 32; s++ {
					nextColor <- sadcolor.HSL{
						H: float64(h) / 31.0,
						S: float64(31-s) / 31.0,
						L: float64(31-l) / 31.0,
					}
				}
			}
		}
		close(nextColor)
	}()
	return nextColor
}
