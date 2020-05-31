package api

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/methods"
	"github.com/Alvarios/kushuh-go-utils/map_utils"
	nefts_config "github.com/Alvarios/nefts-go/config"
	"net/http"
)

func Search(w http.ResponseWriter, r *http.Request) {
	// Load file data from body.
	var file interface{}
	err := json.NewDecoder(r.Body).Decode(&file)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	mapFile, ok := file.(map[string]interface{})
	if ok == false {
		http.Error(w, "Body is not valid JSON format.", http.StatusBadRequest)
		return
	}

	assertOptions := config.Main.Routes.SearchOptions

	if userOptions, ok := mapFile["options"].(map[string]interface{}); ok {
		mergeOptions, err := map_utils.Merge(config.Main.Routes.SearchOptions, userOptions)
		if err != nil {
			http.Error(w, "Body is not valid JSON format.", http.StatusBadRequest)
			return
		}

		assertOptions, ok = mergeOptions.(nefts_config.Options)
		if ok == false {
			http.Error(w, "Cannot parse options.", http.StatusInternalServerError)
			return
		}
	}

	start := int64(0)
	if startValue, ok := mapFile["start"].(int64); ok {
		start = startValue
	}

	end := int64(30)
	if endValue, ok := mapFile["end"].(int64); ok {
		end = endValue
	}

	queryResults, qerr := methods.Search(start, end, assertOptions)

	if qerr != nil {
		http.Error(w, qerr.Message, qerr.Code)
		return
	}

	b, err := json.Marshal(queryResults)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// write response
	_, err = w.Write(b)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	return
}
