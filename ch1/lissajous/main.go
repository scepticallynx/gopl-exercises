// Lissajous generates GIF animations of random Lissajous figures.
// Exercise 1.5: Change the Lissajous program's color palette to green on black.
// Exercise 1.6: Modify the program to produce images in multiple colors by adding them to the palette.
// Exercise 1.12: Modify the Lissajous server to read parameters from URL...
package main

import (
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Packages not needed by version in book.

// Exercise 1.5
// Exercise 1.6
var palette = []color.Color{
	color.Black,
	color.RGBA{R:0, G: 255, B:0, A:1},
	color.RGBA{R:255, G: 0, B:0, A:1},
	color.RGBA{R:0, G: 0, B:255, A:1},
	color.RGBA{R:40, G: 80, B:120, A:10},
	}

func main() {
	// The sequence of images is deterministic unless we seed
	// the pseudo-random number generator using the current time.
	// Thanks to Randall McPherson for pointing out the omission.
	rand.Seed(time.Now().UTC().UnixNano())

	if len(os.Args) > 1 && os.Args[1] == "web" {

		handler := func(w http.ResponseWriter, r *http.Request) {
			// Exercise 1.12
			if err := r.ParseForm(); err != nil {
				log.Print(err)
			}

			cycles := 5
			size := 100
			delay := 8

			for k, v := range r.Form {
				fmt.Printf("Form[%q] = %q\n", k, v)
				switch k {
				case "cycles":
					cycles, _ = strconv.Atoi(v[0])
				case "size":
					size, _ = strconv.Atoi(v[0])
				case "delay":
					delay, _ = strconv.Atoi(v[0])
				}
			}

			lissajous(w, cycles, size, delay)
		}
		http.HandleFunc("/", handler)

		log.Fatal(http.ListenAndServe("localhost:8000", nil))
		return
	}

	lissajous(os.Stdout, 5, 100, 8)
}

func lissajous(out io.Writer, cyclesN, canvasSize, delayMS int) {
	const (
		// cycles  = 5     // number of complete x oscillator revolutions
		res     = 0.001 // angular resolution
		// size    = 100   // image canvas covers [-size..+size]
		nframes = 64    // number of animation frames
		// delay   = 8     // delay between frames in 10ms units
	)
	cycles := float64(cyclesN)
	size := canvasSize
	delay := delayMS

	freq := rand.Float64() * 3.0 // relative frequency of y oscillator
	anim := gif.GIF{LoopCount: nframes}
	phase := 0.0 // phase difference
	for i := 0; i < nframes; i++ {
		rect := image.Rect(0, 0, 2*size+1, 2*size+1)
		img := image.NewPaletted(rect, palette)
		for t := 0.0; t < cycles*2*math.Pi; t += res {
			x := math.Sin(t)
			y := math.Sin(t*freq + phase)
			img.SetColorIndex(size+int(x*float64(size)+0.5), size+int(y*float64(size)+0.5), 4)
		}
		phase += 0.1
		anim.Delay = append(anim.Delay, delay)
		anim.Image = append(anim.Image, img)
	}
	err := gif.EncodeAll(out, &anim) // NOTE: ignoring encoding errors
	if err != nil {
		log.Fatal(err)
	}
}

