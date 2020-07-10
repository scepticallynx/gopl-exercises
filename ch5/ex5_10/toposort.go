/*Exercise 5_10: Rewrite topoSort to use maps instead of slices and eliminate the
initial sort. Verify that the results, though nondeterministic, are valid
topological orderings.
*/
package main

import (
	"fmt"
	"sort"
)

// prereqs maps computer science courses to their prerequisites.
// Elements are intentionally unsorted to match the exercise requirements.
var prereqs = map[string][]string{
	"calculus":              {"linear algebra"},
	"programming languages": {"data structures", "computer organization"},
	"formal languages":      {"discrete math"},
	"data structures":       {"discrete math"},
	"databases":             {"data structures"},
	"algorithms":            {"data structures"},
	"discrete math":         {"intro to programming"},
	"operating systems":     {"data structures", "computer organization"},
	"networks":              {"operating systems"},
	"compilers": {
		"data structures",
		"formal languages",
		"computer organization",
	},
}

func main() {
	for i, course := range topoSortMap(prereqs) {
		fmt.Printf("%d:\t%s\n", i, course)
	}
}

// topoSortMap is almost the same as topoSort (from the book example). Instead
// of returning slice of ordered strings, it returns a map with string order to string.
// Results are nondeterministic but items has the right topological order.
func topoSortMap(m map[string][]string) map[int]string {
	var n int
	order := make(map[int]string)

	seen := make(map[string]bool)
	var visitAll func(items []string)

	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				visitAll(m[item])
				n++
				order[n] = item
			}
		}
	}

	// here we do need to sort items
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	visitAll(keys)

	return order
}
