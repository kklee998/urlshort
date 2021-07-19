package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/kklee998/urlshort"
	urldb "github.com/kklee998/urlshort/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	m := setupMux()

	fmt.Println("connecting to DB")
	db, err := urldb.Open("sqlite3", "./urldb.sqlite3")
	checkErr(err)
	defer db.Close()

	err = db.StartDB()
	checkErr(err)

	sqlHandler := urlshort.SQLHandler(db, m)

	fmt.Println("Starting the server on :8080")
	http.ListenAndServe(":8080", sqlHandler)
}

func setupMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", index)
	mux.HandleFunc("*", notFound)

	return mux
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello There")
}

// TODO: add 404
func notFound(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "THE PATH YOU ARE LOOKING FOR IS IN ANOTHER CASTLE")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
