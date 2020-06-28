// Exercise 2.2: Write a general purpose unit-conversion program analogous to cf (gopl.io/ch2/cf/main.go)
// that read numbers from its command-line arguments or from the standard input
// if there are no arguments, and converts each number into units like temperature
// in Celsius and Fahrenheit, length in feet and meters, weight in pounds and
// kilograms, and the like.
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	"gopl.io/ch2/tempconv"
)

// length
const lengthC float64 = 3.281

type Foot float64
type Meter float64

func (f Foot) String() string  { return fmt.Sprintf("%g ft", f) }
func (m Meter) String() string { return fmt.Sprintf("%g m", m) }

// MToF converts Meter to Foot
func MToF(m Meter) Foot { return Foot(float64(m) * lengthC) }

// FToM converts Foot to Meter
func FToM(f Foot) Meter { return Meter(float64(f) / lengthC) }

// weight
const weightC float64 = 2.205

type Kilogram float64
type Pound float64

func (f Pound) String() string    { return fmt.Sprintf("%g lb", f) }
func (k Kilogram) String() string { return fmt.Sprintf("%g kg", k) }

// KToP converts Kilogram to Pound
func KToP(k Kilogram) Pound { return Pound(float64(k) * weightC) }

// PToK converts Pound to Kilogram
func PToK(p Pound) Kilogram { return Kilogram(float64(p) / weightC) }

func convert(v float64) {
	// temperature
	celsius := tempconv.Celsius(v)
	fahrenheit := tempconv.Fahrenheit(v)
	fmt.Printf("%s = %s, %s = %s\n",
		celsius, tempconv.CToF(celsius), fahrenheit, tempconv.FToC(fahrenheit))

	// weight
	k := Kilogram(v)
	p := Pound(v)
	fmt.Printf("%s = %s, %s = %s\n", k, KToP(k), p, PToK(p))

	// length
	m := Meter(v)
	f := Foot(v)
	fmt.Printf("%s = %s, %s = %s\n", m, MToF(m), f, FToM(f))
}

func main() {
	var v float64

	if len(os.Args[1:]) > 0 {
		for n, arg := range os.Args[1:] {
			var err error
			v, err = strconv.ParseFloat(arg, 64)
			if err != nil {
				log.Printf("Failed to parse arg #%d %q to float64: %v", n, arg, err)
				continue
			}
			convert(v)
		}
	} else {
		input := bufio.NewScanner(os.Stdin)
		for input.Scan() {
			var err error
			v, err = strconv.ParseFloat(input.Text(), 64)
			if err != nil {
				log.Printf("Failed to parse value to float from stdin: %v", err)
			}
			convert(v)
		}
	}
}
