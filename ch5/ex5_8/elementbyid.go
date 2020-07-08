/*Exercise 5.8: Modify forEachNode so that pre and post functions return a boolean
result indicating whether to continue the traversal. Use it to write a function
ElementByID with the following signature that finds the first HTML element with
the specified id attribute. The function should stop the traversal as soon as
a match found.
	func ElementByID(doc *html.Node, id string) *html.Node
*/
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

// idMatch answers whether forEachNode need traverse further.
// If node is not nil and is html.ElementNode, the function iterates attributes
// of current node in search of key "id", and if such key exists and its value
// equals specified id, false is returned. In all other cases function returns true.
func idMatch(node *html.Node, id string) bool {
	if node != nil && node.Type == html.ElementNode {
		for _, attr := range node.Attr {
			if attr.Key == "id" && attr.Val == id {
				return false // enough traversing
			}
		}
	}
	return true // traverse more
}

// weird game with local variables.
// It seems I really cannot get what exactly the task wants from me, otherwise
// I'd explained the function properly.
func forEachNode(node *html.Node, id string, pre, post func(*html.Node, string) bool) (n *html.Node) {
	// "node" is that node given to function in initial call,
	// in all other cases it is node provided by recursive calls:
	// any of its children or children of children...
	if !pre(node, id) {
		// n (result variable that has been nil recently) is set to current node
		// only if isMatch returns false inside pre.
		n = node
		return
	}

	// standard loop from book (and html package docs): as current node is already
	// checked, get its first child, than child's sibling and so on until any of
	// them is nil (end of branch) or has matched "id" value (n is set in pre,
	// recursive call returns "not nil n" and execution terminates
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		if _n := forEachNode(c, id, pre, post); _n != nil {
			// if n is set in one of recursive call, when traversing deeper into
			// child's or sibling's or children of child... branches (see "pre" block)
			// set "current" n and return
			n = _n
		}
	}

	// nil until isMatch inside pre returns false (might never happen)
	return
}

// ElementByID returns first tag with attribute id matching specified id
func ElementByID(doc *html.Node, id string) *html.Node {
	return forEachNode(doc, id, idMatch, nil)
}

const usage = `elementbyid [URL] [ID]`

func main() {
	log.SetPrefix("elementbyid: ")
	log.SetFlags(0)

	if len(os.Args) != 3 || os.Args[1] == "" || os.Args[2] == "" {
		fmt.Println(usage)
		os.Exit(1)
	}

	resp, err := http.Get(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("parsing html: %v", err)
	}

	node := ElementByID(doc, os.Args[2])
	if node != nil {
		fmt.Printf("%s\t%s\n", node.Data, node.Attr)
	} else {
		log.Printf("Node with id %q is not found in source of %s", os.Args[2], os.Args[1])
	}
}
