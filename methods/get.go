package methods

import (
	"gcfs"
	"gcfs/config"
	"github.com/couchbase/gocb/v2"
	"time"
)

func Get(fileId string) (*interface{}, error) {
	data, getErr := gcfs.Cluster.Bucket(config.Main.Database.BucketName).DefaultCollection().Get(
		fileId,
		&gocb.GetOptions{
			Timeout: 10 * time.Second,
		},
	)

	if getErr != nil {
		return nil, getErr
	}

	var file interface{}
	parseErr := data.Content(&file)

	if parseErr != nil {
		return nil, parseErr
	}

	return &file, nil
}
