package methods

import (
	"fmt"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/gcfs/database/metadata"
	mapUtils "github.com/Alvarios/kushuh-go-utils/map-utils"
	"github.com/Alvarios/kushuh-go-utils/number-utils"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"net/http"
	"regexp"
	"strconv"
)

type InsertFlags struct {
	AutoProvide bool `json:"auto_provide"`
	Force bool `json:"force"`
	Strict bool `json:"strict"`
}

func Insert(v interface{}, id string) (string, *responses.Error) {
	return InsertF(v, id, InsertFlags{
		AutoProvide: config.Main.Global.AutoProvide,
		Strict: config.Main.Global.Strict,
	})
}

func InsertF(v interface{}, id string, f InsertFlags) (string, *responses.Error) {
	timestamp := numberUtils.Timestamp()

	autoProvided, err := v, &responses.Error{}
	if f.AutoProvide {
		autoProvided, err = metadata.AutoProvide(v)
	}

	if err != (*responses.Error)(nil) {
		return "", err
	}

	if !f.Force {
		_, err = metadata.CheckIntegrity(autoProvided)
		if err != (*responses.Error)(nil) {
			return "", err
		}
	}

	if f.Strict && !f.Force {
		mv, err := mapUtils.ToMap(v)
		if err != nil {
			return "", err
		}

		fv := mapUtils.Flatten(mv, "")
		fe := mapUtils.Flatten(config.Main.Metadata, "")

		for k, _ := range fv {
			match := false
			for ek, _ := range fe {
				if k == ek {
					match = true
					break
				}
			}

			if match == false {
				return "", &responses.Error{
					Code:    400,
					Message: fmt.Sprintf("found non authorized key %s in strict mode", k),
				}
			}
		}
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
		return "", &responses.Error{
			Code:    http.StatusInternalServerError,
			Message: cerr.Error(),
		}
	}

	return fileId, nil
}
