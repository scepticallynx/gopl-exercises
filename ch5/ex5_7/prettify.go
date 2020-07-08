/*Exercise 5.7: Develop startElement and endElement (chapter 5.5) into a general
HTML pretty-printer. Print comment nodes, text nodes, and the attributes of each
element (<a href='...'>). Use short forms like <img/> instead of <img></img> when
an element has no children. Write a test to ensure that the output can be parsed
successfully.

How to run and test:
1) download some url (with or without http|s://): go run prettify.go gopl.io
	This step will produce:
		test.html - a file with prettified HTML
		nodes_count.json - a file with mapped nodes to their number of occurrences
		original HTML
2) (in ex5_7 directory) go test . or go test -v .
	If prettified HTML contains any nodes that does not exist in original HTML
	or missing any nodes from original HTML, or quantity of some nodes does not
	match, you will see a detailed error message (testing.T error)

Some results:
gopl.io is fine.
google.com and duckduck.go fails only by quantity of <script> tags. Maybe they are
empty in original HTML or, most likely, my script failed with Text|RawNode type.
*/
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// prettified HTML output
var b bytes.Buffer

var depth int

func startElement(node *html.Node) {
	switch node.Type {
	case html.DoctypeNode: // looks like DOCTYPE needs special handling
		var attrs = ""
		for _, attr := range node.Attr {
			if attr.Key == "public" {
				attrs += "PUBLIC \"" + attr.Val + "\""
			} else {
				attrs += " \"" + attr.Val + "\""
			}
		}

		b.WriteString(fmt.Sprintf("<!DOCTYPE %s %s>\n", node.Data, attrs))
	case html.ElementNode:
		// collect attributes into single string
		var attrBuff strings.Builder

		for _, attr := range node.Attr {
			attrBuff.WriteString(" " + attr.Key + "='" + attr.Val + "'")
		}

		if node.FirstChild == nil {
			b.WriteString(fmt.Sprintf("%*s<%s%s/>\n", depth*2, "", node.Data, attrBuff.String()))
		} else {
			b.WriteString(fmt.Sprintf("%*s<%s%s>\n", depth*2, "", node.Data, attrBuff.String()))
			depth++ // increment depth only if tag has children
		}
	case html.TextNode, html.RawNode:
		if s := strings.TrimSpace(node.Data); len(s) > 0 { // skip empty rows
			b.WriteString(fmt.Sprintf("%*s%s\n", depth*2, "", s))
		}
	case html.CommentNode:
		b.WriteString(fmt.Sprintf("%*s<!-- %s -->\n", depth*2, "", node.Data))
	case html.ErrorNode: // dump error nodes to stderr, to compare q-ty of errors in test
		_, _ = fmt.Fprintf(os.Stderr, "%s\n", node.Data)
	}
}

func endElement(node *html.Node) {
	if node.Type == html.ElementNode {
		depth--
		b.WriteString(fmt.Sprintf("%*s</%s>\n", depth*2, "", node.Data))
	}
}

// forEachNode calls the functions pre(x) and post(x) for each node
// x in the tree rooted at n. Both functions are optional.
// pre is called before the children are visited (preorder) and
// post is called after (postorder).
func forEachNode(node *html.Node, pre, post func(node *html.Node)) {
	if pre != nil {
		pre(node)
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}

	if post != nil && node.FirstChild != nil { // do not call post for tags without children
		post(node)
	}
}

// outline is slightly modified to write down a map with numbers of nodes to be
// compared later in test.
func outline(url string) error {
	// add https:// prefix to url like in fetch from ch1
	if !(strings.HasPrefix(url, "http://") && strings.HasPrefix(url, "https://")) {
		url = "https://" + url
	}

	// query given url for HTML
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	// parse source of the queried url
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return err
	}

	// for testing purposes: count each node and store as json file
	nodes := make(map[string]int, 0)
	nodes = countNodes(nodes, doc)

	f, err := os.Create("nodes_count.json")
	if err != nil {
		return fmt.Errorf("dumping json with nodes quantity: %v", err)
	}

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(nodes)
	if err != nil {
		log.Fatalf("dumping json data: %v", err)
	}

	// recursively traverse into each tag of the source
	forEachNode(doc, startElement, endElement)

	err = ioutil.WriteFile("test.html", b.Bytes(), os.ModePerm)
	if err != nil {
		return fmt.Errorf("writing html: %v", err)
	}

	return nil
}

// countNodes is a function like in exercise 5.3. We use it to count number of
// occurrences in the original HTML and compare results with parsing generated HTML.
func countNodes(nodes map[string]int, n *html.Node) map[string]int {
	switch n.Type {
	case html.DoctypeNode:
		nodes["doctype"]++
	case html.ElementNode:
		nodes[n.Data]++
	case html.CommentNode:
		nodes["comment"]++
	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) > 0 {
			nodes["text"]++
		}
	case html.RawNode:
		nodes["raw"]++
	case html.ErrorNode:
		nodes["error"]++
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = countNodes(nodes, c)
	}

	return nodes
}

func main() {
	log.SetPrefix("outline: ")
	log.SetFlags(0)

	if len(os.Args) == 1 {
		log.Printf("No URL given")
		os.Exit(1)
	}

	// We need to test resulting HTML somehow, this version works with
	// single url, because after every run 2 files will be created:
	// nodes_count.json - with node name and quantity of this node occurrences
	// test.html - prettified HTML
	err := outline(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
}
