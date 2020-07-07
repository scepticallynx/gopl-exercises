/* Exercise 5.4: Modify findlinks to print all links from img, style etc.

0) go build alllinks.go
1) any command to get html into stdin:
	./fetch gopl.io
	curl ...
	or even: cat local.html
2) redirect to (|)
3) ./alllinks

{HTML to stdin} | ./alllinks
*/
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

// https://html.spec.whatwg.org/index.html#attributes-1
func visit(links []string, node *html.Node) []string {
	if node == nil {
		return links
	}

	if node.Type == html.ElementNode {
		var attr string

		switch node.Data {
		case "a", "link", "area", "base": // href
			// a, area -> ping
			attr = "href"
		case "audio", "embed", "iframe", "img", "input", "script", "source", "track", "video": // src
			// video -> poster
			attr = "src"
		case "form": // action
			attr = "action"
		case "object": // data
			attr = "object"
		default:
			break
		}

		if attr != "" {
			for _, a := range node.Attr {
				if a.Key == attr {
					links = append(links, a.Val)
				}
			}
		}
	}

	links = visit(links, node.FirstChild)
	links = visit(links, node.NextSibling)

	return links
}

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
