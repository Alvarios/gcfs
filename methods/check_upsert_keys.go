package methods

import (
	"fmt"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"reflect"
)

func checkKeyType(key string, value interface{}, expected interface{}, empty interface{}) *responses.Error {
	var ev []string
	if s, ok := expected.(string); ok {
		ev = []string{s}
	} else if m, ok := expected.([]string); ok {
		ev = m
	} else {
		return &responses.Error{
			Message: fmt.Sprintf(
				"wrong expectation value provided : expected []string or string, got %s",
				reflect.TypeOf(expected).String(),
			),
			Code: 500,
		}
	}

	valueType := reflect.TypeOf(value).String()
	match := false
	for _, e := range ev {
		if e == valueType {
			match = true
		}
	}

	if match == false {
		return &responses.Error{
			Message: fmt.Sprintf(
				"trying to update %s key with %v of type %s, which is forbidden",
				key,
				value,
				valueType,
			),
			Code: 400,
		}
	}

	if value == empty {
		return &responses.Error{
			Message: fmt.Sprintf("trying to assign empty value to %s key", key),
			Code: 400,
		}
	}

	return nil
}

func checkUpsertKeys(uk [][]interface{}) *responses.Error {
	for _, kv := range uk {
		key, ok := kv[0].(string)
		valueType := reflect.TypeOf(kv[1])
		if ok == false {
			return &responses.Error{
				Message: fmt.Sprintf(
					"non string key provided in upsert parameters : encountered %v of type %s",
					kv[0],
					reflect.TypeOf(kv[0]).String(),
				),
				Code: 400,
			}
		}

		if key == "" {
			return &responses.Error{
				Message: "empty key not allowed",
				Code: 400,
			}
		}

		// map[string] are flattened, so general is not valid.
		if key == "general" {
			return &responses.Error{
				Message: fmt.Sprintf(
					"trying to update general key with %v of type %s, which is forbidden",
					kv[1],
					valueType,
				),
				Code: 400,
			}
		}

		if key == "url" ||
			key == "general.name" ||
			key == "general.format" {
			err := checkKeyType(key, kv[1], "string", "")
			if err != (*responses.Error)(nil) {
				return err
			}
		}

		if key == "general.creation_time" {
			err := checkKeyType(key, kv[1], []string{"float64", "uint64", "int"}, 0)
			if err != (*responses.Error)(nil) {
				return err
			}
		}

		if key == "general.modification_time" {
			err := checkKeyType(key, kv[1], []string{"float64", "uint64", "int"}, nil)
			if err != (*responses.Error)(nil) {
				return err
			}
		}
	}

	return nil
}
