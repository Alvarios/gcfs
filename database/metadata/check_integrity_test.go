package metadata

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/config/errors"
	"github.com/Alvarios/kushuh-go-utils/number-utils"
	"testing"
)

func TestCheckIntegrityArguments(t *testing.T) {
	config.LoadConfig(config.Configuration{})

	// Should reject non interface/map arguments.
	_, err := CheckIntegrity(nil)
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error for nil argument call.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with nil parameter : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity("fakeString")
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error for string argument call.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with string parameter : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(645)
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error for number argument call.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with number parameter : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(true)
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error for boolean argument call.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with boolean parameter : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}
}

func TestCheckIntegrityWithCastError(t *testing.T) {
	// Should reject invalid when any option is given an incorrect value (cannot cast to metadata).
	_, err := CheckIntegrity(map[string]interface{}{
		"url": 998,
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non string url.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non string url : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"general": 998,
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non map general.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non map general : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"general": map[string]interface{}{
			"name": []string{},
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non string general.name.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non string general.name : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"general": map[string]interface{}{
			"format": []string{},
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non string general.format.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non string general.format : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"general": map[string]interface{}{
			"creation_time": "timestamp",
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non number general.creation_time.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non number general.creation_time : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"general": map[string]interface{}{
			"modification_time": "timestamp",
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non number general.modification_time.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non number general.modification_time : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"general": map[string]interface{}{
			"size": "size in string",
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error with non number general.size.")
	} else if err != nil && err.Message != errors.Err.Metadata.Invalid.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error with non number general.size : expected '%s', got '%s'",
			errors.Err.Metadata.Invalid.Message,
			err.Message,
		)
	}
}

func TestCheckIntegrityWithWrongMetadata(t *testing.T) {
	// Should reject noUrl.
	_, err := CheckIntegrity(map[string]interface{}{
		"general": map[string]interface{}{
			"name":          "my awesome video.",
			"format":        "mkv",
			"creation_time": numberUtils.Timestamp(),
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error when missing url.")
	} else if err != nil && err.Message != errors.Err.Metadata.MissingUrl.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error when missingurl : expected '%s', got '%s'",
			errors.Err.Metadata.MissingUrl.Message,
			err.Message,
		)
	}

	// Should reject noName.
	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"format":        "mkv",
			"creation_time": numberUtils.Timestamp(),
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error when missing general.name.")
	} else if err != nil && err.Message != errors.Err.Metadata.General.MissingName.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error when missing general.name : expected '%s', got '%s'",
			errors.Err.Metadata.General.MissingName.Message,
			err.Message,
		)
	}

	// Should reject noFormat.
	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"name":          "my awesome video.",
			"creation_time": numberUtils.Timestamp(),
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error when missing general.format.")
	} else if err != nil && err.Message != errors.Err.Metadata.General.MissingFormat.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error when missing general.format : expected '%s', got '%s'",
			errors.Err.Metadata.General.MissingFormat.Message,
			err.Message,
		)
	}

	// Should reject noCreationTime.
	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"name":              "my awesome video.",
			"format":            "mkv",
			"modification_time": numberUtils.Timestamp(),
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error when missing general.creation_time.")
	} else if err != nil && err.Message != errors.Err.Metadata.General.MissingCreationTime.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error when missing general.creation_time : expected '%s', got '%s'",
			errors.Err.Metadata.General.MissingCreationTime.Message,
			err.Message,
		)
	}

	// Should reject invalidModificationTime.
	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"name":              "my awesome video.",
			"format":            "mkv",
			"creation_time":     numberUtils.Timestamp(),
			"modification_time": numberUtils.Timestamp() - 1000,
		},
	})
	if err == (*errors.Error)(nil) {
		t.Error("CheckIntegrity returned no error when general.modification_time is less than general.creation_time.")
	} else if err != nil && err.Message != errors.Err.Metadata.General.InvalidModificationTime.Message {
		t.Errorf(
			"CheckIntegrity returned wrong error when general.modification_time is less than general.creation_time : expected '%s', got '%s'",
			errors.Err.Metadata.General.InvalidModificationTime.Message,
			err.Message,
		)
	}
}

func TestCheckIntegrity(t *testing.T) {
	// All following tests should run.
	_, err := CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"name": "my awesome video",
			"format": "mkv",
			"creation_time": numberUtils.Timestamp(),
		},
	})

	if err != (*errors.Error)(nil) {
		t.Errorf("error checking metadata integrity : %s", err.Error())
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"name": "my awesome video",
			"format": "mkv",
			"creation_time": numberUtils.Timestamp(),
		},
	})

	if err != (*errors.Error)(nil) {
		t.Errorf("error checking metadata integrity : %s", err.Error())
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"key": "value",
		"general": map[string]interface{}{
			"name": "my awesome video",
			"format": "mkv",
			"creation_time": numberUtils.Timestamp() - 1000,
			"modification_time": numberUtils.Timestamp(),
			"another_key": "test",
		},
	})

	if err != (*errors.Error)(nil) {
		t.Errorf("error checking metadata integrity : %s", err.Error())
	}

	config.LoadConfig(config.Configuration{
		Metadata: map[string]interface{}{
			"required_key": "string",
		},
	})
	// Further test with custom metadata are provided in check_custom_metadata_test.go file.
	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"required_key": "value",
		"general": map[string]interface{}{
			"name": "my awesome video",
			"format": "mkv",
			"creation_time": numberUtils.Timestamp() - 1000,
			"another_key": "test",
		},
	})

	if err != (*errors.Error)(nil) {
		t.Errorf("error checking metadata integrity : %s", err.Error())
	}

	_, err = CheckIntegrity(map[string]interface{}{
		"url": "/path/to/myfile",
		"general": map[string]interface{}{
			"name": "my awesome video",
			"format": "mkv",
			"creation_time": numberUtils.Timestamp() - 1000,
			"modification_time": numberUtils.Timestamp(),
			"another_key": "test",
		},
	})

	if err == (*errors.Error)(nil) {
		t.Errorf("error checking metadata integrity : no error raised when not providing custom required metadata")
	}
}