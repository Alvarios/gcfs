package database

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/couchbase/gocb/v2"
	"time"
)

var Cluster *gocb.Cluster
var Bucket *gocb.Bucket

func Connect() error {
	if config.Main.Database.Bucket == (*gocb.Bucket)(nil) {
		var err error
		Cluster, err = gocb.Connect(
			config.Main.Database.Address,
			gocb.ClusterOptions{
				Username: config.Main.Database.Username,
				Password: config.Main.Database.Password,
			},
		)

		if err != nil {
			return err
		}

		Bucket = Cluster.Bucket(config.Main.Database.BucketName)
		err = Bucket.WaitUntilReady(10 * time.Second, nil)

		return err
	} else {
		Bucket = config.Main.Database.Bucket
		return nil
	}
}