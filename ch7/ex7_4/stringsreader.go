/* Exercise 7.4: ... Implement a simple version of strings.NewReader, and use it
to make the HTML parser (5.2) take input from string.
*/
package main

import (
	"fmt"
	"io"
	"log"

	"golang.org/x/net/html"
)

type StringsReader struct {
	s string
}

func (s *StringsReader) Read(b []byte) (n int, err error) {
	n = copy(b, s.s)
	s.s = s.s[n:]
	if len(s.s) == 0 {
		err = io.EOF
	}

	return n, err
}

func NewReader(s string) *StringsReader {
	return &StringsReader{s}
}

func main() {
	reader := NewReader(`<html><body><head><title>The title</title></head></body></html>`)

	doc, err := html.Parse(reader)
	if err != nil {
		log.Fatal(err)
	}

	var printHTML func(n *html.Node)

	printHTML = func(n *html.Node) {
		fmt.Println(n.Data)

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			printHTML(c)
		}
	}

	printHTML(doc)
}
