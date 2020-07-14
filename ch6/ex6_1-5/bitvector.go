/*Package intset provides a set of integers based on a bit vector.

Exercise 6.1: Implement this additional methods:
	func (*IntSet) Len() int // return the number of elements
	func (*IntSet) Remove(x int) // remove x from set
	func (*IntSet) Clear() // remove all elements from the set
	func (*IntSet) Copy() *IntSet // return a copy of a set

Exercise 6.2: Define a variadic (*IntSet).AddAll(...int) method that allows a list
of values to be added, such as s.AddAll(1, 2, 3).

Exercise 6.3: ... Implement methods for IntersectWith, DifferenceWith, and
SymmetricDifference for the corresponding set operations.

Exercise 6.4: Add the method Elems that returns a slice containing the elements
of the set, suitable for iterating over with a range loop.

Exercise 6.5: ... Modify the program to use uint type, which is the most efficient
unsigned integer type for the platform ...
*/
package main

import (
	"bytes"
	"fmt"
)

const effectiveUintSize = 32 << (^uint(0) >> 63)

// An IntSet is a set of small non-negative integers.
// Its zero value represents the empty set.
//
// https://en.wikipedia.org/wiki/Bit_array
// To obtain the bit mask needed for Add, Remove, Clear... operations, we can use
// a bit shift operator to shift the number 1 to the left by the appropriate number
// of places, as well as bitwise negation if necessary.
type IntSet struct {
	words []uint
}

// Has reports whether the set contains the non-negative value x.
//
// AND together with zero-testing can be used to determine if a bit is set:
//    11101010 AND 00000001 = 00000000 = 0
//    11101010 AND 00000010 = 00000010 ≠ 0
func (s *IntSet) Has(x int) bool {
	word, bit := x/effectiveUintSize, uint(x%effectiveUintSize)
	return word < len(s.words) && s.words[word]&(1<<bit) != 0
}

// Add adds the non-negative value x to the set.
//
// OR can be used to set a bit to one:
//  11101010 OR 00000100 = 11101110
func (s *IntSet) Add(x int) {
	word, bit := x/effectiveUintSize, uint(x%effectiveUintSize) // index in slice, bit index
	for word >= len(s.words) {
		s.words = append(s.words, 0)
	}
	s.words[word] |= 1 << bit
}

// AddAll adds all positive integers from given slice, calling Add for each if
// int is not yet in s
func (s *IntSet) AddAll(xs ...int) {
	for _, x := range xs {
		if x > 0 && !s.Has(x) {
			s.Add(x)
		}
	}
}

// Remove removes the non-negative value x from the set.
//
// AND can be used to set a bit to zero:
//  11101010 AND 11111101 = 11101000
// NOT can be used to invert all bits
// &^ - AND (NOT )
func (s *IntSet) Remove(x int) {
	// if s.Has(x) {
	word, bit := x/effectiveUintSize, uint(x%effectiveUintSize)
	s.words[word] &^= 1 << bit
	// }
}

// Len is popcount
func (s *IntSet) Len() int {
	n := 0
	for _, word := range s.words {
		// popcount
		for word != 0 {
			n++
			word &= word - 1
		}
	}
	return n
}

// Copy returns a copy of s
func (s *IntSet) Copy() *IntSet {
	_s := new(IntSet)
	_s.words = make([]uint, len(s.words))
	copy(_s.words, s.words)
	return _s
}

// UnionWith sets s to the union of s and t.
// https://en.wikipedia.org/wiki/Intersection_(set_theory)
//
//  for i from 0 to n/w-1
//      union[i] := a[i] or b[i]
func (s *IntSet) UnionWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] |= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// IntersectWith sets s to the intersection of s and t (elements that exist in both sets)
// https://en.wikipedia.org/wiki/Intersection_(set_theory)
//
//  for i from 0 to n/w-1
//      intersection[i] := a[i] and b[i]
func (s *IntSet) IntersectWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// DifferenceWith (complement) returns all elements from a not in b.
func (s *IntSet) DifferenceWith(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] &^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// SymmetricDifference returns elements unique to a and unique b sets.
func (s *IntSet) SymmetricDifference(t *IntSet) {
	for i, tword := range t.words {
		if i < len(s.words) {
			s.words[i] ^= tword
		} else {
			s.words = append(s.words, tword)
		}
	}
}

// Elems returns a slice containing the elements of the set.
func (s *IntSet) Elems() (result []int) {
	for i, word := range s.words {
		if word == 0 {
			continue
		}

		for j := 0; j < effectiveUintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				result = append(result, 64*i+j)
			}
		}
	}

	return
}

// String returns the set as a string of the form "{1 2 3}".
func (s *IntSet) String() string {
	var buf bytes.Buffer
	buf.WriteByte('{')
	for i, word := range s.words {
		if word == 0 {
			continue
		}
		for j := 0; j < effectiveUintSize; j++ {
			if word&(1<<uint(j)) != 0 {
				if buf.Len() > len("{") {
					buf.WriteByte(' ')
				}
				fmt.Fprintf(&buf, "%d", 64*i+j)
			}
		}
	}
	buf.WriteByte('}')
	return buf.String()
}

func main() {
	var x IntSet

	x.Add(1)
	x.Add(2)
}
