package popcount

import "testing"

func BenchmarkPopCount(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PopCount(10)
	}
}

func BenchmarkPopCountLoop(b *testing.B) {
	for n := 0; n < b.N; n++ {
		PopCountLoop(10)
	}
}
