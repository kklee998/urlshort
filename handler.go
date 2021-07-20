package urlshort

import (
	"encoding/json"
	"log"
	"net/http"

	urldb "github.com/kklee998/urlshort/db"
)

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

func URLHandler(db *urldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var input urldb.URL
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				http.Error(w, "Unable to parse JSON", http.StatusUnprocessableEntity)
				return
			}
			err = db.UpdateUrlAndPath(input)
			if err != nil {
				log.Printf("URLHandler %s: Unable to insert due to: %s", http.MethodPut, err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"errors": "Unable to insert into DB",
				})
				return
			}

			var success = map[string]string{
				"message": "succesfully updated",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(success)

		case http.MethodPost:
			var input urldb.URL
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				http.Error(w, "Unable to parse JSON", http.StatusUnprocessableEntity)
				return
			}
			err = db.SaveUrlAndPath(input)
			if err != nil {
				log.Printf("URLHandler %s: Unable to insert due to: %s", http.MethodPost, err.Error())
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"errors": "Unable to insert into DB",
				})
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
