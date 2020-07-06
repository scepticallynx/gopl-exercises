/* Exercise 5.1: Change the findlinks to traverse the n.FirstChild linked list
using recursive calls to visit instead of a loop.

How to run:
1) go build -o fetch ../../ch1/fetch/main.go && go build findlinks.go
2) ./fetch golang.org | ./findlinks
*/
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

func main() {
	doc, err := html.Parse(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks error: %v", err)
		os.Exit(1)
	}

	for _, link := range visit(nil, doc) {
		fmt.Println(link)
	}
}

func visit(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				links = append(links, a.Val)
			}
		}
	}

	if c := n.FirstChild; c != nil {
		links = visit(links, c)
	}

	if c := n.NextSibling; c != nil {
		links = visit(links, c)
	}

	return links
}
