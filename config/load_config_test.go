package config

import (
	"testing"
)

func TestLoadConfig(t *testing.T) {
	LoadConfig(Configuration{})

	if Main.Database.Address != "couchbase://127.0.0.1" {
		t.Errorf("Unexpected value during config initialization : expected couchbase://127.0.0.1 for Database.Address, got %s", Main.Database.Address)
	}

	if Main.Database.BucketName != "metadata" {
		t.Errorf("Unexpected value during config initialization : expected metadata for Database.BucketName, got %s", Main.Database.BucketName)
	}

	if Main.Database.Username != "" {
		t.Errorf("Unexpected value during config initialization : expected empty for Database.Username, got %s", Main.Database.Username)
	}

	if Main.Database.Password != "" {
		t.Errorf("Unexpected value during config initialization : expected empty for Database.Password, got %s", Main.Database.Password)
	}
}
