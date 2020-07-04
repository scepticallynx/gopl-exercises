package main

import "testing"

var s = []string{"simple", "word", "word", "word", "word", "test", "test", "simple"}

func BenchmarkEliminateDuplicates(b *testing.B) {
	for n := 0; n < b.N; n++ {
		eliminateDuplicates(s)
	}
}
