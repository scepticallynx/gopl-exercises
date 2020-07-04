// Exercise 4.6: Write an in-place function that squashes each run of adjacent
// Unicode spaces (unicode.IsSpace) in a UTF-8-encoded []byte slice into a single
// ASCII space.
package main

import (
	"fmt"
	"unicode"
)

func squashSpaces(s []byte) []byte {
	var w rune // last written
	n := 0
	for _, b := range s {
		r := rune(b)
		if unicode.IsSpace(r) && unicode.IsSpace(w) {
			continue
		}

		w = r
		s[n] = b
		n++
	}

	return s[:n]
}

func main() {
	s := []byte("one two two   three  four  five  five   ")
	s1 := squashSpaces(s)
	fmt.Printf("%q\n%q\n", s, s1)
}
