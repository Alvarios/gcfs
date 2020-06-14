package api

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/methods"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
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
	fileId, ierr := methods.InsertF(file, "", true)

	if ierr != (*responses.Error)(nil) {
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
