package main

import (
	"fmt"
	"os"
	"testing"
)

func TestCountingWriter(t *testing.T) {
	w, n := CountingWriter(os.Stderr)

	fmt.Fprintln(w, "123 456 abc")

	fmt.Println(*n)
}
