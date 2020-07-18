/* Exercise 7.11: Add additional handlers so that clients can create, read, update,
and delete database entries. For example, a request of the form /update?item=socks&price=6
will update the price of an item in the inventory and report an error if the
item does not exist or if the price is invalid. (Warning: this change introduces
concurrent variable updates.)

Exercise 7.12: Change the handler for /list to print its output as an HTML table, not text.

Http4 is an e-commerce server that registers the /list and /price
endpoint by calling http.HandleFunc.
*/
package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

func main() {
	db := database{"shoes": 50, "socks": 5}
	http.HandleFunc("/list", db.list)
	http.HandleFunc("/price", db.price)

	http.HandleFunc("/create", db.create)
	http.HandleFunc("/read", db.read)
	http.HandleFunc("/update", db.update)
	http.HandleFunc("/delete", db.delete)

	log.Fatal(http.ListenAndServe("localhost:8000", nil))
}

type dollars float32

func (d dollars) String() string { return fmt.Sprintf("$%.2f", d) }

const itemsTemplate string = `
<html>
<head>
	<title>Inventory:</title>
</head>
<body>
	<h1>Inventory:</h1>
	<table>
	<tr style='text-align: center'>
		<th>Item</th>
		<th>Price</th>
	</tr>
{{range $item, $price := .}}
	<tr>
		<td>{{$item}}</td>
		<td>{{$price}}</td>
	</tr>
{{end}}
	</table>
</body>
</html>
`

var t = template.Must(template.New("items").Parse(itemsTemplate))

type database map[string]dollars

func (db database) list(w http.ResponseWriter, req *http.Request) {
	if err := t.Execute(w, db); err != nil {
		log.Fatal(err)
	}
}

func (db database) price(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%s\n", price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}

func parsePrice(price string) (float64, error) {
	if price == "" {
		return 0, fmt.Errorf("no price specified")
	}

	p, err := strconv.ParseFloat(price, 32)
	if err != nil {
		return 0, err
	}

	if p < 0 {
		return 0, fmt.Errorf("price is negative")
	}

	return p, nil
}

func (db database) create(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	p, err := parsePrice(price)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, " for: %q\n", item)
		return
	}

	if _, ok := db[item]; ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item already exists in database: %q\n", item)
	} else {
		w.WriteHeader(http.StatusOK)
		db[item] = dollars(p)
		fmt.Fprintf(w, "item %q with price %.2f created\n", item, p)
	}
}

func (db database) update(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	price := req.URL.Query().Get("price")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item does not exist in database: %q\n", item)
	} else {
		p, err := parsePrice(price)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, " for: %q\n", item)
			return
		}

		db[item] = dollars(p)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "item %q modified: new price %.2f\n", item, p)
	}
}

func (db database) delete(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")

	if _, ok := db[item]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "item does not exist in database: %q\n", item)
	} else {
		delete(db, item)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "item %q deleted\n", item)
	}
}

func (db database) read(w http.ResponseWriter, req *http.Request) {
	item := req.URL.Query().Get("item")
	if price, ok := db[item]; ok {
		fmt.Fprintf(w, "%q: %s\n", item, price)
	} else {
		w.WriteHeader(http.StatusNotFound) // 404
		fmt.Fprintf(w, "no such item: %q\n", item)
	}
}
