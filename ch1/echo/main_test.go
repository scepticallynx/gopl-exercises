package main

import "testing"

func BenchmarkEcho1(t *testing.B) {
	for i:=0; i<t.N; i++ {
		echo1()
	}
}

func BenchmarkEcho2(t *testing.B) {
	for i:=0; i<t.N; i++ {
		echo2()
	}
}

func BenchmarkEcho3(t *testing.B) {
	for i:=0; i<t.N; i++ {
		echo3()
	}
}