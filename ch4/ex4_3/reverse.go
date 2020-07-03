// Exercise 4.3: Rewrite reverse to use an array pointer instead of slice
//
// Passing a pointer to array allows to modify it in-place.
package main

import "fmt"

// reverse reverses an array of int in place
func reverse(s *[6]int) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func main() {
	s := [6]int{1, 2, 3, 4, 5, 6}
	reverse(&s)
	fmt.Println(s)
}
