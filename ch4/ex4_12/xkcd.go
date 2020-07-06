/* Exercise 4.12: A request to https://xkcd.com/571/info.0.json produces a detailed
 description of comic 571. Download each URL (once!) and build an offline index.
 Write a tool xkcd that, using this index, prints the URL and transcript of each
 comic that matches a search term provided on the command line.

performs search by 1+ specified terms:
go run xkcd.go term1 term2 ... termN

index comics only missing in offline storage, if any:
go run xkcd.go -reindex

index ALL comics from scratch:
rm xkcd_index.json && go run xkcd.go -reindex
*/

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"unicode"
)

const baseURL string = `https://xkcd.com` // root
const endpoint string = `info.0.json`     // default endpoint for all comics urls

// Comic represents a singe API response for requested comic
type Comic struct {
	Month      string          `json:"month"`
	Num        uint64          `json:"num"`
	Link       string          `json:"link"`
	Year       string          `json:"year"`
	News       string          `json:"news"`
	SafeTitle  string          `json:"safe_title"`
	Transcript string          `json:"transcript"`
	Alt        string          `json:"alt"`
	Img        string          `json:"img"`
	Title      string          `json:"title"`
	Day        string          `json:"day"`
	Tokens     map[string]bool `json:"tokens,omitempty"`
}

// ---- fetch ----

// getLatestComicID queries API without comic ID what returns the latest available comic on xkcd.
func getLatestComicID() uint64 {
	u, err := url.Parse(fmt.Sprintf("%s/%s", baseURL, endpoint))
	if err != nil {
		log.Printf("xkcd: error parsing url: %v", err)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Printf("xkcd: query error: %v", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("xkcd: HTTP Status: %d", resp.StatusCode)
	}

	comic := new(Comic)
	err = json.NewDecoder(resp.Body).Decode(comic)

	if err != nil {
		log.Printf("xkcd: unmarshal error: %v", err)
	}

	return comic.Num
}

// fetchComic queries API for specific comic by its id
func fetchComic(comicID uint64) (*Comic, error) {
	log.Printf("Fetching comic %d", comicID)
	u, err := url.Parse(fmt.Sprintf("%s/%d/%s", baseURL, comicID, endpoint))
	if err != nil {
		return nil, fmt.Errorf("xkcd: comic %d: error parsing url: %v", comicID, err)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("xkcd: comic %d: query error: %v", comicID, err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("xkcd: comic %d: HTTP Status: %d", comicID, resp.StatusCode)
	}

	comic := new(Comic)
	err = json.NewDecoder(resp.Body).Decode(comic)

	if err != nil {
		return nil, fmt.Errorf("xkcd: comic %d: unmarshal error: %v", comicID, err)
	}

	return comic, nil
}

// ----- tokenize ------

// TODO: separate attachable to tokenizeComic tokenizeFunction
// tokenizeComic populates Comic.Tokens with lowercase words taken from
// Comic.Title, Comic.Alt, Comic.Transcript. In-place.
func tokenizeComic(c *Comic) {
	tokens := make(map[string]bool, 0)

	scanner := bufio.NewScanner(strings.NewReader(c.Title + c.Alt + c.Transcript))
	scanner.Split(bufio.ScanWords)

	for scanner.Scan() {
		// to lowercase
		_text := strings.ToLower(scanner.Text())
		if len(_text) < 3 {
			continue
		}

		text := make([]rune, 0, len(_text))

		// cleanup: nothing but digits and letters
		for _, r := range _text {
			switch {
			case unicode.IsLetter(r), unicode.IsDigit(r), unicode.IsNumber(r):
				text = append(text, r)
			default:
				continue
			}
		}

		_text = string(text)
		if _, ok := tokens[_text]; ok {
			continue
		}

		tokens[_text] = true
	}

	c.Tokens = tokens
}

// reTokenize reads all comics stored offline and performs tokenization (in case
// of tokenization function update).
// func reTokenize() {
//
// }

// -------- index related functions -------

// dumpIndex takes slice of Comic to dump it into json file
func dumpIndex(comics []*Comic) error {
	// without os.O_WRONLY "bad file descriptor" error occurs
	file, err := os.OpenFile("xkcd_index.json", os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Print("file error")
		return err
	}

	defer file.Close()

	var b bytes.Buffer
	encoder := json.NewEncoder(&b)
	encoder.SetIndent("", "\t")

	err = encoder.Encode(comics)
	if err != nil {
		log.Print("xkcd: json error while encoding comics slice")
		return err
	}

	_, err = file.Write(b.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// unmarshalIndex unpacks list of comics stored in default index file, if any exist.
func unmarshalIndex() ([]*Comic, error) {
	f, err := ioutil.ReadFile("xkcd_index.json")
	if err != nil {
		return nil, err
	}

	results := make([]*Comic, 0)

	err = json.Unmarshal(f, &results)
	if err != nil {
		return nil, err
	}

	return results, nil
}

// reindexComics queries xkcd.com to obtain missing in index comics. To perform full
// reindex, delete index file (xkcd_index.json).
func reindexComics() error {
	latestAvailable := getLatestComicID()
	log.Printf("Latest comic available: %d\n", latestAvailable)
	latestIndexed := uint64(1)

	indexedComics, err := unmarshalIndex()

	// TODO: separate file not found error = reindex everything and read file error
	// if index already exists, read the latest indexed comic
	// here we assume that even comics are unordered, the index does not contain
	// "holes" (skipped comics)
	if err == nil {
		for _, comic := range indexedComics {
			if comic.Num > latestIndexed {
				latestIndexed = comic.Num
			}
		}
	} else {
		log.Printf("%v", err)
	}

	if latestIndexed >= latestAvailable {
		log.Printf("xkcd: all available comics are already indexed\n")
		return nil
	}

	log.Printf("Fetching missing comics:\n")

	for i := latestIndexed; i <= latestAvailable; i++ {
		// query comic
		c, err := fetchComic(i)
		if err != nil {
			log.Printf("xkcd: fetching comic: %v", err)
			continue
		}

		// tokenize successfully parsed comic
		tokenizeComic(c)

		indexedComics = append(indexedComics, c)
	}

	if err := dumpIndex(indexedComics); err != nil {
		return fmt.Errorf("xkcd: error dumping comics index to json: %v", err)
	}

	return nil
}

// ------ search related functions -------

// search looks for each given term if it's not an empty string in tokens of
// every comic stored.
func search(terms ...string) {
	// read index from file
	comics, err := unmarshalIndex()
	if err != nil {
		log.Fatal(err)
	}

	for _, term := range terms {
		if term == "" {
			continue
		}

		// TODO: map[token][]ComicID ? to search faster
		fmt.Printf("Search term: %s\n%s\n", term, strings.Repeat("-", 20))
		for _, comic := range comics {
			if _, ok := comic.Tokens[term]; ok {
				fmt.Printf("Comic ID: %d\nURL: %s\nTranscript: %q\n%s\n", comic.Num, comic.Img, comic.Transcript, strings.Repeat("-", 40))
			}
		}
	}
}

// --------- main ----------

const usage string = `build entire index from scratch: go run xkcd -reindex
search specified "term": go run xkcd search_term`

func main() {
	var reindex bool
	flag.BoolVar(&reindex, "reindex", false, "index all comics from scratch")
	flag.Parse()

	args := flag.Args()
	if !reindex && len(args) == 0 {
		fmt.Println(usage)
		os.Exit(1)
	}

	if reindex {
		err := reindexComics()
		if err != nil {
			log.Fatalf("xkcd: %v", err)
		}

		os.Exit(0)
	}

	search(args...)
}
