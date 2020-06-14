package api

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/methods"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"github.com/gorilla/mux"
	"net/http"
)

func Get(w http.ResponseWriter, r *http.Request) {
	// fileId is sent via url parameter.
	fileId := mux.Vars(r)["id"]

	// Don't allow empty ids, no-op.
	if fileId == "" {
		http.Error(w, "no file id was provided in url", http.StatusBadRequest)
		return
	}

	file, lerr := methods.Get(fileId)

	if lerr != (*responses.Error)(nil) {
		http.Error(w, lerr.Error(), http.StatusInternalServerError)
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
