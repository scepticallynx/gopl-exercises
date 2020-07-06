/* Exercise 5.2: Write a function to populate a mapping from element names - p,
div, span, and so on - to the number of elements with that name in an HTML document tree.
*/
package main

import (
	"fmt"
	"os"
	"sort"

	"golang.org/x/net/html"
)

// countElements traverses HTML tree and increments each tag quantity in m
func countElements(m map[string]int, node *html.Node) map[string]int {
	if node.Type == html.ElementNode {
		m[node.Data]++
	}

	if c := node.FirstChild; c != nil {
		m = countElements(m, c)
	}

	if c := node.NextSibling; c != nil {
		m = countElements(m, c)
	}

	return m
}

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks error: %v", err)
		os.Exit(1)
	}

	// a map for elements
	elements := make(map[string]int, 0)

	// populating the map
	elements = countElements(elements, doc)

	// sort by quantity, descending
	byN := make(map[int][]string, len(elements))
	for e, n := range elements {
		if v, exists := byN[n]; !exists {
			byN[n] = []string{e}
		} else {
			v = append(v, e)
			byN[n] = v
		}
	}

	// slice of numbers of elements to be sorted
	sorted := make([]int, 0, len(byN))

	for n := range byN {
		sorted = append(sorted, n)
	}

	sort.Sort(sort.Reverse(sort.IntSlice(sorted)))

	for _, n := range sorted {
		for _, el := range byN[n] {
			fmt.Printf("%-10s\t%4d\n", el, n)
		}
	}
}
