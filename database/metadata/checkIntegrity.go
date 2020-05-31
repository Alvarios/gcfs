package metadata

import (
	"fmt"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/config/data"
	"reflect"
)

func CheckCustomMeta(
	provided interface{},
	expected map[string]interface{},
	currentKey string,
) *data.Error {
	providedMap, ok := provided.(map[string]interface{})

	// Expect a json like object.
	if ok == false {
		return &data.Error{
			Message: fmt.Sprintf(
				"Unexpected metadata value at %s: expected a nested json object, but got %T.",
				currentKey,
				reflect.TypeOf(providedMap),
			),
			Code: 400,
		}
	}

	for key, value := range expected {
		// Key is absent so metadata are incorrect.
		if _, ok := providedMap[key]; ok == false {
			return &data.Error{
				Message: fmt.Sprintf("Missing key %s in %s.", key, currentKey),
				Code:    400,
			}
		}

		// If a string is the value, then it represent the expected type. Otherwise we expect another map.
		expectedType, ok := value.(string)

		// Expect a map.
		if ok == false {
			newMap, ok := value.(map[string]interface{})

			// Other values are not allowed.
			if ok == false {
				return &data.Error{
					Message: fmt.Sprintf(
						"Wrong expectation map provided: key %s is not a string or a nested object.",
						currentKey,
					),
					Code: 500,
				}
			}

			err := CheckCustomMeta(providedMap[key], newMap, currentKey + "." + key)

			// Nested error breaks the loop.
			if err != nil {
				return err
			}
		}

		// Check if type matches.
		t := reflect.TypeOf(providedMap[key]).String()
		if t != expectedType {
			return &data.Error{
				Message: fmt.Sprintf(
					"Unexpected metadata value at %s.%s: expected %s, but got %s.",
					currentKey,
					key,
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

func CheckIntegrity(metadata interface{}) (*fileMetadata, *data.Error) {
	required, ok := metadata.(fileMetadata)

	if ok == false {
		return nil, &data.Err.Metadata.Invalid
	}

	if required.Url == "" {
		return nil, &data.Err.Metadata.MissingUrl
	}

	if required.General.Name == "" {
		return nil, &data.Err.Metadata.General.MissingName
	}

	if required.General.Format == "" {
		return nil, &data.Err.Metadata.General.MissingFormat
	}

	if required.General.CreationTime == 0 {
		return nil, &data.Err.Metadata.General.MissingCreationTime
	}

	if required.General.ModificationTime < required.General.CreationTime {
		return nil, &data.Err.Metadata.General.InvalidModificationTime
	}

	var err *data.Error

	// User can add some required metadata to check for.
	if config.Main.Metadata != nil {
		// Function is recursive.
		err = CheckCustomMeta(metadata, config.Main.Metadata, "Metadata")
	}

	return &required, err
}
