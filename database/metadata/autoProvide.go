package metadata

import (
	"github.com/Alvarios/gcfs/config/data"
	"time"
)

func AutoProvide(metadata interface{}) (map[string]interface{}, *data.Error) {
	v, ok := metadata.(map[string]interface{})

	if ok == false {
		return nil, &data.Err.Metadata.Invalid
	}

	// General object has to be provided, as it contains some data that cannot be automatically generated.
	if mv, ok := v["general"].(map[string]interface{}); ok == true {
		if _, ok := mv["creation_time"]; ok == false {
			mv["creation_time"] = time.Now().UnixNano() / int64(time.Millisecond)
		} else if _, ok := mv["modification_time"]; ok == false {
			mv["modification_time"] = time.Now().UnixNano() / int64(time.Millisecond)
		}
	}

	return v, nil
}
