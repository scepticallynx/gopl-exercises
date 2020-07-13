/*Exercise 5.17: Write a variadic function ElementsByTagName that, given an HTML
node tree and zero or more names, returns all the elements that match one of those
 names.
Example:
	func ElementsByTagName(doc *html.Node, name ...string) []*html.Node

	images := ElementsByTagName(doc, "img")
	headings := ElementsByTagName(doc, "h1", "h2", "h3", "h4")
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

func ElementsByTagName(doc *html.Node, name ...string) []*html.Node {
	var matched []*html.Node
	var forEachNode func(node *html.Node)

	forEachNode = func(n *html.Node) {
		if n.Type == html.ElementNode {
			for _, tag := range name {
				if n.Data == tag {
					matched = append(matched, n)
				}
			}
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			forEachNode(c)
		}
	}

	forEachNode(doc)

	return matched
}

func fetch(url string) (*html.Node, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}

	return doc, nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: elementsbytagname url tag1 ... tagN")
		os.Exit(1)
	}

	url := os.Args[1]
	tags := os.Args[2:]

	doc, err := fetch(url)
	if err != nil {
		log.Fatalf("%q: %v\n", url, err)
	}

	match := ElementsByTagName(doc, tags...)
	if l := len(match); l > 0 {
		fmt.Printf("%q: %d tags matched\n", url, l)
		for _, el := range match {
			fmt.Printf("%s, ", el.Data)
		}
		fmt.Println()
	}
}
