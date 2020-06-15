package methods

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/gcfs/database/metadata"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"reflect"
	"testing"
)

func TestUpdate(t *testing.T) {
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

	fileId, err := InsertF(data, "", InsertFlags{AutoProvide: true})
	if err != (*responses.Error)(nil) {
		t.Fatalf("unable to insert document in database : %s", err.Error())
	}

	fileData, err := Get(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't get file from database : %s", err.Error())
	} else if fileData == nil {
		t.Error("fetch returned empty data")
	}

	mapFileData, ok := fileData.(map[string]interface{})
	if ok == false {
		t.Error("unable to cast fetched data to valid map[string]")
	}

	if views, ok := mapFileData["views"].(float64); ok == false || views != 0 {
		t.Errorf("fetch returned wrong views : expected 0, got %v", mapFileData["views"])
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"views"},
	})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to remove field : %s", err.Error())
	}

	fileData, err = Get(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't get file from database : %s", err.Error())
	} else if fileData == nil {
		t.Error("fetch returned empty data")
	}

	_, err = metadata.CheckIntegrity(fileData)
	if err != (*responses.Error)(nil) {
		t.Errorf("returned data failed integrity check : %s", err.Error())
	}

	mapFileData, ok = fileData.(map[string]interface{})
	if ok == false {
		t.Error("unable to cast fetched data to valid map[string]")
	}

	if _, ok := mapFileData["views"]; ok == true {
		t.Error("deletion failed on update : field views is still present in database")
	}

	_, err = Update(fileId, UpdateSpec{
		Upsert: map[string]interface{}{
			"slice_key": []int{10,20,30},
			"general": map[string]interface{}{
				"name": "my new awesome file",
			},
		},
	})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to upsert fields : %s", err.Error())
	}

	fileData, err = Get(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't get file from database : %s", err.Error())
	} else if fileData == nil {
		t.Error("fetch returned empty data")
	}

	mapFileData, ok = fileData.(map[string]interface{})
	if ok == false {
		t.Error("unable to cast fetched data to valid map[string]")
	}

	fileMetadata, err := metadata.CheckIntegrity(fileData)
	if err != (*responses.Error)(nil) {
		t.Errorf("returned data failed integrity check : %s", err.Error())
	}

	if fileMetadata == nil {
		t.Error("CheckIntegrity returned nil parsed metadata")
		return
	}

	if fileMetadata.General.Name != "my new awesome file" {
		t.Errorf(
			"upsert failed on update : expected general.name to be 'my new awesome file', got %s",
			fileMetadata.General.Name,
		)
	}

	if slice, ok := mapFileData["slice_key"].([]interface{}); ok == true {
		if slice[0] != float64(10) || slice[1] != float64(20) || slice[2] != float64(30) {
			t.Errorf("fetch returned wrong values for slice_key : expected [10 20 30], got %v", slice)
		}
	} else {
		t.Errorf(
			"fetch returned wrong slice_key : expected []interface {}, got %s",
			reflect.TypeOf(mapFileData["slice_key"]).String(),
		)
	}

	_, err = Update(fileId, UpdateSpec{
		Upsert: map[string]interface{}{
			"general.name": "another name",
		},
		Append: map[string]interface{}{
			"slice_key": []int{40,50},
		},
	})
	if err != (*responses.Error)(nil) {
		t.Errorf("unable to upsert or append fields : %s", err.Error())
	}

	fileData, err = Get(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't get file from database : %s", err.Error())
	} else if fileData == nil {
		t.Error("fetch returned empty data")
	}

	mapFileData, ok = fileData.(map[string]interface{})
	if ok == false {
		t.Error("unable to cast fetched data to valid map[string]")
	}

	fileMetadata, err = metadata.CheckIntegrity(fileData)
	if err != (*responses.Error)(nil) {
		t.Errorf("returned data failed integrity check : %s", err.Error())
	}

	if fileMetadata == nil {
		t.Error("CheckIntegrity returned nil parsed metadata")
		return
	}

	if fileMetadata.General.Name != "another name" {
		t.Errorf(
			"upsert failed on update : expected general.name to be 'another name', got %s",
			fileMetadata.General.Name,
		)
	}

	if slice, ok := mapFileData["slice_key"].([]interface{}); ok == true {
		if slice[0] != float64(10) ||
			slice[1] != float64(20) ||
			slice[2] != float64(30) ||
			slice[3] != float64(40) ||
			slice[4] != float64(50) {
			t.Errorf("fetch returned wrong values for slice_key : expected [10 20 30 40 50], got %v", slice)
		}
	} else {
		t.Errorf(
			"fetch returned wrong slice_key : expected []interface {}, got %s",
			reflect.TypeOf(mapFileData["slice_key"]).String(),
		)
	}

	err = Delete(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't delete file from database : %s", err.Error())
	}

	err = database.Cluster.Close(nil)
	if err != nil {
		t.Errorf("cannot close database instance : %s", err.Error())
		return
	}
}

func TestUpdateDeleteError(t *testing.T) {
	config.LoadConfigForTest(config.Configuration{})
	err := database.Connect()
	if err != nil {
		t.Errorf("cannot connect to database : %s", err.Error())
		return
	}

	data := map[string]interface{}{
		"url": "path/to/my/file",
		"general": map[string]interface{}{
			"name": "my awesome file",
			"format": "txt",
			"size": 2048,
		},
	}

	fileId, err := InsertF(data, "", InsertFlags{AutoProvide: true})
	if err != (*responses.Error)(nil) {
		t.Fatalf("unable to insert document in database : %s", err.Error())
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"views"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove non existing path")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"general"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path general")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"general.name"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path general.name")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"general.format"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path general.format")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"general.size"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path general.size")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"general.modification_time"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path general.modification_time")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"general.creation_time"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path general.creation_time")
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"url"},
	})
	if err == (*responses.Error)(nil) {
		t.Errorf("returned no error when trying to remove core path url")
	}

	err = Delete(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't delete file from database : %s", err.Error())
	}

	err = database.Cluster.Close(nil)
	if err != nil {
		t.Errorf("cannot close database instance : %s", err.Error())
		return
	}
}