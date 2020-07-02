/*
Comma prints its argument numbers with a comma at each power of 1000.

 Example:
 	$ go build gopl.io/ch3/comma
	$ ./comma 1 12 123 1234 1234567890
 	1
 	12
 	123
 	1,234
 	1,234,567,890

	Exercise 3.10: Write a non-recursive version of comma, using bytes.Buffer
	instead of string concatenation.

	Exercise 3.11: Enhance comma so that is deals correctly with floating-point
	numbers and an optional sign.
*/
package main

import (
	"bytes"
	"fmt"
	"os"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", commaBuffer(os.Args[i]))
	}
}

const splitAfter int = 3

// comma inserts commas in a non-negative decimal integer string.
func comma(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}
	return comma(s[:n-3]) + "," + s[n-3:]
}

// commaBuffer is a non-recursive version of comma (Exercise 3.10)
func commaBuffer(s string) string {
	var b bytes.Buffer

	n := len(s)
	if n <= splitAfter {
		return s
	}

	// cycle through string writing comma after every 3d letters
	for i := 0; i < n; i++ {
		if i >= splitAfter && i%splitAfter == 0 {
			b.WriteString(", ")
		}

		b.WriteByte(s[i])
	}

	return b.String()
}

// Exercise 3.11
// commaEnhanced also parses floating point numbers from specified string with
// optional sign.
