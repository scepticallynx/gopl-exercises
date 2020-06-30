// Surface computes an SVG rendering of a 3-D surface function.
// Try:
// http://localhost:8080/?surface=f&cells=100&xyrange=20&width=600&height=300
// http://localhost:8080/?surface=saddle&cells=50&xyrange=20&width=600&height=400
// http://localhost:8080/?surface=eggbox&cells=40&xyrange=10&width=550&height=200
package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
)

// defaults
var (
	width, height                                 float64 = 500, 500 // canvas size in pixels
	cells                                                 = 200      // number of grid cells
	xyrange, xyscale, zscale, angle, sin30, cos30 float64            // default placeholders
)

func calculateVariables() {
	if xyrange == 0 {
		xyrange = 15.0 // axis ranges (-xyrange..+xyrange)
	}
	xyscale = width / 2 / xyrange                   // pixels per x or y unit
	zscale = height * 0.5                           // pixels per z unit
	angle = math.Pi / 6                             // angle of x, y axes (=30°)
	sin30, cos30 = math.Sin(angle), math.Cos(angle) // sin(30°), cos(30°)
}

// variants of the built surface
type theFunc func(x, y float64) float64

var functions = map[string]theFunc{"f": f, "saddle": saddle, "eggbox": eggbox}

// default surface
var chosenFunc = "f"

// Red*65536 + Green*256 + Blue
const colorMax = 255

var r, g, b uint8

func colorFormula(r, g, b uint8) int {
	return int(r)*65536 + int(g)*256 + int(b)
}

func main() {
	// Exercise 3.3: server like ch1/lissajous with parameters in query
	if len(os.Args) > 1 && os.Args[1] == "web" {
		handler := func(w http.ResponseWriter, r *http.Request) {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}

			for k, v := range r.Form {
				var err error
				switch k {
				case "xyrange":
					xyrange, err = strconv.ParseFloat(v[0], 64)
					if err != nil {
						log.Printf("Failed to parse float for xyrange: %v", err)
					}
				case "height":
					height, err = strconv.ParseFloat(v[0], 64)
					if err != nil {
						log.Printf("Failed to parse float for height: %v", err)
					}
				case "width":
					width, err = strconv.ParseFloat(v[0], 64)
					if err != nil {
						log.Printf("Failed to parse float for width: %v", err)
					}
				case "cells":
					cells, err = strconv.Atoi(v[0])
					if err != nil {
						log.Printf("Failed to parse float for cells: %v", err)
					}
				case "surface":
					v := v[0]
					_, exists := functions[v]
					if !exists {
						log.Printf("Unknown surface %q. Use one of 'f', 'eggbox', 'saddle'", v)
					} else {
						chosenFunc = v
					}
				default:
					log.Printf("Unsupported parameter passes in query: %s=%v", k, v)
				}
			}

			w.Header().Set("Content-Type", "image/svg+xml")
			calculateVariables()
			computeSurface(w)
		}

		http.HandleFunc("/", handler)

		log.Fatal(http.ListenAndServe("localhost:8080", nil))
		return
	} else { // write to file
		out, err := os.Create("img.svg")
		if err != nil {
			log.Fatal(err)
		}

		calculateVariables()
		computeSurface(out)
	}
}

func computeSurface(out io.Writer) {
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("<svg xmlns='http://www.w3.org/2000/svg' "+
		"style='stroke: grey; fill: white; stroke-width: 0.7' "+
		"width='%d' height='%d'>\n", int(width), int(height)))

	for i := 0; i < cells; i++ {
		for j := 0; j < cells; j++ {
			ax, ay := corner(i+1, j)
			bx, by := corner(i, j)
			cx, cy := corner(i, j+1)
			dx, dy := corner(i+1, j+1)

			buf.WriteString(fmt.Sprintf("<polygon points='%g,%g %g,%g %g,%g %g,%g' fill='#%06x' stroke-width='0.2' stroke='#DCDCDC'/>\n",
				ax, ay, bx, by, cx, cy, dx, dy, colorFormula(r, g, b)))
		}
	}

	buf.WriteString(fmt.Sprintln("</svg>"))

	// dump svg from buffer
	_, err := out.Write(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}
}

func corner(i, j int) (float64, float64) {
	// Find point (x,y) at corner of cell (i,j).
	x := xyrange * (float64(i)/float64(cells) - 0.5)
	y := xyrange * (float64(j)/float64(cells) - 0.5)

	// Compute surface height z.
	// One of predefined function from functions array
	z := functions[chosenFunc](x, y)

	// Exercise 3.1: ignore non-finite return value from f to avoid building
	// invalid polygons in svg
	if math.IsInf(z, 1) || math.IsNaN(z) { // +Inf | NaN
		z = 1
	} else if math.IsInf(z, -1) { // -Inf
		z = -1
	}

	/* Exercise 3.2 explanation:
	z is in range -1 <= z <= 1.
	Positive range of z (0..1, 0 inclusive) must be painted red.
	Negative (-1..0, 0 exclusive) - blue.
	To make a smooth transition from +1 to -1 (red-white-blue)
	*/

	v := colorMax * math.Abs(z)
	switch {
	case z > 0:
		r = colorMax
		b = uint8(colorMax - v)
		g = uint8(colorMax - v)
	case z < 0:
		b = colorMax
		r = uint8(colorMax - v)
		g = uint8(colorMax - v)
	default:
		r = colorMax
		b = colorMax
		g = colorMax
	}

	// Project (x,y,z) isometrically onto 2-D SVG canvas (sx,sy).
	sx := width/2 + (x-y)*cos30*xyscale
	sy := height/2 + (x+y)*sin30*xyscale - z*zscale

	return sx, sy
}

// Exercise 3.2: Experiment with visualizations of other functions from the math
// package. Produce eggbox, saddle, moguls.
func eggbox(x, y float64) float64 {
	return math.Cos(x) * math.Cos(y) / (angle * math.Pi)
}

func saddle(x, y float64) float64 {
	return math.Pow(y, 2)/(height/2) - math.Pow(x, 2)/(width/2)
}

// Exercise 3.1: If the function f returns a non-finite float64 value, the SVG
// file contains invalid <polygon> elements. ... Modify the program to skip invalid
// polygons.
// Return 0 when r == 0 to avoid 0 / 0 (f returns NaN)
func f(x, y float64) float64 {
	r := math.Hypot(x, y) // distance from (0,0)
	return math.Sin(r) / r
}
