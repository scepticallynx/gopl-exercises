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

	Exercise 3.12: Write a function that reports whether two strings are anagrams
	of each other.
	See benchmarks: cd ch3/comma; go test -bench=.
*/
package main

import (
	"bytes"
	"fmt"
	"os"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
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

// ---- Exercise 3.12 ----

// isAnagram answers whether 2 given strings are anagrams.
// Current approach used just for the sake of trying something new and different from
// sorting letters method, mapping letters to q-ty of occurrences etc.
func isAnagram(s1, s2 string) bool {
	if s1 == s2 {
		return false
	}

	n1, n2 := len(s1), len(s2)
	if n1 != n2 {
		return false
	}

	f := func(s string) int32 {
		var result int32

		for _, c := range s {
			if unicode.IsSpace(c) {
				continue
			}

			if !unicode.IsLower(c) {
				c = unicode.ToLower(c)
			}

			result += c
		}

		return result
	}

	if f(s1) != f(s2) {
		return false
	}

	return true
}

// isAnagramMap checks whether two strings are anagrams counting occurrences of their letters.
func isAnagramMap(s1, s2 string) bool {
	if s1 == s2 || utf8.RuneCountInString(s1) != utf8.RuneCountInString(s2) {
		return false
	}

	m1, m2 := make(map[rune]int), make(map[rune]int)

	for _, c := range s1 {
		m1[c]++
	}

	for _, c := range s2 {
		m2[c]++
	}

	for r, n := range m1 {
		if m2[r] != n {
			return false
		}
	}

	return true
}

// isAnagramSort sorts letters of given string to answer whether they are anagrams of each other.
func isAnagramSortRune(s1, s2 string) bool {
	if s1 == s2 || len(s1) != len(s2) {
		return false
	}

	a1 := []rune(s1)
	a2 := []rune(s2)

	sort.Slice(a1, func(i, j int) bool {
		return a1[i] < a1[j]
	})

	sort.Slice(a2, func(i, j int) bool {
		return a2[i] < a2[j]
	})

	for n, r := range a1 {
		if a2[n] != r {
			return false
		}
	}

	return true
}

// isAnagramSortString also sorts strings slices of string not runes
func isAnagramSortString(s1, s2 string) bool {
	if s1 == s2 || len(s1) != len(s2) {
		return false
	}

	ss1 := strings.Split(s1, "")
	ss2 := strings.Split(s2, "")

	sort.Strings(ss1)
	sort.Strings(ss2)

	if strings.Join(ss1, "") != strings.Join(ss2, "") {
		return false
	}

	return true
}
