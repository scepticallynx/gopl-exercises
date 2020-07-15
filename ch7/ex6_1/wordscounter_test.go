package main

import (
	"fmt"
	"testing"
)

func TestWordsLinesCounter_Write(t *testing.T) {
	testCases := []struct {
		input  []byte
		expect [2]int // lines, words
	}{
		{input: []byte("1"), expect: [2]int{1, 1}},
		{input: []byte("one two three 123\nfour 5 6"), expect: [2]int{2, 7}},
		{input: []byte("   \n    \n   3"), expect: [2]int{3, 1}},
	}

	for n, test := range testCases {
		var c WordsLinesCounter
		fmt.Fprintf(&c, "%s", test.input)

		if c.Totals() != test.expect {
			t.Errorf("Case #%2d\tInput: %s\tExpected: %v\tGot: %v\n", n, test.input, test.expect, c.Totals())
		}
	}
}
