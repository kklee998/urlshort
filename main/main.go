package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/kklee998/urlshort"
	urldb "github.com/kklee998/urlshort/db"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	r := setupRouter()

	r.NotFoundHandler = http.HandlerFunc(NotFound)

	log.Println("connecting to DB")
	db, err := urldb.Open("sqlite3", "./urldb.sqlite3")
	checkErr(err)
	defer db.Close()

	err = db.StartDB()
	checkErr(err)
	log.Println("Succesfully connected to DB")

	r.HandleFunc("/urls", urlshort.InsertUpdateHandler(db))
	r.HandleFunc("/urls/{path}", urlshort.DeleteHandler(db))
	r.HandleFunc("/{path}", urlshort.SQLHandler(db))

	srv := &http.Server{
		Handler:      r,
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Println("Starting the server on :8080")
	log.Fatal(srv.ListenAndServe())
}

func setupRouter() *mux.Router {
	mux := mux.NewRouter()
	mux.HandleFunc("/", index)

	return mux
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello There")
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "THE PAGE YOU ARE LOOKING FOR IS IN ANOTHER CASTLE")
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
