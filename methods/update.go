package methods

import (
	"fmt"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	mapUtils "github.com/Alvarios/kushuh-go-utils/map-utils"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"github.com/couchbase/gocb/v2"
	"net/http"
	"time"
)

type UpdateSpec struct {
	Remove []string `json:"remove"`
	Upsert map[string]interface{} `json:"upsert"`
	Append map[string]interface{} `json:"append"`
	Force bool `json:"strict"`
}

func Update(id string, params UpdateSpec) (uint64, *responses.Error) {
	var specs []gocb.MutateInSpec

	if params.Remove != nil && config.Main.Global.Strict {
		return 0, &responses.Error{
			Code: http.StatusNotAcceptable,
			Message: "you are forbidden delete any key in strict mode",
		}
	}

	if params.Remove != nil {
		for _, value := range params.Remove {
			// Cannot remove critical data.
			if !params.Force && value == "general" ||
				value == "general.name" ||
				value == "general.size" ||
				value == "general.creation_time" ||
				value == "general.modification_time" ||
				value == "general.format" ||
				value == "url" {
				return 0, &responses.Error{
					Code: http.StatusNotAcceptable,
					Message: fmt.Sprintf("you are trying to delete %s, which is a critical metadata", value),
				}
			}

			specs = append(specs, gocb.RemoveSpec(value, &gocb.RemoveSpecOptions{}))
		}
	}

	if params.Upsert != nil {
		upsertSpecs := mapUtils.Flatten(params.Upsert, "")

		err := checkUpsertKeys(upsertSpecs)
		if err != (*responses.Error)(nil) {
			return 0, err
		}

		for key, value := range upsertSpecs {
			if config.Main.Global.Strict && !params.Force {
				match := false
				for k, _ := range config.Main.Metadata {
					if k == key {
						match = true
						break
					}
				}

				if match == false {
					return 0, &responses.Error{
						Code: http.StatusNotAcceptable,
						Message: fmt.Sprintf("you are trying to add key %s, which is not an allowed metadata", key),
					}
				}
			}

			specs = append(specs, gocb.UpsertSpec(key, value, &gocb.UpsertSpecOptions{}))
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
		return 0, &responses.Error{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return uint64(time.Now().UnixNano() / int64(time.Millisecond)), nil
}
