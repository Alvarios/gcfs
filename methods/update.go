package methods

import (
	"fmt"
	"gcfs"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/config/data"
	"github.com/couchbase/gocb/v2"
	"net/http"
	"time"
)

type UpdateSpec struct {
	Remove interface{} `json:"remove"`
	Upsert interface{} `json:"upsert"`
	Append interface{} `json:"append"`
}

func parseUpsertKeys(params map[string]interface{}, parentKey string) ([][]interface{}, *data.Error) {
	var output [][]interface{}

	for key, value := range params {
		fKey := fmt.Sprintf("%s%s", parentKey, key)

		if fKey == "general" {
			mValue, ok := value.(map[string]interface{})
			if ok == false {
				return nil, &data.Error{
					Code: http.StatusBadRequest,
					Message: fmt.Sprintf("Illegal value %v for general metadata.", value),
				}
			}

			subOutput, err := parseUpsertKeys(mValue, fmt.Sprintf("%s.", fKey))
			if err != nil {
				return nil, err
			}

			output = append(output, subOutput...)
		} else if fKey == "general.name" || fKey == "url" || fKey == "format" {
			strValue, ok := value.(string)
			if ok == false {
				return nil, &data.Error{
					Code: http.StatusBadRequest,
					Message: fmt.Sprintf("Illegal value %v for %s metadata.", value, fKey),
				}
			}

			output = append(output, []interface{}{fKey, strValue})
		} else if fKey == "general.creation_time" || fKey == "general.size" {
			uintValue, ok := value.(uint64)
			if ok == false {
				return nil, &data.Error{
					Code: http.StatusBadRequest,
					Message: fmt.Sprintf("Illegal value %v for %s metadata.", value, fKey),
				}
			}

			output = append(output, []interface{}{fKey, uintValue})
		} else if mValue, ok := value.(map[string]interface{}); ok {
			subOutput, err := parseUpsertKeys(mValue, fmt.Sprintf("%s.", fKey))
			if err != nil {
				return nil, err
			}

			output = append(output, subOutput...)
		} else {
			output = append(output, []interface{}{fKey, value})
		}
	}

	return output, nil
}

func Update(id string, params UpdateSpec) (uint64, *data.Error) {
	castRemove, okRemove := params.Remove.([]string)
	if okRemove == false {
		return 0, &data.Error{
			Code: http.StatusBadRequest,
			Message: "Remove is not a valid array of string.",
		}
	}

	castUpsert, okUpsert := params.Upsert.(map[string]interface{})
	if okUpsert == false {
		return 0, &data.Error{
			Code: http.StatusBadRequest,
			Message: "Upsert is not a valid JSON object.",
		}
	}

	castAppend, okAppend := params.Append.(map[string]interface{})
	if okAppend == false {
		return 0, &data.Error{
			Code: http.StatusBadRequest,
			Message: "Append is not a valid JSON object.",
		}
	}

	var specs []gocb.MutateInSpec

	if castRemove != nil {
		for _, value := range castRemove {
			if value == "general" ||
				value == "general.name" ||
				value == "general.size" ||
				value == "general.creation_time" ||
				value == "general.format" ||
				value == "url" {
				return 0, &data.Error{
					Code: http.StatusNotAcceptable,
					Message: fmt.Sprintf("You are trying to delete %s, which is a critical metadata.", value),
				}
			}

			specs = append(specs, gocb.RemoveSpec(value, &gocb.RemoveSpecOptions{}))
		}
	}

	if castUpsert != nil {
		upsertSpecs, err := parseUpsertKeys(castUpsert, "")
		if err != nil {
			return 0, err
		}

		for _, value := range upsertSpecs {
			key, ok := value[0].(string)
			if ok == false {
				return 0, &data.Error{
					Code: http.StatusInternalServerError,
					Message: "Unable to parse key.",
				}
			}

			specs = append(specs, gocb.UpsertSpec(key, value[1], &gocb.UpsertSpecOptions{}))
		}
	}

	if castAppend != nil {
		for key, value := range castAppend {
			_, hasMultiple := value.([]interface{})
			specs = append(specs, gocb.ArrayAppendSpec(key, value, &gocb.ArrayAppendSpecOptions{
				HasMultiple: hasMultiple,
			}))
		}
	}

	_, err := gcfs.Cluster.
		Bucket(config.Main.Database.BucketName).
		DefaultCollection().
		MutateIn(id, specs, &gocb.MutateInOptions{})

	if err != nil {
		return 0, &data.Error{
			Code: http.StatusInternalServerError,
			Message: err.Error(),
		}
	}

	return uint64(time.Now().UnixNano() / int64(time.Millisecond)), nil
}
