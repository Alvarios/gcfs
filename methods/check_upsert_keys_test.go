package methods

import (
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"testing"
)

func TestCheckKeyType(t *testing.T) {
	err := checkKeyType("", "a value", "string", "")
	if err != (*responses.Error)(nil) {
		t.Errorf("error while checking data : %s", err.Error())
	}

	err = checkKeyType("", 10, []string{"float64", "uint64", "int"}, 0)
	if err != (*responses.Error)(nil) {
		t.Errorf("error while checking data : %s", err.Error())
	}

	err = checkKeyType("", 0, []string{"float64", "uint64", "int"}, 0)
	if err == (*responses.Error)(nil) {
		t.Error("returned no error when checking nil type")
	} else if err.Message != "trying to assign empty value to  key" {
		t.Errorf("unexpected error message : expected 'trying to assign empty value to  key', got %s", err.Message)
	}

	err = checkKeyType("", 0, []string{"float64", "uint64", "int"}, nil)
	if err != (*responses.Error)(nil) {
		t.Errorf("returned error when value not matching empty value : '%s'", err.Error())
	}

	err = checkKeyType("", 0, "string", "")
	if err == (*responses.Error)(nil) {
		t.Error("returned no error when checking non matching type")
	} else if err.Message != "trying to update  key with 0 of type int, which is forbidden" {
		t.Errorf("unexpected error message : expected 'trying to update  key with 0 of type int, which is forbidden', got '%s'", err.Message)
	}
}

func TestCheckUpsertKeys(t *testing.T) {
	err := checkUpsertKeys(map[string]interface{}{
		"key": "value",
	})
	if err != (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys failed : %s", err.Error())
	}

	err = checkUpsertKeys(map[string]interface{}{
		"key": "value",
		"url": "",
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys allowed non valid url value")
	} else if err.Message != "trying to assign empty value to url key" {
		t.Errorf("unexpected error message : expected 'trying to assign empty value to url key', got %s", err.Message)
	}

	err = checkUpsertKeys(map[string]interface{}{
		"key": "value",
		"general.name": "",
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys allowed non valid general.name value")
	} else if err.Message != "trying to assign empty value to general.name key" {
		t.Errorf("unexpected error message : expected 'trying to assign empty value to general.name key', got %s", err.Message)
	}

	err = checkUpsertKeys(map[string]interface{}{
		"key": "value",
		"general.format": "",
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys allowed non valid general.format value")
	} else if err.Message != "trying to assign empty value to general.format key" {
		t.Errorf("unexpected error message : expected 'trying to assign empty value to general.format key', got %s", err.Message)
	}

	err = checkUpsertKeys(map[string]interface{}{
		"key": "value",
		"general.creation_time": "",
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys allowed non valid general.creation_time value")
	} else if err.Message != "trying to update general.creation_time key with  of type string, which is forbidden" {
		t.Errorf("unexpected error message : expected 'trying to update general.creation_time key with  of type string, which is forbidden', got %s", err.Message)
	}

	err = checkUpsertKeys(map[string]interface{}{
		"key": "value",
		"general.modification_time": "",
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys allowed non valid general.modification_time value")
	} else if err.Message != "trying to update general.modification_time key with  of type string, which is forbidden" {
		t.Errorf("unexpected error message : expected 'trying to update general.modification_time key with  of type string, which is forbidden', got %s", err.Message)
	}

	err = checkUpsertKeys(map[string]interface{}{
		"key": "value",
		"general": 18,
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("checkUpsertKeys allowed non valid general value")
	} else if err.Message != "trying to update general key with 18 of type int, which is forbidden" {
		t.Errorf("unexpected error message : expected 'trying to update general key with 18 of type int, which is forbidden', got %s", err.Message)
	}
}