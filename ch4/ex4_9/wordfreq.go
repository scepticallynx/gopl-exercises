// Exercise 4.9: Write a program wordfreq to report frequency of each word in an
// input text file.
package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	var words = make(map[string]int)

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		words[scanner.Text()]++
	}

	sorted := make([]string, 0, len(words))
	for w := range words {
		sorted = append(sorted, w)
	}

	sort.Strings(sorted)

	fmt.Printf("%-20s\tQ-ty\n", "Word")
	for _, w := range sorted {
		if n := words[w]; n > 1 {
			fmt.Printf("%-20q\t%-6d\n", w, n)
		}
	}
}
