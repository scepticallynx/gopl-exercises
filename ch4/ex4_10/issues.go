// Exercise 4.10: Modify issues to report the result in age categories, say less
// than a month old, less than a year old, and more than a year old.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// github
const IssuesURL = "https://api.github.com/search/issues"

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

// SearchIssues queries the GitHub issue tracker.
func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	q := url.QueryEscape(strings.Join(terms, " "))
	resp, err := http.Get(IssuesURL + "?q=" + q)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("search query failed: %s", resp.Status)
	}

	var result IssuesSearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		resp.Body.Close()
		return nil, err
	}
	resp.Body.Close()
	return &result, nil
}

// Issues prints a table of GitHub issues matching the search terms.
func main() {
	var ageCategories = make(map[string][]*Issue)

	var now = time.Now()

	result, err := SearchIssues(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	var key string

	for _, item := range result.Items {
		switch dh := now.Sub(item.CreatedAt).Hours() / 24; {
		case dh > 365:
			key = ">1 year"
		case dh <= 7:
			key = "<1 week"
		case dh <= 30:
			key = "<1 month"
		case dh <= 365:
			key = "<1 year"
		}

		if v, ok := ageCategories[key]; !ok {
			ageCategories[key] = []*Issue{item}
		} else {
			ageCategories[key] = append(v, item)
		}
	}

	for k := range ageCategories {
		fmt.Println(k, len(ageCategories[k]))
		for _, item := range ageCategories[k] {
			fmt.Printf("%s #%-5d %9.9s %.55s\n",
				item.CreatedAt, item.Number, item.User.Login, item.Title)
		}
	}
}
