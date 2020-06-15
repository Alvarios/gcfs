package methods

import (
	"github.com/Alvarios/gcfs/config"
	"github.com/Alvarios/gcfs/database"
	"github.com/Alvarios/gcfs/database/metadata"
	"github.com/Alvarios/kushuh-go-utils/router-utils/responses"
	"testing"
)

func TestDbInteractionBasicCycle(t *testing.T) {
	config.LoadConfigForTest(config.Configuration{})
	err := database.Connect()
	if err != nil {
		t.Errorf("cannot connect to database : %s", err.Error())
		return
	}

	// Basic Insert - Get - Update - Get - Delete test.
	originalData := map[string]interface{}{
		"url": "path/to/my/file",
		"views": 0,
		"general": map[string]interface{}{
			"name": "my awesome file",
			"format": "txt",
			"size": 2048,
		},
	}

	// Return on test since one failure will broke the entire cycle.
	fileId, err := InsertF(originalData, "", InsertFlags{AutoProvide: true})
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't insert file into database : %s", err.Error())
		return
	} else if fileId == "" {
		t.Error("insertion returned empty fileId")
		return
	}

	fileData, err := Get(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't get file from database : %s", err.Error())
	} else if fileData == nil {
		t.Error("fetch returned empty data")
	}

	fileMetadata, err := metadata.CheckIntegrity(fileData)
	if err != (*responses.Error)(nil) {
		t.Errorf("returned data failed integrity check : %s", err.Error())
	}

	if fileMetadata == nil {
		t.Error("CheckIntegrity returned nil parsed metadata")
		return
	}

	mapFileData, ok := fileData.(map[string]interface{})
	if ok == false {
		t.Error("unable to cast fetched data to valid map[string]")
	}

	if fileMetadata.Url != "path/to/my/file" {
		t.Errorf("fetch returned wrong url : expected 'path/to/my/file', got %s", fileMetadata.Url)
	}

	if fileMetadata.General.Name != "my awesome file" {
		t.Errorf("fetch returned wrong general.name : expected 'my awesome file', got %s", fileMetadata.General.Name)
	}

	if fileMetadata.General.Format != "txt" {
		t.Errorf("fetch returned wrong general.format : expected 'txt', got %s", fileMetadata.General.Format)
	}

	if fileMetadata.General.Size != 2048 {
		t.Errorf("fetch returned wrong general.size : expected 2048, got %v", fileMetadata.General.Size)
	}

	if views, ok := mapFileData["views"].(float64); ok == false || views != 0 {
		t.Errorf("fetch returned wrong views : expected 0, got %v", mapFileData["views"])
	}

	_, err = Update(fileId, UpdateSpec{
		Remove: []string{"views"},
		Upsert: map[string]interface{}{
			"likes":    20,
			"dislikes": 2,
		},
	})

	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't update file into database : %s", err.Error())
	}

	fileData, err = Get(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't get file from database : %s", err.Error())
	} else if fileData == nil {
		t.Error("fetch returned empty data")
	}

	fileMetadata, err = metadata.CheckIntegrity(fileData)
	if err != (*responses.Error)(nil) {
		t.Errorf("returned data failed integrity check : %s", err.Error())
	}

	if fileMetadata == nil {
		t.Error("CheckIntegrity returned nil parsed metadata")
		return
	}

	mapFileData, ok = fileData.(map[string]interface{})
	if ok == false {
		t.Error("unable to cast fetched data to valid map[string]")
	}

	if _, ok := mapFileData["views"]; ok == true {
		t.Error("deletion failed on update : field views is still present in database")
	}

	if fileMetadata.Url != "path/to/my/file" {
		t.Errorf("fetch returned wrong url : expected 'path/to/my/file', got %s", fileMetadata.Url)
	}

	if fileMetadata.General.Name != "my awesome file" {
		t.Errorf("fetch returned wrong general.name : expected 'my awesome file', got %s", fileMetadata.General.Name)
	}

	if fileMetadata.General.Format != "txt" {
		t.Errorf("fetch returned wrong general.format : expected 'txt', got %s", fileMetadata.General.Format)
	}

	if fileMetadata.General.Size != 2048 {
		t.Errorf("fetch returned wrong general.size : expected 2048, got %v", fileMetadata.General.Size)
	}

	if l, ok := mapFileData["likes"].(float64); ok == false || l != 20 {
		t.Errorf("fetch returned wrong likes : expected 20, got %v", mapFileData["likes"])
	}

	if dl, ok := mapFileData["dislikes"].(float64); ok == false || dl != 2 {
		t.Errorf("fetch returned wrong dislikes : expected 2, got %v", mapFileData["dislikes"])
	}

	err = Delete(fileId)
	if err != (*responses.Error)(nil) {
		t.Errorf("couldn't delete file from database : %s", err.Error())
	}

	_, err = Get(fileId)
	if err == (*responses.Error)(nil) {
		t.Errorf("file %s wasn't correctly deleted from database", fileId)
	}

	err = database.Cluster.Close(nil)
	if err != nil {
		t.Errorf("cannot close database instance : %s", err.Error())
		return
	}
}
