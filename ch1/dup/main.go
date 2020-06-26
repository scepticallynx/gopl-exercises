package main

import (
	"bufio"
	"fmt"
	"os"
)

// Exercise 1.4: Modify dup2 to print the names of all files in which each duplicate line occurs.
//  run: echo -e "1\n1\n1\n2\n3\n2\n4\n5\n6" | go run ch1/dup/main.go
//  or: go run ch1/dup/main.go ch1/dup/dup1.txt ch1/dup/dup2.txt
func main() {
	counts := make(map[string]map[string]int)
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts)
			f.Close()
		}
	}
	for f, duplicates := range counts {
		fmt.Printf("Filename\tQ-ty\tString\n")
		for line, n := range duplicates {
			if n > 1 {
				fmt.Printf("%-8s\t%4d\t%6s\n", f, n, line)
			}
		}
	}
}

func countLines(f *os.File, counts map[string]map[string]int) {
	input := bufio.NewScanner(f)
	fName := f.Name()
	counts[fName] = map[string]int{}
	for input.Scan() {
		counts[f.Name()][input.Text()]++
	}
	// NOTE: ignoring potential errors from input.Err()
}
