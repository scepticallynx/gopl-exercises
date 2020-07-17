package main

import (
	"testing"
)

func TestIsPalindrome(t *testing.T) {
	testCases := []struct {
		in  somethingSortable
		out bool
	}{
		{[]rune("malayalam"), true},
		{[]rune("abba"), true},
		{[]rune("test"), false},
		{[]rune(""), false},
		{[]rune("1"), false},
		{[]rune("1234567890987654321"), true},
	}

	for n, test := range testCases {
		if result := IsPalindrome(test.in); result != test.out {
			t.Errorf("%2d: IsPalindrome(%s) != %t", n, string(test.in), test.out)
		}
	}
}
