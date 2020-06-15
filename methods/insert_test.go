package methods

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"regexp"
	"testing"
)

func TestInsertf(t *testing.T) {
	config.LoadConfigForTest(config.Configuration{})
	err := database.Connect()
	if err != nil {
		t.Errorf("cannot connect to database : %s", err.Error())
		return
	}

	data := map[string]interface{}{
		"url": "path/to/my/file",
		"views": 0,
		"general": map[string]interface{}{
			"name": "my awesome file",
			"format": "txt",
			"size": 2048,
		},
	}

	_, err = InsertF(data, "", InsertFlags{AutoProvide: false})
	if err == (*responses.Error)(nil) {
		t.Error("insertf doesn't reject when provided incomplete metadata")
	}

	fileId, err := InsertF(data, "my_awesome_file_id", InsertFlags{AutoProvide: true})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to insert data with custom id : %s", err.Error())
	}

	if fileId != "my_awesome_file_id" {
		t.Errorf("insert with custom id failed : expected id to be ' my_awesome_file_id', got %s", fileId)
	}

	err = Delete(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to delete file %s : %s", fileId, err.Error())
	}

	fileId, err = InsertF(data, "my_awesome_file_id_{ts}", InsertFlags{AutoProvide: true})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to insert data with custom id : %s", err.Error())
	}

	reId := regexp.MustCompile(`my_awesome_file_id_\d+`)
	if !reId.MatchString(fileId) {
		t.Errorf("insert with custom id failed : expected id to be ' my_awesome_file_id_' + timestamp, got %s", fileId)
	}

	err = Delete(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to delete file %s : %s", fileId, err.Error())
	}

	err = database.Cluster.Close(nil)
	if err != nil {
		t.Errorf("cannot close database instance : %s", err.Error())
		return
	}
}

func TestInsert(t *testing.T) {
	config.LoadConfigForTest(config.Configuration{})
	config.Main.Global.AutoProvide = false
	err := database.Connect()
	if err != nil {
		t.Errorf("cannot connect to database : %s", err.Error())
		return
	}

	data := map[string]interface{}{
		"url": "path/to/my/file",
		"views": 0,
		"general": map[string]interface{}{
			"name": "my awesome file",
			"format": "txt",
			"size": 2048,
		},
	}

	_, err = Insert(data, "")
	if err == (*responses.Error)(nil) {
		t.Error("insert doesn't reject when provided incomplete metadata")
	}

	config.Main.Global.AutoProvide = true
	fileId, err := Insert(data, "")
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to insert data : %s", err.Error())
	}

	err = Delete(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to delete file %s : %s", fileId, err.Error())
	}

	err = database.Cluster.Close(nil)
	if err != nil {
		t.Errorf("cannot close database instance : %s", err.Error())
		return
	}
}