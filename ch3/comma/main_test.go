package main

import "testing"

var cases = []struct {
	S1, S2 string
	want   bool
}{
	{S1: "listen", S2: "silent", want: true},
	{S1: "test", S2: "testify", want: false},
	{S1: "eleven plus two", S2: "twelve plus one", want: true},
	{"bad credit", "debit card", true},
	{"act", "cat", true},
	{"abcddefg", "cd def gab", false}, // spaces
}

// testing

func TestIsAnagram(t *testing.T) {
	for n, test := range cases {
		if result := isAnagram(test.S1, test.S2); result != test.want {
			t.Errorf("Case #%2d: 1: %q 2: %q Want: %t Got %t\n", n, test.S1, test.S2, test.want, result)
		}
	}
}

func TestIsAnagramMap(t *testing.T) {
	for n, test := range cases {
		if result := isAnagramMap(test.S1, test.S2); result != test.want {
			t.Errorf("Case #%2d: 1: %q 2: %q Want: %t Got %t\n", n, test.S1, test.S2, test.want, result)
		}
	}
}

func TestIsAnagramSortRune(t *testing.T) {
	for n, test := range cases {
		if result := isAnagramSortRune(test.S1, test.S2); result != test.want {
			t.Errorf("Case #%2d: 1: %q 2: %q Want: %t Got %t\n", n, test.S1, test.S2, test.want, result)
		}
	}
}

func TestIsAnagramSortString(t *testing.T) {
	for n, test := range cases {
		if result := isAnagramSortString(test.S1, test.S2); result != test.want {
			t.Errorf("Case #%2d: 1: %q 2: %q Want: %t Got %t\n", n, test.S1, test.S2, test.want, result)
		}
	}
}

// benchmarking

var b1, b2 = cases[3].S1, cases[3].S2

func BenchmarkIsAnagram(b *testing.B) {
	for n := 0; n < b.N; n++ {
		isAnagram(b1, b2)
	}
}

func BenchmarkIsAnagramMap(b *testing.B) {
	for n := 0; n < b.N; n++ {
		isAnagram(b1, b2)
	}
}

func BenchmarkIsAnagramSortRune(b *testing.B) {
	for n := 0; n < b.N; n++ {
		isAnagram(b1, b2)
	}
}

func BenchmarkIsAnagramSortString(b *testing.B) {
	for n := 0; n < b.N; n++ {
		isAnagram(b1, b2)
	}
}
