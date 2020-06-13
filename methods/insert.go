package methods

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/config/errors"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/gcfs/database/metadata"
	"github.com/Alvarios/kushuh-go-utils/number-utils"
	"net/http"
	"regexp"
	"strconv"
)

func Insert(v interface{}, id string) (string, *errors.Error) {
	return InsertF(v, id, config.Main.Global.AutoProvide)
}

func InsertF(v interface{}, id string, forceAutoProvide bool) (string, *errors.Error) {
	timestamp := numberUtils.Timestamp()

	autoProvided, err := v, &errors.Error{}
	if config.Main.Global.AutoProvide || forceAutoProvide {
		autoProvided, err = metadata.AutoProvide(v)
	}

	if err != (*errors.Error)(nil) {
		return "", err
	}

	_, err = metadata.CheckIntegrity(autoProvided)
	if err != (*errors.Error)(nil) {
		return "", err
	}

	fileId := id
	if id == "" {
		fileId = strconv.FormatUint(timestamp, 10)
	} else {
		reTs := regexp.MustCompile(`\{ts}`)
		fileId = reTs.ReplaceAllString(fileId, strconv.FormatUint(timestamp, 10))
	}

	_, cerr := database.Bucket.DefaultCollection().Upsert(
		fileId,
		autoProvided,
		nil,
	)

	if cerr != nil {
		return "", &errors.Error{
			Code:    http.StatusInternalServerError,
			Message: cerr.Error(),
		}
	}

	return fileId, nil
}
