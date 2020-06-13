package middlewares

import (
	"log"
	"net/http"
	"time"
)

//Log : Middleware used to display the API log
// format [Method] - URL DATETIME
func Log(next http.Handler) http.Handler {
	f := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		end := time.Now()
		log.Printf("[%s] - %s %s \n", r.Method, r.URL.String(), end.Sub(start))
	}

	return http.HandlerFunc(f)
}
