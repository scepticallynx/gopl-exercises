/*Exercise 5.11: The instructor of the linear algebra course decides that calculus
is now a prerequisite. Extend the topoSort function to report cycles.
*/
package main

import (
	"fmt"
	"sort"
)

// prereqs maps computer science courses to their prerequisites.
var prereqs = map[string][]string{
	"linear algebra":        {"calculus"},
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
	for i, course := range topoSort(prereqs) {
		fmt.Printf("%d:\t%s\n", i+1, course)
	}
}

func topoSort(m map[string][]string) []string {
	var order []string
	seen := make(map[string]bool)
	var visitAll func(items []string)

	inSlice := func(s []string, item string) (string, bool) {
		for _, i := range s {
			if i == item {
				return i, true
			}
		}

		return "", false
	}

	// example: visitAll receives "linear algebra", its prerequisite is "calculus".
	// "linear algebra" is not in seen, add it and call recursively visitAll with
	// "calculus".
	// visitAll receives "calculus", its prerequisite is "linear algebra".
	// "linear algebra" is already processed as course, now it's a prerequisite.
	visitAll = func(items []string) {
		for _, item := range items {
			if !seen[item] {
				seen[item] = true

				visitAll(m[item])
				order = append(order, item)
			} else {
				for _, v := range m[item] {
					if _, ok := inSlice(m[v], item); ok {
						fmt.Printf("Cycle: %s <-> %s\n", item, v)
						break
					}
				}
			}
		}
	}

	var keys []string
	for key := range m {
		keys = append(keys, key)
	}

	sort.Strings(keys)
	visitAll(keys)
	return order
}
