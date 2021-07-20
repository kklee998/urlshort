package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/kklee998/urlshort"
	urldb "github.com/kklee998/urlshort/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := setupMux()

	r.NotFoundHandler = http.HandlerFunc(notFound)

	log.Println("connecting to DB")
	db, err := urldb.Open("sqlite3", "./urldb.sqlite3")
	checkErr(err)
	defer db.Close()

	err = db.StartDB()
	checkErr(err)
	log.Println("Succesfully connected to DB")

	sqlHandler := urlshort.SQLHandler(db, r)
	r.HandleFunc("/urls", urlshort.URLHandler(db))

	log.Println("Starting the server on :8080")
	log.Fatal(http.ListenAndServe(":8080", sqlHandler))
}

func setupMux() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", index)

	return mux
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello There")
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "THE PATH YOU ARE LOOKING FOR IS IN ANOTHER CASTLE")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
