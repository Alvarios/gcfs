package methods

import (
	"gcfs"
	"gcfs/config"
	"gcfs/config/data"
	"gcfs/database/metadata"
	"github.com/couchbase/gocb/v2"
	"net/http"
	"strconv"
	"time"
)

func Insert(v interface{}, flag ...bool) (string, *data.Error) {
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	autoProvided, err := v, &data.Error{}
	if config.Main.Global.AutoProvide || (len(flag) > 0 && flag[0]) {
		autoProvided, err = metadata.AutoProvide(v)
	}

	if err != nil {
		return "", err
	}

	m, err := metadata.CheckIntegrity(v)
	if err != nil {
		return "", err
	}

	if m.Id == "" {
		m.Id = strconv.FormatInt(timestamp, 10)
	}

	fileId := m.Id

	_, cerr := gcfs.Cluster.Bucket(config.Main.Database.BucketName).DefaultCollection().Upsert(
		fileId,
		autoProvided,
		&gocb.UpsertOptions{
			Timeout: 5 * time.Second,
		},
	)

	if cerr != nil {
		return "", &data.Error{
			Code:    http.StatusInternalServerError,
			Message: cerr.Error(),
		}
	}

	return fileId, nil
}
