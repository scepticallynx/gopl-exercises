/*Exercise 5.5: Implement countWordsAndImages
 */
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

// CountWordsAndImages does an HTTP GET request for the HTML document url and
// returns the number of words and images in it.
func CountWordsAndImages(url string) (words, images int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		err = fmt.Errorf("parsing HTML: %s", err)
		return
	}

	words, images = countWordsAndImages(doc)
	return
}

func countWordsAndImages(node *html.Node) (words, images int) {
	switch node.Type {
	case html.ElementNode:
		if node.Data == "img" {
			images++
		}
	case html.TextNode:
		words += len(strings.Fields(node.Data))
	}

	for c := node.FirstChild; c != nil; c = c.NextSibling {
		_words, _images := countWordsAndImages(c)
		words += _words
		images += _images
	}

	return
}

func main() {
	fmt.Printf("URL%37s\tWords\tImages\n", "")

	for _, url := range os.Args[1:] {
		words, images, err := CountWordsAndImages(url)
		if err != nil {
			log.Printf("%q error: %s", url, err)
			continue
		}

		fmt.Printf("%-40q\t%-5d\t%-5d\n", url, words, images)
	}
}
