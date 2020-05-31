package data

import (
	"fmt"
	"github.com/jinzhu/configor"
)

var Err Errors

type MetadataGeneralErrors struct {
	MissingName             Error `json:"missing_name"`
	MissingFormat           Error `json:"missing_format"`
	MissingCreationTime     Error `json:"missing_creation_time"`
	InvalidModificationTime Error `json:"invalid_modification_time"`
}

type MetadataErrors struct {
	Invalid    Error                 `json:"empty"`
	MissingUrl Error                 `json:"missing_url"`
	General    MetadataGeneralErrors `json:"general"`
}

type Errors struct {
	Metadata MetadataErrors
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("Code %q: %s", e.Code, e.Message)
}

func init() {
	if err := configor.Load(&Err, "./errors.ENV.json"); err != nil {
		panic(err.Error())
	}
}