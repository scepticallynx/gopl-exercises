package main

import (
	"fmt"
	"testing"
)

func TestTree_String(t *testing.T) {
	var tree tree

	for i := 0; i <= 50; i++ {
		add(&tree, i)
	}

	fmt.Println(tree.String())
}
