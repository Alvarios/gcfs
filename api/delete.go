package api

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/api/responses"
	"github.com/Alvarios/gcfs/methods"
	"github.com/gorilla/mux"
	"net/http"
)

func Delete(w http.ResponseWriter, r *http.Request) {
	// fileId is sent via url parameter.
	fileId := mux.Vars(r)["id"]

	// Don't allow empty ids, no-op.
	if fileId == "" {
		http.Error(w, "No file id was provided in url.", http.StatusBadRequest)
		return
	}

	deleteErr := methods.Delete(fileId)

	if deleteErr != nil {
		http.Error(w, deleteErr.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	response := responses.Responses{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Body:   "",
	}

	_ = json.NewEncoder(w).Encode(response)
	return
}
