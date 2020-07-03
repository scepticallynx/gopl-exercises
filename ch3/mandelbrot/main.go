/*
Mandelbrot emits a PNG image of the Mandelbrot fractal.

	Exercise 3.5: Implement full-color Mandelbrot set using the function image.NewRGBA
	and the type color.RGBA or color.YCbCr

	Exercise 3.6: Supersampling. pixel / 4
	A bit of googling + tuning from
	https://stackoverflow.com/questions/45094282/how-to-reduce-the-effect-of-pixelation-by-computing-the-color-value-at-several-p

	Exercise 3.7: Newton's method (z⁴ - 1 = 0). Shade each starting point by the
	number of iterations required to get close to one of the four roots. Color
	each point by the root it approaches.

	Colorize: https://www.paridebroggi.com/blogpost/2015/05/06/fractal-continuous-coloring/
*/
package main

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/cmplx"
	"os"
)

var (
	xmin, ymin, xmax, ymax = -2.1, -2.1, +2.1, +2.1 // left and right boundaries for plotting
	width, height          = 1000, 1000             // number of grid points
	iterations             = 50                     // number of iterations
	precision              = 1e-10                  // desired precision goal
)

const (
	_w, _h   = 1000, 1000
	wx2, hx2 = _w * 2, _h * 2 // number of grid point for antialiasing
)

type fractalF func(z complex128) color.RGBA

func main() {
	img := image.NewRGBA(image.Rect(0, 0, width, height))

	noSampling(img)

	f, err := os.Create("img.png")
	if err != nil {
		log.Fatalf("%v", err)
	}
	err = png.Encode(f, img) // NOTE: ignoring errors
	if err != nil {
		log.Fatal(err)
	}
}

func supersampling(img *image.RGBA) {
	var x2colors [wx2][hx2]color.Color

	for py := 0; py < hx2; py++ {
		y := float64(py)/hx2*(ymax-ymin) + ymin
		for px := 0; px < wx2; px++ {
			x := float64(px)/wx2*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			x2colors[px][py] = mandelbrot(z)
		}
	}

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			ix2, jx2 := i*2, j*2

			r1, b1, g1, a1 := x2colors[ix2][jx2].RGBA()
			r2, b2, g2, a2 := x2colors[ix2+1][jx2].RGBA()
			r3, b3, g3, a3 := x2colors[ix2+1][jx2+1].RGBA()
			r4, b4, g4, a4 := x2colors[ix2][jx2+1].RGBA()

			// ((x1+x2+x3+x4)/4)*(255/65535) -> (x1+x2+x3+x4)/1028
			avg := color.RGBA{
				R: uint8((r1 + r2 + r3 + r4) / 1028),
				G: uint8((g1 + g2 + g3 + g4) / 1028),
				B: uint8((b1 + b2 + b3 + b4) / 1028),
				A: uint8((a1 + a2 + a3 + a4) / 1028),
			}

			img.Set(i, j, avg)
		}
	}
}

func noSampling(img *image.RGBA) {
	// xmin = xmin+math.Floor(width/4)/zoom
	// ymin = -math.Floor(height/4)/zoom+iterations/zoom+ymin

	for py := 0; py < height; py++ {
		y := float64(py)/float64(height)*(ymax-ymin) + ymin
		for px := 0; px < width; px++ {
			x := float64(px)/float64(width)*(xmax-xmin) + xmin
			z := complex(x, y)
			// Image point (px, py) represents complex value z.
			img.Set(px, py, newton(z))
		}
	}
}

func mandelbrot(z complex128) color.Color {
	var v complex128

	for n := 0; n < iterations; n++ {
		v = v*v + z
		if cmplx.Abs(v) > 2 {
			i := math.Log(float64(n)) / math.Log(float64(iterations))
			c := uint8(i * 255)
			return color.RGBA{R: c - 20, G: c, B: c - 40, A: uint8(255 * i)}
		}
	}

	return color.Black
}

// newton's

var f = func(z complex128) complex128 { return z*z*z*z - 1 }
var derivative = func(z complex128) complex128 { return 4 * z * z * z }
var funcZ = func(z complex128) complex128 { return f(z) / derivative(z) }

func newton(z complex128) color.Color {
	// f := func(z complex128) complex128 { return z*z*z*z - 1	}
	// derivative := func(z complex128) complex128 { return 4 * z*z*z }
	// funcZ := func(z complex128) complex128 {return f(z)/derivative(z)}

	for n := 0; n < iterations; n++ {
		z -= funcZ(z)
		if cmplx.Abs(f(z)) < precision {
			// (frequency * continuous index + phase) * 127.5 + 127.5
			// r, i := real(z), imag(z)
			return color.RGBA{
				R: uint8(math.Sin(0.016*float64(n+4))*230 + 25),
				G: uint8(math.Sin(0.027*float64(n+2))*230 + 25),
				B: uint8(math.Sin(0.01*float64(n+1))*230 + 25),
				A: 255,
			}
		}
	}

	return color.Black
}

// Some other interesting functions:

func acos(z complex128) color.Color {
	v := cmplx.Acos(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{192, blue, red}
}

func sqrt(z complex128) color.Color {
	v := cmplx.Sqrt(z)
	blue := uint8(real(v)*128) + 127
	red := uint8(imag(v)*128) + 127
	return color.YCbCr{128, blue, red}
}
