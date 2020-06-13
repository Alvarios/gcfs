package metadata

import (
	"github.com/Alvarios/gcfs/config/errors"
	"testing"
)

func TestCheckCustomMetadataWithWrongExpectationMap(t *testing.T) {
	// 1-level deep map.
	expected := map[string]interface{}{
		"boolean_key": true,
		"uint64_key": "uint64",
		"int_key": "int",
		"string_key": "string",
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": "map[string]int",
	}
	value := map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"string_key": "fakeValue",
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	expectedError := "wrong expectation map provided : key .boolean_key is not a string or a map[string]"
	err := CheckCustomMetadata(value, expected, "")
	if err == (*errors.Error)(nil) {
		t.Error("test pass with boolean value in expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}

	expected = map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": uint64(324),
		"int_key": "int",
		"string_key": "string",
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": "map[string]int",
	}
	value = map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"string_key": "fakeValue",
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	expectedError = "wrong expectation map provided : key .uint64_key is not a string or a map[string]"
	err = CheckCustomMetadata(value, expected, "")
	if err == (*errors.Error)(nil) {
		t.Error("test pass with uint64 value in expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}

	expected = map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": "uint64",
		"int_key": "int",
		"string_key": "string",
		"char_key": "int32",
		"slice_key": []interface{}{},
		"map_key": "map[string]int",
	}
	value = map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"string_key": "fakeValue",
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	expectedError = "wrong expectation map provided : key .slice_key is not a string or a map[string]"
	err = CheckCustomMetadata(value, expected, "")
	if err == (*errors.Error)(nil) {
		t.Error("test pass with []interface {} value in expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}

	expected = map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": "uint64",
		"int_key": "int",
		"string_key": "string",
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": map[string]int{},
	}
	value = map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"string_key": "fakeValue",
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	expectedError = "wrong expectation map provided : key .map_key is not a string or a map[string]"
	err = CheckCustomMetadata(value, expected, "")
	if err == (*errors.Error)(nil) {
		t.Error("test pass with map[string]int value in expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}
}

func TestCheckCustomMetadataWithWrongValue(t *testing.T) {
	// 1-level deep map.
	expected := map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": "uint64",
		"int_key": "int",
		"string_key": "string",
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": "map[string]int",
	}
	value := map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"string_key": []string{"a", "b", "c"},
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	err := CheckCustomMetadata(value, expected, "")
	expectedError := "unexpected metadata value at .string_key : expected string, got []string"
	if err == (*errors.Error)(nil) {
		t.Error("test pass with value non matching expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}

	// x-level deep map.
	expected = map[string]interface{}{
		"key": "string",
		"map": map[string]interface{}{
			"subkey": "int",
			"submap": map[string]interface{}{
				"abyss": "bool",
			},
		},
	}
	value = map[string]interface{}{
		"key": "a string",
		"map": map[string]interface{}{
			"subkey": 42,
			"submap": map[string]interface{}{
				"abyss": "ok",
			},
		},
	}
	err = CheckCustomMetadata(value, expected, "")
	expectedError = "unexpected metadata value at .map.submap.abyss : expected bool, got string"
	if err == (*errors.Error)(nil) {
		t.Error("test pass with value non matching expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}

	// Empty map is equivalent to the string type.
	expected = map[string]interface{}{
		"key": "string",
		"map": map[string]interface{}{},
	}
	value = map[string]interface{}{
		"key": "a string",
		"map": []string{},
	}
	err = CheckCustomMetadata(value, expected, "")
	expectedError = "unexpected metadata value at .map : expected map[string]interface {}, got []string"
	if err == (*errors.Error)(nil) {
		t.Error("test pass with value non matching expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}
}

func TestCheckCustomMetadataWithCoreRedefinition(t *testing.T) {
	expected := map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": "uint64",
		"int_key": "int",
		"url": "string",
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": "map[string]int",
	}
	value := map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"url": "/path/to/file",
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	err := CheckCustomMetadata(value, expected, "")
	expectedError := "trying to redefine the core key .url"
	if err == (*errors.Error)(nil) {
		t.Error("test pass with value non matching expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}

	expected = map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": "uint64",
		"int_key": "int",
		"general": map[string]interface{}{
			"name": "int",
		},
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": "map[string]int",
	}
	value = map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"general": map[string]interface{}{
			"name": 200,
		},
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	err = CheckCustomMetadata(value, expected, "")
	expectedError = "trying to redefine the core key .general.name"
	if err == (*errors.Error)(nil) {
		t.Error("test pass with value non matching expectation map")
	} else if err.Message != expectedError {
		t.Errorf("wrong error returned : expected '%s', got '%s'", expectedError, err.Message)
	}
}

func TestCheckCustomMetadata(t *testing.T) {
	// All below tests should pass.

	// 1-level deep map.
	expected := map[string]interface{}{
		"boolean_key": "bool",
		"uint64_key": "uint64",
		"int_key": "int",
		"string_key": "string",
		"char_key": "int32",
		"slice_key": "[]interface {}",
		"map_key": "map[string]int",
	}
	value := map[string]interface{}{
		"boolean_key": true,
		"uint64_key": uint64(42),
		"int_key": -20,
		"string_key": "fakeValue",
		"char_key": 'x',
		"slice_key": []interface{}{2080, "a value"},
		"map_key": map[string]int{
			"from": 2048,
			"to":   2056,
		},
	}
	err := CheckCustomMetadata(value, expected, "")
	if err != (*errors.Error)(nil) {
		t.Error(err.Error())
	}

	// x-level deep map.
	expected = map[string]interface{}{
		"key": "string",
		"map": map[string]interface{}{
			"subkey": "int",
			"submap": map[string]interface{}{
				"abyss": "bool",
			},
		},
	}
	value = map[string]interface{}{
		"key": "a string",
		"map": map[string]interface{}{
			"subkey": 42,
			"submap": map[string]interface{}{
				"abyss": true,
			},
		},
	}
	err = CheckCustomMetadata(value, expected, "")
	if err != (*errors.Error)(nil) {
		t.Error(err.Error())
	}

	// Empty map is equivalent to the string type.
	expected = map[string]interface{}{
		"key": "string",
		"map": map[string]interface{}{},
	}
	value = map[string]interface{}{
		"key": "a string",
		"map": map[string]interface{}{
			"subkey": 42,
			"submap": map[string]interface{}{
				"abyss": true,
			},
		},
	}
	err = CheckCustomMetadata(value, expected, "")
	if err != (*errors.Error)(nil) {
		t.Error(err.Error())
	}
}
