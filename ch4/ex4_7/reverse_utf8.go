// Exercise 4.8: Modify reverse to reverse the characters of a []byte slice that
// represents a UTF-8 encoded string, in place. Can you do it without allocating new memory?
package main

import (
	"fmt"
	"unicode/utf8"
)

// string(s) already allocates new memory
func reverseUTF8(s []byte) {
	for i := 0; i < len(s); {
		_, size := utf8.DecodeRune(s[i:])
		// reverse
		for j := 0; j < len(s[i:i+size])/2; j++ {
			s[j+i], s[len(s[j+i:j+i+size])-j+i-1] = s[len(s[j+i:j+i+size])-j+i-1], s[j+i]
		}
		i += size
	}

	// reverse
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-i-1] = s[len(s)-i-1], s[i]
	}
}

func main() {
	s := []byte("Привет, 世界")
	reverseUTF8(s)
	fmt.Printf("%s\n", s)
}
