package middlewares

import (
	"encoding/json"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"log"
	"net/http"
)

// Recover Handler : Middleware that handle panic and return error 500 message
func Recover(next http.Handler) http.Handler {

	fn := func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("panic : %+v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Header().Set("Content-Type", "json/application")
				panicError := &responses.Error{
					Code:    http.StatusInternalServerError,
					Message: "panic internal server error",
				}

				encodeError := json.NewEncoder(w).Encode(panicError)
				if encodeError != nil {
					// A recover error can be critical, though it is better to avoid shutting down a running server.
					log.Printf("\n\nfailed to recover: %s", encodeError.Error())
					return
				}

				return
			}
		}()
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
