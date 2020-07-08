package main

import (
	"encoding/json"
	"log"
	"os"
	"testing"

	"golang.org/x/net/html"
)

// Since html.Parse can silently parse most of errors, let's try to parse
// generated HTML and print only html.ErrorNode (no results)
//
// Another possible test is to count quantity of nodes in original HTML and compare
// it with nodes in prettified HTML.
func TestPrettify(t *testing.T) {
	file, err := os.Open("test.html")
	if err != nil {
		log.Fatalf("missing prettified html (test.html): %s", err)
	}

	doc, err := html.Parse(file)
	if err != nil {
		t.Fatalf("parsing prettified HTML failed")
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ErrorNode {
			t.Errorf("error node in prettified HTML: %s", n.Data)
		}

		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)

	// local storage for nodes unmarshalling
	prettifiedNodes := make(map[string]int, 0)
	prettifiedNodes = countNodes(prettifiedNodes, doc)

	file, err = os.Open("nodes_count.json")
	if err != nil {
		log.Fatal(err)
	}

	var originalNodes map[string]int
	err = json.NewDecoder(file).Decode(&originalNodes)
	if err != nil {
		log.Fatal(err)
	}

	// compare existence of nodes from original HTML and number of their occurrences
	// with nodes from prettified HTML
	for node, quantity := range originalNodes {
		if quantityPrettified, exists := prettifiedNodes[node]; !exists {
			t.Errorf("missing node in prettified HTML: %q", node)
		} else if quantity != quantityPrettified {
			t.Errorf("different quantity of node %q occurrences in original (%d) and prettified (%d) HTML",
				node, quantity, quantityPrettified)
		}
	}

	// compare existence of nodes from prettified HTML in original HTML
	for node := range prettifiedNodes {
		if _, exists := originalNodes[node]; !exists {
			t.Errorf("extra node in prettified HTML: %s", node)
		}
	}
}
