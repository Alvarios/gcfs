package methods

import (
	"github.com/Alvarios/gcfs/config/errors"
	"github.com/Alvarios/gcfs/database"
	"github.com/couchbase/gocb/v2"
	"time"
)

func Get(fileId string) (interface{}, *errors.Error) {
	data, getErr := database.Bucket.DefaultCollection().Get(
		fileId,
		&gocb.GetOptions{
			Timeout: 10 * time.Second,
		},
	)

	if getErr != nil {
		return nil, &errors.Error{
			Code: 500,
			Message: getErr.Error(),
		}
	}

	var file interface{}
	parseErr := data.Content(&file)

	if parseErr != nil {
		return nil, &errors.Error{
			Code: 500,
			Message: parseErr.Error(),
		}
	}

	return file, nil
}
