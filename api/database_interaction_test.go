package api

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"testing"
)

func TestDatabaseInteraction(t *testing.T) {
	config.LoadConfig(config.Configuration{})
	err := database.Connect()
	if err != (*responses.Error)(nil) {
		t.Fatalf("unable to connect to database : %s", err.Error())
	}

	/*
	mid := alice.New(middlewares.Log, middlewares.Recover)
	r := mux.NewRouter()*/
}
