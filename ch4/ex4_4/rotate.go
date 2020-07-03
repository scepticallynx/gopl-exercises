// Exercise 4.3: Write a version of rotate that operates in single pass.
package main

import (
	"fmt"
)

// rotate rotates given slice left by n elements
func rotate(s []int, n int) {
	l := len(s)

	// switch {
	// case l == 0, n == 0, n == l:
	// 	return
	// case n > l, n < 0: // uint for n ?
	// 	log.Fatalf("number of elements n=%d is bigger than length of slice %d or < 0", n, l)
	// default:
	temp := make([]int, l, cap(s))

	for i, d := range s {
		if i < n {
			temp[l-n+i] = d
		} else {
			temp[i-n] = d
		}
	}

	copy(s, temp)
	// }
}

// simpler but slower version of rotate. Anyway original approach (3 times rotate) is the fastest
func rotate2(s []int, n int) {
	temp := append(s, s[:n]...)
	copy(s, temp[n:])
}

func main() {
	s := []int{0, 1, 2, 3, 4, 5}
	rotate(s, 2)
	fmt.Println(s)
}
