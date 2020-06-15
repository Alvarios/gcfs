package gcfs

import (
	"fmt"
	"github.com/Alvarios/gcfs/api"
	"github.com/Alvarios/gcfs/api/middlewares"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/gcfs/database/metadata"
	"github.com/Alvarios/gcfs/methods"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"log"
	"net/http"
	"time"
)

type Configuration = config.Configuration
type DbConfig = config.Database
type GlobalConfig = config.Global
type ServerConfig = config.Server

type GeneralMetadata = metadata.GeneralData

type Error = responses.Error

type UpdateSpecs = methods.UpdateSpec
type InsertFlags = methods.InsertFlags

var Insert = methods.Insert
var InsertF = methods.InsertF
var Get = methods.Get
var Update = methods.Update
var Delete = methods.Delete

var AutoProvide = metadata.AutoProvide
var CheckIntegrity = metadata.CheckIntegrity

func Setup(c Configuration) {
	// Avoid unused var error.
	_, _, _, _, _, _, _ = Insert, InsertF, Get, Update, Delete, AutoProvide, CheckIntegrity

	config.LoadConfig(c)
	err := database.Connect()

	if err != nil {
		log.Fatal(err.Error())
	}

	// Launch server.
	if config.Main.Server.Provided() {
		fmt.Println("Starting server...")

		if config.Main.Routes.Provided() {
			mid := alice.New(middlewares.Log, middlewares.Recover)
			r := mux.NewRouter()

			if config.Main.Routes.Ping != "" {
				r.Handle(config.Main.Routes.Ping, mid.ThenFunc(api.Ping)).Methods("GET")
			}

			if config.Main.Routes.PingDatabase != "" {
				r.Handle(config.Main.Routes.PingDatabase, mid.ThenFunc(api.PingCouchbase)).Methods("GET")
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
/*
			if config.Main.Routes.Search != "" {
				r.Handle(config.Main.Routes.Search, mid.ThenFunc(api.Search)).Methods("POST")
			}*/

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