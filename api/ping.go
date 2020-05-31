package api

import (
	"gcfs/api/responses"
	"net/http"
)

func Ping(w http.ResponseWriter, _ *http.Request) {
	// write header first to avoid flush error
	w.WriteHeader(http.StatusOK)
	// write response
	_, err := w.Write([]byte("pong"))

	// couldn't write the response
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
