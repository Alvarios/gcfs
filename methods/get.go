package methods

import (
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"github.com/couchbase/gocb/v2"
	"time"
)

func Get(fileId string) (interface{}, *responses.Error) {
	data, getErr := database.Bucket.DefaultCollection().Get(
		fileId,
		&gocb.GetOptions{
			Timeout: 10 * time.Second,
		},
	)

	if getErr != nil {
		return nil, &responses.Error{
			Code: 500,
			Message: getErr.Error(),
		}
	}

	var file interface{}
	parseErr := data.Content(&file)

	if parseErr != nil {
		return nil, &responses.Error{
			Code: 500,
			Message: parseErr.Error(),
		}
	}

	return file, nil
}
