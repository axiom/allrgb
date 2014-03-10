package main

import (
	"fmt"
	"github.com/axiom/allrgb"
	. "github.com/axiom/allrgb/examples"
	"image"
	"runtime"
	"time"
	//"net/http"
	//_ "net/http/pprof"
)

func main() {
	runtime.GOMAXPROCS(8)

	/*
		go func() {
			fmt.Println(http.ListenAndServe("localhost:8080", nil))
		}()
	*/

	rect := image.Rectangle{Max: image.Point{X: 256, Y: 128}}

	configurations := map[string]image.Image{
		"trivial.png": allrgb.ColorDetermined(
			rect,
			NewHSLColorProducer(H, L, S, true, false, true),
			NewTrivialPlacer(rect)),

		"hilbert.png": allrgb.ColorDetermined(
			rect,
			NewHSLColorProducer(S, H, L, true, false, true),
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
		NewHSLColorProducer(H, L, S, false, false, true),
		//allrgb.ColorProducerFunc(RGBColorProducer),
		NewFrontier(rect),
		100,
		fmt.Sprintf("frames/frontier-%v", time.Now().Format("20060102T150405")),
	)

}
