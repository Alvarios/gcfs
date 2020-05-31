package api

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/api/responses"
	"github.com/Alvarios/gcfs/methods"
	"github.com/gorilla/mux"
	"net/http"
)

func Update(w http.ResponseWriter, r *http.Request) {
	// fileId is sent via url parameter.
	fileId := mux.Vars(r)["id"]

	// Don't allow empty ids, no-op.
	if fileId == "" {
		http.Error(w, "No file id was provided in url.", http.StatusBadRequest)
		return
	}

	var file methods.UpdateSpec
	err := json.NewDecoder(r.Body).Decode(&file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	timestamp, uerr := methods.Update(fileId, file)

	if uerr != nil {
		http.Error(w, uerr.Message, uerr.Code)
		return
	}

	w.WriteHeader(http.StatusOK)
	response := responses.Responses{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Body:   string(timestamp),
	}

	_ = json.NewEncoder(w).Encode(response)
	return
}
