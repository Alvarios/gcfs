package methods

import (
	"fmt"
	"github.com/Alvarios/gcfs/config/errors"
	"github.com/Alvarios/gcfs/database"
	"github.com/couchbase/gocb/v2"
	"net/http"
	"time"
)

type UpdateSpec struct {
	Remove []string `json:"remove"`
	Upsert map[string]interface{} `json:"upsert"`
	Append map[string]interface{} `json:"append"`
}

func Update(id string, params UpdateSpec) (uint64, *errors.Error) {
	var specs []gocb.MutateInSpec

	if params.Remove != nil {
		for _, value := range params.Remove {
			// Cannot remove critical data.
			if value == "general" ||
				value == "general.name" ||
				value == "general.size" ||
				value == "general.creation_time" ||
				value == "general.modification_time" ||
				value == "general.format" ||
				value == "url" {
				return 0, &errors.Error{
					Code: http.StatusNotAcceptable,
					Message: fmt.Sprintf("you are trying to delete %s, which is a critical metadata", value),
				}
			}

			specs = append(specs, gocb.RemoveSpec(value, &gocb.RemoveSpecOptions{}))
		}
	}

	if params.Upsert != nil {
		upsertSpecs, err := flattenUpsertKeys(params.Upsert, "")
		if err != (*errors.Error)(nil) {
			return 0, err
		}

		err = checkUpsertKeys(upsertSpecs)
		if err != (*errors.Error)(nil) {
			return 0, err
		}

		for _, value := range upsertSpecs {
			key, ok := value[0].(string)
			if ok == false {
				return 0, &errors.Error{
					Code: http.StatusInternalServerError,
					Message: "unable to parse key",
				}
			}

			specs = append(specs, gocb.UpsertSpec(key, value[1], &gocb.UpsertSpecOptions{}))
		}
	}

	if params.Append != nil {
		for key, value := range params.Append {
			specs = append(specs, gocb.ArrayAppendSpec(key, value, &gocb.ArrayAppendSpecOptions{
				HasMultiple: true,
			}))
		}
	}

	_, err := database.Bucket.
		DefaultCollection().
		MutateIn(id, specs, &gocb.MutateInOptions{})

	if err != nil {
		return 0, &errors.Error{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return uint64(time.Now().UnixNano() / int64(time.Millisecond)), nil
}
