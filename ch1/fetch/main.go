/* Fetch prints the content found at each specified URL.

Exercise 1.7: Use io.Copy(dst, src) instead of ioutil.ReadAll() to copy
		the response body to os.Stdout without a buffer large enough to hold the
		entire stream.
		go run ch1/fetch/main.go http://google.com

Exercise 1.8: Modify fetch to add prefix http:// to each argument URL if it is missing.

Exercise 1.9: Modify fetch to also print the HTTP status code, found in resp.Status.
 */
package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func main() {
	const prefix string = "http://"

	for _, url := range os.Args[1:] {
		// Exercise 1.8 (see package doc)
		if !strings.HasPrefix(url, prefix) && !strings.HasPrefix(url, "https://") {
			url = prefix + url
		}

		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: %v\n", err)
			os.Exit(1)
		}
		// Exercise 1.7 (see package doc)
		written, err := io.Copy(os.Stdout, resp.Body)

		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(os.Stderr, "fetch: reading %s: %v\n", url, err)
			os.Exit(1)
		}

		// Exercise 1.8: (see package doc)
		fmt.Printf("\n%d bytes in response.\nStatus: %s\n", written, resp.Status)
	}
}

