package database

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/couchbase/gocb/v2"
)

func Connect() (*gocb.Cluster, error) {
	address := "couchbase://" + config.Main.Database.Address
	return gocb.Connect(
		address,
		gocb.ClusterOptions{
			Username: config.Main.Database.Username,
			Password: config.Main.Database.Password,
		},
	)
}