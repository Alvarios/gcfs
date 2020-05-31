package gcfs

import (
	"fmt"
	"gcfs/api"
	"gcfs/api/middlewares"
	"gcfs/config"
	"gcfs/database"
	"github.com/couchbase/gocb/v2"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

var Cluster *gocb.Cluster

func init() {
	Cluster, err := database.Connect()

	// Launch server.
	if config.Main.Server.Provided() {
		fmt.Println("Starting server...")

		if err != nil {
			panic(err.Error())
		}

		if config.Main.Routes.Provided() {
			mid := alice.New(middlewares.Log, middlewares.Recover)
			r := mux.NewRouter()

			if config.Main.Routes.Ping != "" {
				r.Handle(config.Main.Routes.Ping, mid.ThenFunc(api.Ping)).Methods("GET")
			}

			if config.Main.Routes.PingDatabase != "" {
				r.Handle(config.Main.Routes.Ping, mid.ThenFunc(api.PingCouchbase)).Methods("GET")
			}

			if config.Main.Routes.Insert != "" {
				r.Handle(config.Main.Routes.Insert, mid.ThenFunc(api.Insert)).Methods("POST")
			}

			if config.Main.Routes.Delete != "" {
				r.Handle(config.Main.Routes.Delete, mid.ThenFunc(api.Delete)).Methods("DELETE")
			}

			if config.Main.Routes.Get != "" {
				r.Handle(config.Main.Routes.Get, mid.ThenFunc(api.Get)).Methods("GET")
			}

			if config.Main.Routes.Search != "" {
				r.Handle(config.Main.Routes.Search, mid.ThenFunc(api.Search)).Methods("POST")
			}

			if config.Main.Routes.Update != "" {
				r.Handle(config.Main.Routes.Update, mid.ThenFunc(api.Update)).Methods("POST")
			}

			// Link router to base route.
			http.Handle("/", r)

			// Start router.
			srv := &http.Server{
				Handler:      r,
				Addr:         config.Main.Server.Port,
				WriteTimeout: 15 * time.Second,
				ReadTimeout:  15 * time.Second,
			}

			log.Printf("Server started on port %s", config.Main.Server.Port)
			log.Fatal(srv.ListenAndServe())
		}
	}

	return
}