package main

import (
	"fmt"
	"testing"
)

func ExampleIntSet_UnionWith() {
	var x, y IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	fmt.Println(x.String()) // "{1 9 144}"

	y.Add(9)
	y.Add(42)
	fmt.Println(y.String()) // "{9 42}"

	x.UnionWith(&y)
	fmt.Println(x.String()) // "{1 9 42 144}"

	fmt.Println(x.Has(9), x.Has(123)) // "true false"

	// Output:
	// {1 9 144}
	// {9 42}
	// {1 9 42 144}
	// true false
}

func ExampleIntSet_Add() {
	var x IntSet
	x.Add(1)
	x.Add(144)
	x.Add(9)
	x.Add(42)

	fmt.Println(&x)         // "{1 9 42 144}"
	fmt.Println(x.String()) // "{1 9 42 144}"
	fmt.Println(x)          // "{[4398046511618 0 65536]}"

	// Output:
	// {1 9 42 144}
	// {1 9 42 144}
	// {[4398046511618 0 65536]}
}

func ExampleIntSet_Remove() {
	var x IntSet
	x.Add(1)
	x.Add(2)
	x.Add(3)
	x.Add(3)
	x.Add(9)
	x.Add(6)

	x.Remove(9)
	x.Remove(3)
	x.Remove(10)

	fmt.Println(&x)

	// Output:
	// {1 2 6}
}

func ExampleIntSet_Len() {
	var x IntSet
	x.Add(1)
	x.Add(2)
	x.Add(2)

	fmt.Println(x.Len())

	// Output:
	// 2
}

func ExampleIntSet_Len2() {
	var x IntSet
	x.Add(1)
	x.Add(2)
	x.Add(3)

	fmt.Println(x.Len())

	// Output:
	// 3
}

func ExampleIntSet_Copy() {
	var x IntSet
	x.Add(23)
	x.Add(32)

	y := x.Copy()

	x.Remove(32)

	fmt.Printf("%s\n%s\n", x.String(), y.String())

	// Output:
	// {23}
	// {23 32}
}

func TestIntSet_IntersectWith(t *testing.T) {
	var a, b IntSet
	want := `{2 3}`
	a.AddAll([]int{1, 2, 3}...)
	b.AddAll([]int{2, 3, 4}...)

	a.IntersectWith(&b)

	if a.String() != want {
		t.Errorf("Want: %s\tGot: %s\n", want, a.String())
	}
}

func TestIntSet_SymmetricDifference(t *testing.T) {
	var a, b IntSet
	want := `{1 5}`
	a.AddAll([]int{1, 2, 2, 2, 3, 3, 4}...)
	b.AddAll([]int{2, 3, 4, 4, 4, 5, 5}...)

	a.SymmetricDifference(&b)

	if a.String() != want {
		t.Errorf("Want: %s\tGot: %s\n", want, a.String())
	}
}

func TestIntSet_DifferenceWith(t *testing.T) {
	var a, b IntSet
	want := `{1}`
	a.AddAll([]int{1, 2, 2, 2, 3, 3, 4}...)
	b.AddAll([]int{2, 3, 4, 4, 4, 5, 5}...)

	a.DifferenceWith(&b)

	if a.String() != want {
		t.Errorf("Want: %s\tGot: %s\n", want, a.String())
	}
}

func TestIntSet_Elems(t *testing.T) {
	var a IntSet
	a.AddAll([]int{1, 2, 2, 2, 3, 3, 4}...)

	want := []int{1, 2, 3, 4}

	result := a.Elems()

	if lr, lw := len(result), len(want); lr != lw {
		t.Fatalf("Sets length differs: want %d\tgot%d\n", lw, lr)
	}

	for n := 0; n < len(result); n++ {
		if r, w := result[n], want[n]; r != w {
			t.Errorf("Elements #%3d differs: want %d\tgot%d\n", n, w, r)
		}
	}

	fmt.Println(result)
}

// Looks like variant with x > 0 && s.Has(x) check is the fastest
func BenchmarkIntSet_AddAll(b *testing.B) {
	var x IntSet
	s := []int{1, 2, 3, 4, 5, 6, 7, 8, 8, 8, 8, 9, 9, 9, -1}

	for n := 0; n < b.N; n++ {
		x.AddAll(s...)
	}
}
