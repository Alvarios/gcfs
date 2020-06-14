package errors

import (
	"encoding/json"
	"github.com/Alvarios/kushuh-go-utils/file-utils"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"io/ioutil"
	"log"
)

var Err Errors

type MetadataGeneralErrors struct {
	MissingName             responses.Error `json:"missing_name"`
	MissingFormat           responses.Error `json:"missing_format"`
	MissingCreationTime     responses.Error `json:"missing_creation_time"`
	InvalidModificationTime responses.Error `json:"invalid_modification_time"`
}

type MetadataErrors struct {
	Invalid    responses.Error                 `json:"empty"`
	MissingUrl responses.Error                 `json:"missing_url"`
	General    MetadataGeneralErrors `json:"general"`
}

type Errors struct {
	Metadata MetadataErrors `json:"metadata"`
	Test responses.Error `json:"test"`
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