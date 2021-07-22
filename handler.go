package urlshort

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	urldb "github.com/kklee998/urlshort/db"
)

// SQLHandler returns a http.HandlerFunc that will attempt to find the
// URL based on the provided path. Path should not include the leading slash.
func SQLHandler(db *urldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		path, ok := vars["path"]
		if !ok {
			pathNotFound(w)
			return
		}
		exists, err := db.FindURLbyPath(path)
		if err != nil {
			log.Fatal(err)
		}
		if exists != nil {
			http.Redirect(w, r, exists.URL, http.StatusFound)
			return
		}
		pathNotFound(w)
	}
}

func pathNotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintln(w, "THE PATH YOU ARE LOOKING FOR IS IN ANOTHER CASTLE")
}

func InsertUpdateHandler(db *urldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPut:
			var input urldb.URLPath
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				log.Printf("Unable to parse JSON due to error: %s", err.Error())
				errMsg := map[string]string{
					"errors": "Unable to parse JSON",
				}
				writeJSONResponse(w, errMsg, http.StatusUnprocessableEntity)
				return
			}

			err = db.UpdateUrlAndPath(input)

			if err != nil {
				log.Printf("URLHandler %s: Unable to insert due to: %s", http.MethodPut, err.Error())
				errMsg := map[string]string{
					"errors": "Unable to insert into DB",
				}
				writeJSONResponse(w, errMsg, http.StatusBadRequest)
				return
			}

			var success = map[string]string{
				"message": "succesfully updated",
			}
			writeJSONResponse(w, success, http.StatusCreated)

		case http.MethodPost:
			var input urldb.URLPath
			err := json.NewDecoder(r.Body).Decode(&input)
			if err != nil {
				log.Printf("Unable to parse JSON due to error: %s", err.Error())
				errMsg := map[string]string{
					"errors": "Unable to parse JSON",
				}
				writeJSONResponse(w, errMsg, http.StatusUnprocessableEntity)
				return
			}

			err = db.SaveUrlAndPath(input)
			if err != nil {
				log.Printf("URLHandler %s: Unable to insert due to: %s", http.MethodPost, err.Error())
				errMsg := map[string]string{
					"errors": "Unable to insert into DB",
				}
				writeJSONResponse(w, errMsg, http.StatusBadRequest)
				return
			}

			var success = map[string]string{
				"message": "succesfully created",
			}
			writeJSONResponse(w, success, http.StatusCreated)

		default:
			http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
		}

	}

}

func DeleteHandler(db *urldb.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			vars := mux.Vars(r)
			path, ok := vars["path"]
			if !ok {
				pathNotFound(w)
			}
			err := db.DeleteURLbyPath(path)
			if err != nil {
				log.Printf("DeleteHandler %s: Unable to insert due to: %s", http.MethodDelete, err.Error())
				errMsg := map[string]string{
					"errors": "Unable to remove path from DB",
				}
				writeJSONResponse(w, errMsg, http.StatusBadRequest)
				return
			}
			emptyMsg := map[string]string{"": ""}
			writeJSONResponse(w, emptyMsg, http.StatusNoContent)

		default:
			http.Error(w, "METHOD NOT ALLOWED", http.StatusMethodNotAllowed)
		}

	}
}

func writeJSONResponse(w http.ResponseWriter, msg map[string]string, statusCode int) {
	err := json.NewEncoder(w).Encode(msg)
	if err != nil {
		log.Printf("Unable to write JSON response due to error: %s", err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(map[string]string{
			"errors": "Unable to send response",
		})
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

}
