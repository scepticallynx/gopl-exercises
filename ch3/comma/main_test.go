package main

import "testing"

func BenchmarkComma(b *testing.B) {
	for n := 0; n < b.N; n++ {
		comma("abcdefghijklmno")
	}
}

func BenchmarkCommaBuffer(b *testing.B) {
	for n := 0; n < b.N; n++ {
		commaBuffer("abcdefghijklmno")
	}
}

func TestCommaEnhanced(t *testing.T) {
	testCases := []struct {
		In, Out string
	}{
		{"123456789", "123, 456, 789"},
		{"3334567751785", "333, 456, 775, 178, 5"},
		{"+234-5+1+2-34", "+234, -5+1+2, -34"},
		{"123.1-4.256+7.78-90", "123.1, -4.256, +7.78-9, 0"},
	}

	for n, test := range testCases {
		if result := commaEnhanced(test.In); result != test.Out {
			t.Errorf("Case #%2d: %s\nExpected: %s\nGot: %s", n, test.In, test.Out, result)
		}
	}
}

func BenchmarkCommaEnhanced(b *testing.B) {
	for n := 0; n < b.N; n++ {
		commaEnhanced("123.1-4.256+7.78-90")
	}
}
