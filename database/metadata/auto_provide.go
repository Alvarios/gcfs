package metadata

import (
	"bytes"
	"encoding/json"
	"github.com/Alvarios/gcfs/config/errors"
	numberUtils "github.com/Alvarios/kushuh-go-utils/number-utils"
	"strconv"
)

func AutoProvide(metadata interface{}) (interface{}, *errors.Error) {
	// Nil values are castable to json, and we don't want that.
	if metadata == nil {
		return nil, &errors.Err.Metadata.Invalid
	}
	// i64, _ := strconv.ParseUint(string(n), 10, 64)

	jsonString, err := json.Marshal(metadata)
	if err != nil {
		return nil, &errors.Error{
			Message: err.Error(),
			Code: 500,
		}
	}

	output := make(map[string]interface{})
	d := json.NewDecoder(bytes.NewBuffer(jsonString))
	d.UseNumber()
	if err := d.Decode(&output); err != nil {
		return nil, &errors.Error{
			Message: err.Error(),
			Code: 500,
		}
	}

	timestamp := numberUtils.Timestamp()

	if mv, ok := output["general"].(map[string]interface{}); ok == true {
		if ct, ok := mv["creation_time"].(json.Number); ok {
			if ct == "0" {
				mv["creation_time"] = timestamp
			} else {
				mv["creation_time"], _ = strconv.ParseUint(string(ct), 10, 64)

				if mt, ok := mv["modification_time"].(json.Number); ok {
					if mt == "0" {
						mv["modification_time"] = timestamp
					} else {
						mv["modification_time"], _ = strconv.ParseUint(string(mt), 10, 64)
					}
				} else {
					mv["modification_time"] = timestamp
				}
			}
		} else {
			mv["creation_time"] = timestamp
		}

		output["general"] = mv
	} else {
		return nil, &errors.Err.Metadata.Invalid
	}

	return output, nil
}
