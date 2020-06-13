package metadata

import (
	"fmt"
	"github.com/Alvarios/gcfs/config/errors"
	"reflect"
	"strings"
)

func CheckCustomMetadata(
	provided interface{},
	expected map[string]interface{},
	currentKey string,
) *errors.Error {
	providedMap, ok := provided.(map[string]interface{})

	// Expect a map[string] object.
	if ok == false {
		return &errors.Error{
			Message: fmt.Sprintf(
				"unexpected metadata value at %s : expected map[string]interface {}, got %s",
				currentKey,
				reflect.TypeOf(provided).String(),
			),
			Code: 400,
		}
	}

	for key, value := range expected {
		buildKey := currentKey + "." + key
		initialPrefix := strings.Split(currentKey, ".")[0]

		// Disable overriding of default parameters.
		if buildKey == initialPrefix + ".url" ||
			buildKey == initialPrefix + ".general.name" ||
			buildKey == initialPrefix + ".general.format" ||
			buildKey == initialPrefix + ".general.size" ||
			buildKey == initialPrefix + ".general.creation_time" ||
			buildKey == initialPrefix + ".general.modification_time" {
			return &errors.Error{
				Message: fmt.Sprintf("trying to redefine the core key %s", buildKey),
				Code:    400,
			}
		}

		// Key is absent so metadata are incorrect.
		if _, ok := providedMap[key]; ok == false {
			return &errors.Error{
				Message: fmt.Sprintf("missing key %s in %s", key, currentKey),
				Code:    400,
			}
		}

		// If a string is the value, then it represent the expected type. Otherwise we expect another map.
		newMap, ok := value.(map[string]interface{})

		// Nested check has to be performed
		if ok == true {
			err := CheckCustomMetadata(providedMap[key], newMap, buildKey)

			// Nested error breaks the loop.
			if err != (*errors.Error)(nil) {
				return err
			}

			continue
		}

		expectedType, ok := value.(string)
		if ok == false {
			return &errors.Error{
				Message: fmt.Sprintf(
					"wrong expectation map provided : key %s is not a string or a map[string]",
					buildKey,
				),
				Code: 500,
			}
		}

		// Check if type matches.
		t := reflect.TypeOf(providedMap[key]).String()
		if t != expectedType {
			return &errors.Error{
				Message: fmt.Sprintf(
					"unexpected metadata value at %s : expected %s, got %s",
					buildKey,
					expectedType,
					t,
				),
				Code: 400,
			}
		}
	}

	// No error. Great !
	return nil
}
