package main

import (
	"fmt"
	"github.com/axiom/allrgb"
	. "github.com/axiom/allrgb/examples"
	"image"
	"time"
	//"net/http"
	//_ "net/http/pprof"
)

func main() {

	/*
		go func() {
			fmt.Println(http.ListenAndServe("localhost:8080", nil))
		}()
	*/

	rect := image.Rectangle{Max: image.Point{X: 256, Y: 128}}

	configurations := map[string]image.Image{
		"trivial.png": allrgb.ColorDetermined(
			rect,
			allrgb.ColorProducerFunc(HSLColorProducer),
			NewTrivialPlacer(rect)),

		"hilbert.png": allrgb.ColorDetermined(
			rect,
			allrgb.ColorProducerFunc(HSLColorProducer),
			NewHilbertPlacer(rect)),

		/*
			"frontier": allrgb.ColorDetermined(
				rect,
				allrgb.ColorProducerFunc(RGBColorProducer),
				NewFrontier(rect)),
		*/
	}

	for name, img := range configurations {
		if err := allrgb.SaveImage(name, img); err != nil {
			fmt.Printf("Could not do %v: %v\n", name, err)
		}
	}

	allrgb.ColorDeterminedFrameSaver(
		rect,
		NewHSLColorProducer(H, S, L, true, true, true),
		NewFrontier(rect),
		100,
		fmt.Sprintf("frames/frontier-%v", time.Now().Format("20060102T150405")),
	)

}
