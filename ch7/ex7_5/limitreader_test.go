package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestLimitReader(t *testing.T) {
	input, expect := `123 456 789`, `123`
	var readN int64 = 3

	var b bytes.Buffer

	reader := LimitReader(strings.NewReader(input), readN)

	n, _ := b.ReadFrom(reader)
	if n != readN {
		t.Errorf("Expected: %d\tGot: %d\n", readN, n)
	}

	if b.String() != expect {
		t.Errorf("Expected: %s\tGot: %s\n", expect, b.String())
	}
}
