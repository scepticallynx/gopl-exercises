/*Exercise 6.1: Using the ideas from ByteCounter, implement counters for words
and for lines.
*/
package main

import (
	"bufio"
	"bytes"
	"fmt"
)

// WordsLinesCounter counts number of line and number of words in each line
type WordsLinesCounter struct {
	words, lines int
}

// countWords increments WordsLinesCounter.words field by number of words found in line p
func (c *WordsLinesCounter) countWords(p []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(p))

	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		c.words++
	}
}

// countLines increments WordsLinesCounter.lines field by number of lines found in input p
func (c *WordsLinesCounter) countLines(p []byte) {
	scanner := bufio.NewScanner(bytes.NewReader(p))

	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		c.lines++

		// for each line count number of words in it
		c.countWords(scanner.Bytes())
	}
}

func (c *WordsLinesCounter) Write(p []byte) (int, error) {
	c.countLines(p)

	return len(p), nil
}

func (c *WordsLinesCounter) Words() int {
	return c.words
}

func (c *WordsLinesCounter) Lines() int {
	return c.lines
}

func (c *WordsLinesCounter) Totals() [2]int {
	return [2]int{c.Lines(), c.Words()}
}

func main() {
	var c WordsLinesCounter

	input := []byte("one two three 123\nfour 5 6")

	fmt.Fprintf(&c, "%s", input)
	fmt.Println(c)
}
