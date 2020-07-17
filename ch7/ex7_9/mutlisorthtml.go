/*Exercise 7.9: Use html/template package (§ 4.6) to replace printTracks with a function
that displays the tracks as an HTML table. Use the solution to the previous exercise
to arrange that each click on a column head makes an HTTP request to sort the table.*/
package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"sort"
	"time"
)

type Track struct {
	Title  string
	Artist string
	Album  string
	Year   int
	Length time.Duration
}

func length(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		panic(s)
	}
	return d
}

type customSort struct {
	t    []*Track
	less func(x, y *Track) bool
}

func (x customSort) Len() int           { return len(x.t) }
func (x customSort) Less(i, j int) bool { return x.less(x.t[i], x.t[j]) }
func (x customSort) Swap(i, j int)      { x.t[i], x.t[j] = x.t[j], x.t[i] }

func MultiTierSort(tracks []*Track, by []string) sort.Interface {
	return customSort{
		tracks,
		func(x, y *Track) bool {
			for i := 0; i < len(by); i++ {
				if by[i] == "Title" && x.Title != y.Title {
					return x.Title < y.Title
				}

				if by[i] == "Year" && x.Year != y.Year {
					return x.Year < y.Year
				}

				if by[i] == "Artist" && x.Artist != y.Artist {
					return x.Artist < y.Artist
				}

				if by[i] == "Album" && x.Album != y.Album {
					return x.Album < y.Album
				}

				if by[i] == "Length" && x.Length != y.Length {
					return x.Length < y.Length
				}
			}

			return true
		},
	}
}

// simple minimal HTML to show the table
const templ = `
<html>
<body>
<table>
<tr style='text-align: left'>
	<th><a href='?sort="Artist"'>Artist</a></th>
	<th><a href='?sort="Title"'>Title</a></th>
	<th><a href='?sort="Album"'>Album</a></th>
	<th><a href='?sort="Year"'>Year</a></th>
	<th><a href='?sort="Length"'>Length</a></th>
</tr>
{{range .}}
<tr>
	<td>{{.Artist}}</td>
	<td>{{.Title}}</td>
	<td>{{.Album}}</td>
	<td>{{.Year}}</td>
	<td>{{.Length}}</td>
</tr>
{{end}}
</table>
</body>
</html>
`

func htmlTracks(tracks []*Track, by string, w io.Writer) {
	sort.Sort(MultiTierSort(tracks, []string{by}))

	t := template.Must(template.New("sort").Parse(templ))
	if err := t.Execute(w, tracks); err != nil {
		log.Fatal(err)
	}
}

var tracks = []*Track{
	{"Go", "Delilah", "From the Roots Up", 2012, length("3m38s")},
	{"Go", "Moby", "Moby", 1992, length("3m37s")},
	{"Go Ahead", "Alicia Keys", "As I Am", 2007, length("4m36s")},
	{"Ready 2 Go", "Martin Solveig", "Smash", 2011, length("4m24s")},
}

func main() {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := r.ParseForm(); err != nil {
			log.Println(err)
		}

		by := ""
		if f := r.Form["sort"]; len(f) > 0 {
			by = f[0]
		}

		htmlTracks(tracks, by, w)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
