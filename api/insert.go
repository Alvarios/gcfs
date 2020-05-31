package api

import (
	"encoding/json"
	"gcfs/api/responses"
	"gcfs/methods"
	"net/http"
)

func Insert(w http.ResponseWriter, r *http.Request) {
	// Load file data from body.
	var file interface{}
	err := json.NewDecoder(r.Body).Decode(&file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Here the magic happens :)
	fileId, ierr := methods.Insert(file, true)

	if ierr != nil {
		http.Error(w, ierr.Message, ierr.Code)
		return
	}

	// Insertion was successful, ready to inform user.
	w.WriteHeader(http.StatusOK)
	response := responses.Responses{
		Code:   http.StatusOK,
		Status: http.StatusText(http.StatusOK),
		Body:   fileId,
	}

	_ = json.NewEncoder(w).Encode(response)
	return
}
