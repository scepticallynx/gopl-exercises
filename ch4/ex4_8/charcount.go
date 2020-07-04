// Exercise 4.8: modify charcount to count letters, digits, and so on in their
// Unicode categories.
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode"
)

// charcount uses simple approach to group Unicode characters to avoid messing
// with unicode.RangeTable
//
// Digits are stored in map["digits"] because their condition comes first
func charcount() {
	counts := make(map[string]map[rune]int)

	in := bufio.NewReader(os.Stdin)
	for {
		r, n, err := in.ReadRune()
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "charcount: %v\n", err)
			os.Exit(1)
		}

		var key string
		switch {
		case unicode.IsDigit(r):
			key = "digits"
		case unicode.IsNumber(r):
			key = "numbers"
		case unicode.IsSymbol(r):
			key = "symbols"
		case unicode.IsMark(r):
			key = "marks"
		case unicode.IsPunct(r):
			key = "puncts"
		case r == unicode.ReplacementChar && n == 1:
			key = "invalid"
		case unicode.IsControl(r):
			key = "control"
		case unicode.IsLetter(r):
			key = "letters"
		case unicode.IsSpace(r):
			key = "space"
		default:
			key = "other"
			fmt.Printf("Other: %q", r)
		}

		if counts[key] == nil {
			counts[key] = make(map[rune]int)
		}

		counts[key][r]++
	}
	fmt.Printf("Category\tRune\tQ-ty\n")
	for category, rn := range counts {
		for r, n := range rn {
			if n > 0 {
				fmt.Printf("%q\t%q\t%6d\n", category, r, n)
			}
		}
	}
}

func main() {
	charcount()
}
