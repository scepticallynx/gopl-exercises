/*Exercise 5.13: Modify crawl to make local copies of the pages it finds, creating
directories as necessary. Don't make copies of pages that come from different domain.
*/
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlp "net/url"
	"os"
	"path/filepath"

	"gopl.io/ch5/links"
)

// breadthFirst calls f for each item in the worklist.
// Any items returned by f are added to the worklist.
// f is called at most once for each item.
func breadthFirst(f func(item string) []string, worklist []string) {
	seen := make(map[string]bool)
	for len(worklist) > 0 {
		items := worklist
		worklist = nil
		for _, item := range items {
			if !seen[item] {
				seen[item] = true
				worklist = append(worklist, f(item)...)
			}
		}
	}
}

func saveOffline(url string) error {
	u, err := urlp.Parse(url)
	if err != nil {
		return fmt.Errorf("parsing url %q: %v", url, err)
	}

	// get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// if baseURL hasn't been stored yet, assign base url to be the first url
	// received (passed as argument)
	if baseURL == "" {
		baseURL = u.Hostname()

		d, err := os.Stat(baseURL)
		if err != nil || !d.IsDir() {
			log.Printf("Root directory %q does not exists, creating.", baseURL)
			err = os.Mkdir(baseURL, os.ModePerm)
			if err != nil {
				return fmt.Errorf("creating root directory %q in base directory %q: %v", baseURL, cwd, err)
			}
		}
	} else if u.Host != baseURL {
		log.Printf("%q is not the same domain as base url %q, ignoring", u, baseURL)
		return nil
	}

	// separate path and endpoint (dirs and file)
	dir, file := filepath.Split(u.Path)

	// save main page
	if file == "" {
		file = "index.html"
	}

	// join current working directory with base dir
	dir = filepath.Join(cwd, dir)

	// create directory structure in current working directory
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("creating directory %s: %v", dir, err)
	}

	// download url source
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filepath.Join(dir, file+".html"), body, os.ModePerm)
}

func crawl(url string) []string {
	fmt.Println(url)
	err := saveOffline(url)
	if err != nil {
		log.Printf("saving page %s: %v", url, err)
	}

	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

var baseURL string

func main() {
	// Crawl the web breadth-first,
	// starting from the command-line arguments.
	breadthFirst(crawl, os.Args[1:])
}
