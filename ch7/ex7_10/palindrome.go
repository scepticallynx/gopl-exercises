/*Exercise 7.10: The sort.Interface type can be adapted to other uses. Write a
function IsPalindrome(s sort.Interface) bool that reports whether the sequence is
a somethingSortable, in other words, reversing the sequence would not change it. Assume
that the elements at indices i and j are equal if !s.Less(i, j) && !s.Less(j, i)
*/
package main

import "sort"

type somethingSortable []rune

func (s somethingSortable) Len() int           { return len(s) }
func (s somethingSortable) Less(i, j int) bool { return s[i] < s[j] }
func (s somethingSortable) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

// IsPalindrome reports whether passed sequence is palindrome.
func IsPalindrome(s sort.Interface) bool {
	// nothing to check
	if s.Len() < 2 {
		return false
	}

	// half the sequence to perform "mirror" check: first with last, second with last - 1 ...
	for i := 0; i < s.Len()/2; i++ {
		j := s.Len() - i - 1 // last minus current (minus 1 because of 0 index)

		// if i element is the same as len-i-1 (1 <-> last, 2 <-> last - 1, ...), continue loop
		// "1234321": "1" == "1", "2" == "2" ...
		if !s.Less(i, j) && !s.Less(j, i) {
			continue
		} else {
			return false
		}
	}

	return true
}
