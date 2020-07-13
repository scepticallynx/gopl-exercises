/*Exercise 5.13: Use breadthFirst function to explore a different structure.
For example: courses from topoSort (graph), filesystem (tree), bus routes (undirected graph).
*/
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func walker(node string) []string {
	var nextDirs []string

	info, err := ioutil.ReadDir(node)
	if err == nil {
		for _, d := range info {
			if d.IsDir() {
				nextDirs = append(nextDirs, filepath.Join(node, d.Name()))
			}
		}
	}

	// fmt.Printf("|%-*s%s\n", depth, )
	fmt.Printf("%s\n", node)

	return nextDirs
}

func main() {
	breadthFirst(walker, os.Args[1:])
}
