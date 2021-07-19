package urlshort

import (
	"encoding/json"
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

func URLCreateHandler(db *urldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			var input urldb.URL
			var error = map[string]string{
				"errors": "Unable to insert into DB",
			}
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				http.Error(w, "Unable to parse JSON", http.StatusUnprocessableEntity)
				return
			}
			err = db.SaveUrlAndPath(input)
			if err != nil {
				log.Printf("URLCreateHandler: Unable to insert due to: %s", err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(error)
				return
			}

			var success = map[string]string{
				"message": "succesfully created",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(success)

		default:
			http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)

		}

	}

}
