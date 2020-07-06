/* Exercise 4.13: OMBD API. ... Write a tool poster that downloads the poster
image for the movie named on the command line.

Poster API is only available to patrons. So we extract poster from default API response.

set api key environment variable from file. File format: OMBD_KEY=key
set -a && . ~/YOUR_API_KEY_FILE && set +a && go run poster.go movie title

Use "" for multi-word movie title or simply pass it as is,
all arguments are joined into single string.
*/
package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const baseurl string = "https://omdbapi.com/?apikey="
const byTitle string = "&t="
const byID string = "&i="
const search string = "&s="

// taken from environment variable in init() function, make sure it is set
var apikey string

// setAPIKeyFromEnv reads OMDB_KEY environment variable to obtain API key.
// Fails if variable is not set or set but empty.
func setAPIKeyFromEnv() {
	key, set := os.LookupEnv("OMDB_KEY")
	if !set {
		log.Fatalf("OMDB_KEY environment variable is not set")
	} else if key == "" {
		log.Fatalf("OMDB_KEY environment variable is empty")
	}

	apikey = key
}

// MovieShort represents a brief info about movie in search results
// as returned by querying with "s=" parameter
type MovieShort struct {
	Title     string // title of the movie
	Year      string // year the movie has been released
	ID        string `json:"imdbID"`
	Type      string // movie, series etc
	PosterURL string `json:"poster" ` // url to poster of the movie
}

// MovieSearchResult is a slice of movies matched search query
type MoviesSearchResult struct {
	Search []*MovieShort `json:"Search"` // uppercase
}

// searchMovie queries API with "s=title" parameter
func searchMovie(title string) (*MoviesSearchResult, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", baseurl+apikey, search, url.QueryEscape(title)))
	if err != nil {
		return nil, err
	}

	// hasn't been describe yet in current and previous chapters
	defer resp.Body.Close()

	result := new(MoviesSearchResult)

	err = json.NewDecoder(resp.Body).Decode(result)
	return result, err
}

// Movie is a full description of requested movie as returned by querying API
// with t= parameter. MovieShort is embedded
type Movie struct {
	MovieShort
	Rated    string // rating of the movie
	Released string // full date of release of the movie. Format: DD Mon YYYY
	Runtime  string // duration of the movie
	Genre    string
	Director string
	Writer   string
	Actors   string
	Plot     string
	Language string
	Country  string
	Awards   string
	Ratings  []struct {
		Source string // where the movie has been rated
		Value  string // score the movie has been given on specific Source
	} // array with sources
	Metascore  string
	imdbRating string
	imdbVotes  string
	DVD        string // date of release on DVD. Format: DD Mon YYYY
	BoxOffice  string // amount $0,123,456 - "," as separator
	Production string
	Website    string
	Response   string
}

// queryMovieByTitle unlike searchMovie, asks API to return specific single movie
func queryMovieByTitle(title string) (*Movie, error) {
	resp, err := http.Get(fmt.Sprintf("%s%s%s", baseurl+apikey, byTitle, url.QueryEscape(title)))
	if err != nil {
		return nil, err
	}

	// hasn't been describe yet in current and previous chapters
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	movie := new(Movie)
	err = json.Unmarshal(r, movie)

	if err != nil {
		return nil, fmt.Errorf("error unmarshaling response body %v", err)
	}

	return movie, nil
}

// downloadPosterImg retrieves poster of the queried movie and stores it to posterFilename file.
// Extension .png will be added if specified filename is without any extensions
// or other extension.
func downloadPosterImg(url, posterFilename string) {
	log.Printf("Downloading poster into: %s", posterFilename)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("poster image request error: %v", err)
	}

	defer resp.Body.Close()

	// check whether specified filename has extension .png, add if hasn't
	if ext := filepath.Ext(posterFilename); ext != ".png" {
		posterFilename = posterFilename + ext
	}

	// create new file
	file, err := os.Create(posterFilename)
	if err != nil {
		log.Fatalf("error creating file %s: %v", posterFilename, err)
	}

	// write image (response body) to created file
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatalf("error writing image data to file: %v", err)
	}
}

// Get API key from environment variable set on github or locally
func init() {
	setAPIKeyFromEnv()
}

func main() {
	if len(os.Args) <= 1 {
		log.Println("no movie title given")
		os.Exit(1)
	}

	// merge all arguments into single string
	title := strings.Join(os.Args[1:], " ")

	movie, err := queryMovieByTitle(title)
	if err != nil {
		log.Printf("error querying movie: %v", err)
		os.Exit(1)
	}

	// TODO: looks like IMDB returns first best-matching movie automatically.
	//  Querying real nonsense returns nothing from both functions. Remove this.
	if movie.Response == "False" { // it's simpler, strconv.ParseBool returns value + error
		log.Printf("requested movie %q does not exist in IMDB database", title)
		fuzzyMatch, err := searchMovie(title)
		if err != nil {
			log.Fatal(err)
		}

		if len(fuzzyMatch.Search) == 0 {
			log.Fatalf("neither direct query nor search by title gave any results for provided movie: %q", title)
		}

		fmt.Printf("Here's the list of movies matching %q somehow:\n", title)
		for _, m := range fuzzyMatch.Search {
			fmt.Printf("Title: %q\nYear: %s\nIMDB ID: %s\nType: %s\nPoster: %s\n",
				m.Title, m.Year, m.ID, m.Type, m.PosterURL)
		}

		os.Exit(0)
	}

	// as we only need a posted
	downloadPosterImg(movie.PosterURL, movie.Title+" "+movie.Year+"_poster.png")
}
