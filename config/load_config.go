package config

import (
	"github.com/Alvarios/gcfs/config/errors"
	"os"
)

var Main Configuration

type Configuration struct {
	Database Database               `json:"database"`
	Server   Server                 `json:"server"`
	Routes   Routes                 `json:"routes"`
	Metadata map[string]interface{} `json:"metadata"`
	Global   Global                 `json:"global"`
}

func LoadConfigForTest(c Configuration) {
	LoadConfig(c)
	username := os.Getenv("GCFS_TEST_USERNAME")
	password := os.Getenv("GCFS_TEST_PASSWORD")
	bucketName := os.Getenv("GCFS_TEST_BUCKETNAME")
	address := os.Getenv("GCFS_TEST_ADDRESS")

	Main.Database.Username = username
	Main.Database.Password = password

	if bucketName != "" {
		Main.Database.BucketName = bucketName
	}

	if address != "" {
		Main.Database.Address = address
	}
}

func LoadConfig(c Configuration) {
	Main = c
	errors.LoadErrors()

	if Main.Database.Address == "" {
		Main.Database.Address = "couchbase://127.0.0.1"
	}

	if Main.Database.BucketName == "" {
		Main.Database.BucketName = "metadata"
	}

	if Main.Server.Port == "" {
		Main.Server.Port = "8080"
	}
}