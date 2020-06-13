package metadata

import (
	"encoding/json"
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/config/errors"
)

func CheckIntegrity(metadata interface{}) (*fileMetadata, *errors.Error) {
	// Nil values are castable to json, and we don't want that.
	if metadata == nil {
		return nil, &errors.Err.Metadata.Invalid
	}

	jsonString, _ := json.Marshal(metadata)
	required := &fileMetadata{}
	err := json.Unmarshal(jsonString, required)

	if err != nil {
		return nil, &errors.Err.Metadata.Invalid
	}

	if required.Url == "" {
		return nil, &errors.Err.Metadata.MissingUrl
	}

	if required.General.Name == "" {
		return nil, &errors.Err.Metadata.General.MissingName
	}

	if required.General.Format == "" {
		return nil, &errors.Err.Metadata.General.MissingFormat
	}

	if required.General.CreationTime == 0 {
		return nil, &errors.Err.Metadata.General.MissingCreationTime
	}

	if required.General.ModificationTime > 0 && required.General.ModificationTime < required.General.CreationTime {
		return nil, &errors.Err.Metadata.General.InvalidModificationTime
	}

	// User can add some required metadata to check for.
	if config.Main.Metadata != nil {
		// Function is recursive.
		err := CheckCustomMetadata(metadata, config.Main.Metadata, "Metadata")

		return required, err
	}

	return required, nil
}
