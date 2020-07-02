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
	"unicode"
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Printf("  %s\n", commaEnhanced(os.Args[i]))
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
// commaEnhanced also parses floating point numbers with 1 decimal digit
// from specified string with optional sign (+/-).
func commaEnhanced(s string) string {
	var b bytes.Buffer

	var digitsCounter int
	var skipRunes = make(map[int]struct{})

	for n, r := range s {
		// skip current rune if it is found in skipRunes map
		if _, exists := skipRunes[n]; exists {
			continue
		}

		// catch non digits with custom handling each case
		// and incrementing digits counter only if rune is digit
		switch {
		case unicode.IsSymbol(r) || r == '-': // minus is not '-', it's '−' \U+2212
			break
		case unicode.IsPunct(r): // float separators: '.' and ','
			b.WriteRune(r)

			// look ahead only if this is not the end of the string
			if next := n + 1; next < len(s)-1 {
				b.WriteByte(s[next])
				skipRunes[next] = struct{}{}
			}

			continue
		default:
			digitsCounter++ // increment counter (need 3 digits in a row to put ',' after them)
		}

		b.WriteRune(r)

		if digitsCounter == 3 {
			if n+1 >= len(s) {
				break
			}

			next := s[n+1]
			if unicode.IsPunct(rune(next)) && rune(next) != '-' { // not a minus
				b.WriteByte(next)
				b.WriteByte(s[n+2])
				skipRunes[n+1] = struct{}{}
				skipRunes[n+2] = struct{}{}
			}

			b.WriteString(", ")
			digitsCounter = 0
		}
	}

	return b.String()
}
