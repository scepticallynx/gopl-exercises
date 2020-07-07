/* Exercise 5.3: Write a function to print contents of all text nodes in HMTL
document tree. Do no descend into <style> and <script> elements, since their
contents are not visible in browser.
*/
package main

import (
	"fmt"
	"os"

	"golang.org/x/net/html"
)

// https://html.spec.whatwg.org/multipage/syntax.html#elements-2
// https://developer.mozilla.org/en-US/docs/Web/API/Node/nodeType#Node_type_constants
// Raw text elements: script, style
// ElementNode - <p>, <div> etc
// TextNode - actual text inside ElementNode or Attr
// ElementNode.FirstChild -> TextNode
//
// Experimenting with goto
func traverse(node *html.Node) {
	if node == nil {
		return
	}

	switch node.Type {
	case html.ElementNode:
		switch node.Data {
		case "script", "style":
			goto OnlySibling
		}
	case html.TextNode:
		_, _ = fmt.Fprintf(os.Stdout, "%s", node.Data)
	}

	traverse(node.FirstChild)
OnlySibling:
	traverse(node.NextSibling)
}

func main() {
	doc, err := html.Parse(os.Stdin)

	if err != nil {
		fmt.Fprintf(os.Stderr, "findlinks error: %v", err)
		os.Exit(1)
	}

	traverse(doc)
}
