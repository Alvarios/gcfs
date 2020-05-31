package api

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/methods"
	"github.com/gorilla/mux"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	// fileId is sent via url parameter.
	fileId := mux.Vars(r)["id"]

	// Don't allow empty ids, no-op.
	if fileId == "" {
		http.Error(w, "No file id was provided in url.", http.StatusBadRequest)
		return
	}

	file, err := methods.Get(fileId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	b, err := json.Marshal(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
