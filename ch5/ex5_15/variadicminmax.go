/*Exercise 5.15: Write variadic functions max and min, analogous to sum. What
should these functions do when called with no arguments? Write variants that
require at least one argument.
*/
package main

import (
	"fmt"
	"log"
)

func max(digits ...int) int {
	var m int

	switch len(digits) {
	case 0:
		log.Fatalf("no arguments given to function. Need at least 1 digit")
	case 1:
		return digits[0]
	default:
		for _, n := range digits {
			if n > m {
				m = n
			}
		}
	}

	return m
}

func min(digits ...int) int {
	var m int

	switch len(digits) {
	case 0:
		log.Fatalf("no arguments given to function. Need at least 1 digit")
	case 1:
		return digits[0]
	default:
		for _, n := range digits {
			if n < m {
				m = n
			}
		}
	}

	return m
}

func main() {
	digits := []int{1, 2, 3, 4, 5, 5, 5, 6, 6, 7, 9, -2}
	fmt.Println(max(digits...), min(digits...))
	max()
}
