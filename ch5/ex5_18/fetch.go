/* Exercise 5.18: Without changing its behaviour, rewrite fetch function to use
defer to close the writable file.
*/
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

// Fetch downloads the URL and returns the
// name and length of the local file.
func fetch(url string) (filename string, n int64, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	local := path.Base(resp.Request.URL.Path)
	if local == "/" {
		local = "index.html"
	}
	f, err := os.Create(local)
	if err != nil {
		return "", 0, err
	}

	n, err = io.Copy(f, resp.Body)

	// 1) anonymous function has access to variables and named results of enclosing function.
	// 2) deferred functions run AFTER return statements have updated the function's result variables.
	// 3) deferred anonymous function can observe the enclosing function's results.
	defer func() {
		// Close file, but prefer error from Copy, if any.
		if closeErr := f.Close(); err == nil {
			err = closeErr
		}
	}()

	// currently, returned err is io.Copy err if any.
	// f.Close error is set only if io.Copy error is nil.
	// f.Close might also be nil.
	return local, n, err
}

func main() {
	for _, url := range os.Args[1:] {
		local, n, err := fetch(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch %s: %v\n", url, err)
			continue
		}
		fmt.Fprintf(os.Stderr, "%s => %s (%d bytes).\n", url, local, n)
	}
}
