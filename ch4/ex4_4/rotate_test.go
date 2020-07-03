package main

import "testing"

var s = []int{0, 1, 2, 3, 4, 5}

func BenchmarkRotate(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rotate(s, 2)
	}
}

func BenchmarkRotate2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		rotate2(s, 2)
	}
}
