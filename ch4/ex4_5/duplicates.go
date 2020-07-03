// Exercise 4.5: Write an in-place function to eliminate adjacent duplicates in a
// []string slice. In-place - use the same underlying array for input and output slices.
package main

import "fmt"

func eliminateDuplicates(s []string) []string {
	var previous string

	i := 0
	for n, str := range s {
		if n > 0 {
			previous = s[n-1]
		}

		if previous != str {
			s[i] = str
			i++
		}
	}

	return s[:i]
}

func main() {
	// s := []string{"f", "a", "a", "a", "a", "f", "f", "c"}
	var s = []string{"simle", "word", "word", "word", "word", "test", "test", "simple"}
	s = eliminateDuplicates(s)
	fmt.Println(s)
}
