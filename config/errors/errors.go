package errors

import (
	"encoding/json"
	"fmt"
	fileUtils "github.com/Alvarios/kushuh-go-utils/file-utils"
	"io/ioutil"
	"log"
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
	Metadata MetadataErrors `json:"metadata"`
	Test Error `json:"test"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("code %v : %s", e.Code, e.Message)
}

func LoadErrors() {
	// Open our jsonFile.
	jsonFile, err := fileUtils.OpenFromProjectRoot("gcfs", "config/errors/errors.json")
	if err != nil {
		log.Fatalf("cannot open error config file : %s", err.Error())
		return
	}

	// Read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// Cast byte array to struct.
	err = json.Unmarshal(byteValue, &Err)
	if err != nil {
		log.Fatalf("cannot parse error config file : %s", err.Error())
		return
	}

	jsonFile.Close()
	return
}