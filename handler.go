package urlshort

import (
	"log"
	"net/http"

	urldb "github.com/kklee998/urlshort/db"
)

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if dest, ok := pathsToUrls[path]; ok {
			http.Redirect(w, r, dest, http.StatusFound)
			return
		}
		fallback.ServeHTTP(w, r)
	}
}

// SQLHandler returns a http.HandlerFunc that will attempt to find the
// URL based on the provided path
func SQLHandler(db *urldb.DB, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		exists, err := db.FindURLbyPath(path)
		if err != nil {
			log.Fatal(err)
		}
		if exists != nil {
			http.Redirect(w, r, exists.URL, http.StatusFound)
			return
		}

		fallback.ServeHTTP(w, r)
	}
}
